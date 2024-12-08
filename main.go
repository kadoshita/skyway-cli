package main

import (
	"fmt"
	"os"

	"github.com/kadoshita/skyway-cli/cmd"
)

func main() {
	if os.Getenv("SKYWAY_CLI_GEN_DOCS") == "true" {
		err := cmd.GenDocs()
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	cmd.Execute()
}
