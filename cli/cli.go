package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/doubledutch/quantum"
)

// CLI wraps the CLIs defined in quantum
type CLI struct {
	cmds map[string]Func
}

// Func wraps a cli run func
type Func func(args []string)

// New returns a new CLI
func New() *CLI {
	cmds := make(map[string]Func)
	cmds["solo"] = SoloCli
	cmds["client"] = ClientCli

	return &CLI{cmds: cmds}
}

// NewRequest creates a request by treating requestData as file and falling back
// to treating it as a string.
func NewRequest(requestType, requestData string) (r quantum.Request, err error) {
	if _, err := os.Stat(requestData); err == nil {
		f, err := os.Open(requestData)
		if err != nil {
			return r, err
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			return r, err
		}

		requestData = string(b)
	}

	return quantum.NewRequest(requestType, requestData), nil
}

// Run runs a CLI with specified arguments
func (c *CLI) Run(args []string) {
	if len(args) < 2 {
		fmt.Println(c.Usage())
		return
	}

	mode := args[1]
	cmdArgs := args[2:]
	cmd, ok := c.cmds[mode]
	if !ok {
		fmt.Println(c.Usage())
		return
	}

	cmd(cmdArgs)
}

// Usage is the usage for the CLI
func (c *CLI) Usage() string {
	return "quantum usage: \n" +
		"\tclient:\trun in client mode\n" +
		"\tsolo:\trun in stand alone mode"
}
