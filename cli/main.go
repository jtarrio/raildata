// An example application to use the raildata library.
package main

import (
	"os"

	"github.com/jtarrio/raildata/raildata-cli/cmd"
)

func main() {
	if err := cmd.App().Run(os.Args); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
