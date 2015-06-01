package cli

import "fmt"

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
