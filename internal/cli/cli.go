package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/raeseoklee/scanrail/internal/app"
	"github.com/raeseoklee/scanrail/internal/exitcode"
	"github.com/raeseoklee/scanrail/internal/mcpserver"
)

func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		usage(stdout)
		return exitcode.OK
	}
	switch args[0] {
	case "--version", "-v", "version":
		return app.Version(stdout)
	case "--help", "-h", "help":
		usage(stdout)
		return exitcode.OK
	case "doctor":
		return app.Doctor(stdout)
	case "init":
		return runInit(args[1:], stdout, stderr)
	case "setup":
		return runSetup(args[1:], stdout, stderr)
	case "run":
		return runScan(args[1:], stdout, stderr)
	case "mcp":
		return runMCP(args[1:], stdout, stderr)
	default:
		fmt.Fprintln(stderr, "unknown command:", args[0])
		usage(stderr)
		return exitcode.ConfigError
	}
}

func runMCP(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) != 1 || args[0] != "serve" {
		fmt.Fprintln(stderr, "usage: scanrail mcp serve")
		return exitcode.ConfigError
	}
	return mcpserver.Serve(context.Background(), os.Stdin, stdout, stderr)
}

func runInit(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var opts app.InitOptions
	nonInteractive := fs.Bool("non-interactive", false, "generate config without prompts")
	fs.StringVar(&opts.ConfigPath, "config", "scanrail.yaml", "config path")
	fs.StringVar(&opts.ProjectName, "project-name", "", "project name")
	fs.StringVar(&opts.Target, "target", "", "web target URL")
	fs.StringVar(&opts.Profile, "profile", "quick", "default profile")
	fs.BoolVar(&opts.Force, "force", false, "overwrite existing files")
	if err := fs.Parse(args); err != nil {
		return exitcode.ConfigError
	}
	if !*nonInteractive {
		fmt.Fprintln(stderr, "interactive init is not implemented yet; use --non-interactive")
		return exitcode.ConfigError
	}
	return app.Init(opts, stdout)
}

func runSetup(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("setup", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var opts app.SetupOptions
	fs.StringVar(&opts.PullPolicy, "pull-policy", "missing", "missing, always, or never")
	if err := fs.Parse(args); err != nil {
		return exitcode.ConfigError
	}
	return app.Setup(opts, stdout)
}

func runScan(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	fs.SetOutput(stderr)
	var opts app.RunOptions
	fs.StringVar(&opts.ConfigPath, "config", "scanrail.yaml", "config path")
	fs.StringVar(&opts.Profile, "profile", "", "profile")
	fs.StringVar(&opts.Target, "target", "", "web target URL")
	fs.StringVar(&opts.Only, "only", "", "run only one tool")
	fs.StringVar(&opts.OutputDir, "output-dir", "", "report output directory")
	if err := fs.Parse(args); err != nil {
		return exitcode.ConfigError
	}
	return app.Run(context.Background(), opts, stdout)
}

func usage(w io.Writer) {
	fmt.Fprintln(w, `Scanrail - developer-first security scan orchestrator

Usage:
  scanrail doctor
  scanrail init --non-interactive --project-name demo --target http://localhost:8080
  scanrail setup [--pull-policy never]
  scanrail run [--profile quick] [--only headers|gitleaks|tls]
  scanrail mcp serve

Options:
  --version    print version
  --help       print help`)
}
