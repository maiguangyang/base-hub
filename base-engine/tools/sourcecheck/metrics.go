package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// metrics 记录一个源码单元的有效行、圈复杂度和控制流嵌套。
type metrics struct {
	Lines      int `json:"lines"`                // 有效代码行数。
	Complexity int `json:"complexity,omitempty"` // 圈复杂度。
	Depth      int `json:"depth,omitempty"`      // 最大控制流嵌套层数。
}

// measurement 用名称关联一个文件或函数的度量结果。
type measurement struct {
	Name    string  // 函数名称；空值表示整个文件。
	Metrics metrics // 对应源码单元的度量结果。
}

// measureFile 只解析指定文件，不写入任何源码。
func measureFile(path string) ([]measurement, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	set := token.NewFileSet()
	file, err := parser.ParseFile(set, path, source, 0)
	if err != nil {
		return nil, err
	}
	active := activeLines(source)
	result := []measurement{{Name: "", Metrics: metrics{Lines: countLines(active, 1, len(active)-1)}}}
	for _, declaration := range file.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok || function.Body == nil {
			continue
		}
		name := function.Name.Name
		if function.Recv != nil {
			name = receiverName(function.Recv) + "." + name
		}
		start := set.Position(function.Pos()).Line
		end := set.Position(function.End()).Line
		result = append(result, measurement{Name: name, Metrics: metrics{
			Lines: countLines(active, start, end), Complexity: complexity(function.Body), Depth: controlDepth(function.Body),
		}})
	}
	closureLines := map[int]int{}
	ast.Inspect(file, func(node ast.Node) bool {
		literal, ok := node.(*ast.FuncLit)
		if !ok {
			return true
		}
		start := set.Position(literal.Pos()).Line
		end := set.Position(literal.End()).Line
		closureLines[start]++
		name := fmt.Sprintf("closure@%d", start)
		if closureLines[start] > 1 {
			name = fmt.Sprintf("%s~%d", name, closureLines[start])
		}
		result = append(result, measurement{Name: name, Metrics: metrics{
			Lines: countLines(active, start, end), Complexity: complexity(literal.Body), Depth: controlDepth(literal.Body),
		}})
		return true
	})
	return result, nil
}

// activeLines 按 Go 词法单元标记非空、非注释行。
func activeLines(source []byte) []bool {
	set := token.NewFileSet()
	file := set.AddFile("source.go", -1, len(source))
	var scan scanner.Scanner
	scan.Init(file, source, nil, scanner.ScanComments)
	active := make([]bool, bytes.Count(source, []byte{'\n'})+2)
	for {
		position, kind, literal := scan.Scan()
		if kind == token.EOF {
			break
		}
		if kind == token.COMMENT || kind == token.SEMICOLON {
			continue
		}
		start := file.Position(position).Line
		end := start
		if literal != "" {
			end = file.Position(position + token.Pos(len(literal)-1)).Line
		}
		for line := start; line <= end && line < len(active); line++ {
			active[line] = true
		}
	}
	return active
}

// countLines 统计指定行区间内的有效行。
func countLines(active []bool, start, end int) int {
	count := 0
	for line := start; line <= end && line < len(active); line++ {
		if active[line] {
			count++
		}
	}
	return count
}

// receiverName 提取方法接收者类型名，供基线定位使用。
func receiverName(fields *ast.FieldList) string {
	if len(fields.List) == 0 {
		return "receiver"
	}
	var name string
	ast.Inspect(fields.List[0].Type, func(node ast.Node) bool {
		if id, ok := node.(*ast.Ident); ok && name == "" {
			name = id.Name
		}
		return name == ""
	})
	return name
}

// complexity 按分支、循环、case 和短路条件计算圈复杂度。
func complexity(body *ast.BlockStmt) int {
	value := 1
	ast.Inspect(body, func(node ast.Node) bool {
		switch current := node.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.CommClause:
			value++
		case *ast.BinaryExpr:
			if current.Op == token.LAND || current.Op == token.LOR {
				value++
			}
		}
		return true
	})
	return value
}

// controlDepth 返回函数体内的最大控制流嵌套层数。
func controlDepth(body *ast.BlockStmt) int {
	return walkDepth(body, 0)
}

// walkDepth 递归检查一个节点及其子控制流。
func walkDepth(node ast.Node, depth int) int {
	if isControl(node) {
		depth++
	}
	maxDepth := depth
	ast.Inspect(node, func(child ast.Node) bool {
		if child == nil || child == node {
			return true
		}
		if isControl(child) {
			candidate := walkDepth(child, depth)
			if candidate > maxDepth {
				maxDepth = candidate
			}
			return false
		}
		return true
	})
	return maxDepth
}

// isControl 判断节点是否增加嵌套层数。
func isControl(node ast.Node) bool {
	switch node.(type) {
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
		return true
	}
	return false
}

// sourceFiles 枚举手写 Go 源码，并跳过根目录下的 gen/。
func sourceFiles(root string) ([]string, error) {
	files := []string{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if excludedDirectory(root, path, entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_gen.go") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// excludedDirectory 集中定义扫描时不读取的目录。
func excludedDirectory(root, path, name string) bool {
	if path == filepath.Join(root, "gen") {
		return true
	}
	return path != root && (name == ".git" || name == "vendor" || name == "testdata")
}
