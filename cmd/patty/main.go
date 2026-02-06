package main

import (
	"fmt"
	"os"

	"github.com/Rahyulty/patty/internal/patty"
)

func usage() {
	fmt.Println(`patty - a modern lua dependency tool (MVP)
Usage:
	patty [command] [options]

Commands:
	init		Initialize a new patty project
	add		    Add a new dependency to the project
	remove		Remove a dependency from the project
	update		Update all dependencies to the latest version
	help		Show help for a command`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	cmd := os.Args[1]

	switch cmd {

	case "init":
		if err := patty.CmdInit(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "add":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: missing package name")
			os.Exit(1)
		}
		pkg := os.Args[2]
		if err := patty.CmdAdd(pkg); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "remove":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: missing package name")
			os.Exit(1)
		}
		pkg := os.Args[2]
		if err := patty.CmdRemove(pkg); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "update":
		if err := patty.CmdUpdate(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	case "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n", cmd)
		usage()
		os.Exit(1)
	}
}
