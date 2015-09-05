package cli

import "fmt"

// Version is the current version of the quantum cli
const Version = "0.0.2"

// VersionCli prints out the current version
func VersionCli(args []string) {
	fmt.Println(Version)
}
