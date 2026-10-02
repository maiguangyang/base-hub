/**
 * @Author: Marlon.M
 * @Email: maiguangyang@163.com
 * @Date: 2026-09-21
 */
package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"base-engine/src/dbup"
	"github.com/urfave/cli"
)

var bootstrapHQCmd = cli.Command{
	Name:  "bootstrap-hq",
	Usage: "create the initial headquarters administrator",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "phone", Usage: "administrator phone"},
		cli.StringFlag{Name: "display-name", Usage: "administrator display name"},
	},
	Action: runBootstrapHQCommand,
}

func runBootstrapHQCommand(cliContext *cli.Context) error {
	phone := strings.TrimSpace(cliContext.String("phone"))
	displayName := strings.TrimSpace(cliContext.String("display-name"))
	if phone == "" || displayName == "" {
		return errors.New("VALIDATION_FAILED")
	}
	db, err := openEngineDBFromEnv()
	if err != nil {
		return err
	}
	defer db.Close()
	result, err := dbup.BootstrapHQ(context.Background(), db.Query(), phone, displayName)
	if err != nil {
		return err
	}
	return json.NewEncoder(os.Stdout).Encode(result)
}
