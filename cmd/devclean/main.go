package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/wangweicheng7/devclean-cli/internal/cli"
	"github.com/wangweicheng7/devclean-cli/internal/version"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	ctx := context.Background()

	// We keep root flags minimal; each subcommand has its own flagset.
	root := flag.NewFlagSet("devclean", flag.ContinueOnError)
	root.SetOutput(os.Stderr)
	showVersion := root.Bool("version", false, "print version and exit")
	showVersionShort := root.Bool("v", false, "print version and exit")
	_ = root.Parse(os.Args[1:])
	if *showVersion || *showVersionShort {
		fmt.Fprintln(os.Stdout, version.String())
		return 0
	}

	args := root.Args()
	if len(args) == 0 {
		cli.PrintHelp(os.Stdout)
		return 2
	}

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case "version", "--version", "-v":
		fmt.Fprintln(os.Stdout, version.String())
		return 0
	case "scan":
		return cli.RunScan(ctx, cmdArgs, os.Stdout, os.Stderr)
	case "plan":
		// Alias for scan (kept for nicer semantics).
		return cli.RunScan(ctx, cmdArgs, os.Stdout, os.Stderr)
	case "clean":
		return cli.RunClean(ctx, cmdArgs, os.Stdout, os.Stderr)
	case "config":
		return cli.RunConfig(ctx, cmdArgs, os.Stdout, os.Stderr)
	case "doctor":
		return cli.RunDoctor(ctx, cmdArgs, os.Stdout, os.Stderr)
	case "help", "-h", "--help":
		cli.PrintHelp(os.Stdout)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		cli.PrintHelp(os.Stderr)
		return 2
	}
}
