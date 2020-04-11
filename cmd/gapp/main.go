package main

import (
	"github.com/eikc/gapp/internal/cli"
	"os"
)

func main() {
	c := cli.NewCLI()

	os.Exit(c.Do(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
