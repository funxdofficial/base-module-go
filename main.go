package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/funxdofficial/base-module-go/funbuild"
)

const (
	cliName    = "base-module-go"
	cliVersion = "1.0.0"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "new" {
		runNew(os.Args[2:])
		return
	}
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "help") {
		printUsage()
		return
	}
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version" || os.Args[1] == "version") {
		fmt.Printf("%s %s\n", cliName, cliVersion)
		return
	}

	runGenerate(os.Args[1:])
}

func runNew(args []string) {
	cfg := parseNewArgs(args)

	if cfg.PkgName == "" {
		fmt.Fprintln(os.Stderr, "error: service name is required")
		printUsage()
		os.Exit(2)
	}

	if err := funbuild.Create(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func parseNewArgs(args []string) funbuild.Config {
	var cfg funbuild.Config
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--output" && i+1 < len(args):
			cfg.Output = args[i+1]
			i++
		case strings.HasPrefix(arg, "--output="):
			cfg.Output = strings.TrimPrefix(arg, "--output=")
		case arg == "--pkg" && i+1 < len(args):
			cfg.PkgName = args[i+1]
			i++
		case strings.HasPrefix(arg, "--pkg="):
			cfg.PkgName = strings.TrimPrefix(arg, "--pkg=")
		case arg == "--type" && i+1 < len(args):
			cfg.ServiceType = args[i+1]
			i++
		case strings.HasPrefix(arg, "--type="):
			cfg.ServiceType = strings.TrimPrefix(arg, "--type=")
		case strings.HasPrefix(arg, "-"):
			fmt.Fprintf(os.Stderr, "error: unknown flag %s\n", arg)
			os.Exit(2)
		case cfg.PkgName == "":
			cfg.PkgName = arg
		default:
			fmt.Fprintf(os.Stderr, "error: unexpected argument %s\n", arg)
			os.Exit(2)
		}
	}
	return cfg
}

func runGenerate(args []string) {
	fs := flag.NewFlagSet(cliName, flag.ExitOnError)
	var cfg funbuild.Config
	fs.StringVar(&cfg.PkgName, "pkg", "", "service name / output folder (required)")
	fs.StringVar(&cfg.Output, "output", "", "output directory (default: ./<pkg>)")
	fs.StringVar(&cfg.ServiceType, "type", funbuild.ServiceTypeREST, "service type: rest or consumer (alias: cons)")
	fs.Parse(args)

	if cfg.PkgName == "" {
		fmt.Fprintln(os.Stderr, "error: --pkg is required")
		printUsage()
		os.Exit(2)
	}

	if err := funbuild.Create(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	usage := `%s %s — Go service scaffold generator

Usage:
  %s --pkg=<service-name> [--type=rest|consumer|cons] [--output=<dir>]
  %s new <service-name> [--type=rest|consumer|cons] [--output=<dir>]

Service types:
  rest      REST API with Echo CRUD (default, no scheduler)
  consumer  Cron/background jobs via golang-module-scheduler (no HTTP)
  cons      Alias for consumer

Examples:
  base-module-go --pkg=order-service --type=rest
  base-module-go --pkg=sync-worker --type=consumer
  make fun-build PKG=order-service TYPE=rest

Install:
  go install github.com/funxdofficial/base-module-go@latest

Generated module path: github.com/funxdofficial/<service-name>
`
	fmt.Printf(usage, cliName, cliVersion, cliName, cliName)
}
