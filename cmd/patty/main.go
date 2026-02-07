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
	init			Initialize a new patty project
	install [packages...]	Install packages (adds to patty.toml and installs)
	remove <package>	Remove a dependency from the project
	update			Reinstall all dependencies
	help			Show this help message

Aliases:
	i			Short for install
	rm			Short for remove

Examples:
	patty init
	patty install luafilesystem
	patty install lua-json@1.0 luasocket
	patty remove luafilesystem
	patty update`)
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
			printError(err)
			os.Exit(1)
		}
	case "install", "i":
		pkgs := os.Args[2:]
		if err := patty.CmdInstall(pkgs); err != nil {
			printError(err)
			os.Exit(1)
		}
	case "remove", "rm":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Error: missing package name")
			fmt.Fprintln(os.Stderr, "  Usage: patty remove <package>")
			fmt.Fprintln(os.Stderr, "  Example: patty remove luafilesystem")
			os.Exit(1)
		}
		pkg := os.Args[2]
		if err := patty.CmdRemove(pkg); err != nil {
			printError(err)
			os.Exit(1)
		}
	case "update":
		if err := patty.CmdUpdate(); err != nil {
			printError(err)
			os.Exit(1)
		}
	case "help", "--help", "-h":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n\n", cmd)
		fmt.Fprintln(os.Stderr, "  Run 'patty help' to see available commands")
		os.Exit(1)
	}
}

func printError(err error) {
	fmt.Fprintf(os.Stderr, "\nError: %s\n", err)
}
