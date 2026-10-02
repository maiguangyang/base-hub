package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// limits 是新手写 Go 源码的默认上限。
var limits = metrics{Lines: 250, Complexity: 10, Depth: 3}

// main 运行只读源码结构检查。
func main() {
	root := flag.String("root", ".", "Engine repository root")
	report := flag.Bool("report", false, "print all measured source sizes")
	flag.Parse()
	if err := run(*root, *report); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// run 逐文件检查，并在任一度量超限时返回错误。
func run(root string, report bool) error {
	baseline, err := readBaseline(filepath.Join(root, "tools/sourcecheck/baseline.json"))
	if err != nil {
		return err
	}
	files, err := sourceFiles(root)
	if err != nil {
		return err
	}
	sort.Strings(files)
	violations := 0
	for _, path := range files {
		measurements, err := measureFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		for _, item := range measurements {
			if checkMeasurement(filepath.ToSlash(relative), item, baseline, report) {
				violations++
			}
		}
	}
	if violations != 0 {
		return fmt.Errorf("source structure: %d violations", violations)
	}
	return nil
}

// checkMeasurement 比较新代码上限或既有单元的冻结基线。
func checkMeasurement(path string, item measurement, baseline map[string]metrics, report bool) bool {
	key := path
	limit := limits
	if item.Name != "" {
		key += "#" + item.Name
		limit.Lines = 50
	}
	if recorded, ok := baseline[key]; ok {
		limit = recorded
	}
	if report {
		fmt.Printf("%s: lines=%d complexity=%d depth=%d\n", key, item.Metrics.Lines, item.Metrics.Complexity, item.Metrics.Depth)
	}
	if item.Metrics.Lines <= limit.Lines && item.Metrics.Complexity <= limit.Complexity && item.Metrics.Depth <= limit.Depth {
		return false
	}
	fmt.Fprintf(os.Stderr, "%s exceeds limit: actual %+v, allowed %+v\n", key, item.Metrics, limit)
	return true
}

// readBaseline 加载既有超限单元的上限记录。
func readBaseline(path string) (map[string]metrics, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var baseline map[string]metrics
	if err := json.Unmarshal(data, &baseline); err != nil {
		return nil, err
	}
	return baseline, nil
}
