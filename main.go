package main

import (
	"os"

	"github.com/doubledutch/quantum-cli/cli"
)

func main() {
	cli := cli.New()
	cli.Run(os.Args)
}
