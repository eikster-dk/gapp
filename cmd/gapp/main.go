package main

import (
	"github.com/eikc/gapp/internal/cli"
	"os"
)

func main() {
	os.Exit(cli.Do(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
