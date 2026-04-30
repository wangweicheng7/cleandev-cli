package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"cleandev-cli/internal/core"
)

func main() {
	cliCfg := core.ParseArgs(os.Args)
	if cliCfg.Command == "config" && hasArg("init", os.Args) {
		path, err := core.InitConfigTemplate()
		if err != nil {
			fmt.Fprintf(os.Stderr, "create config failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Config template created at: %s\n", path)
		return
	}
	if cliCfg.Command == "doctor" {
		runDoctor()
		return
	}

	cfg := core.ResolveConfig(cliCfg)
	candidates := core.ScanCandidates(cfg)
	plan := core.BuildPlan(candidates)

	switch cliCfg.Command {
	case "scan", "plan":
		if cfg.JSONOutput {
			output := map[string]any{
				"profile":      cfg.Profile,
				"scan_results": candidates,
				"plan":         plan,
			}
			b, _ := json.MarshalIndent(output, "", "  ")
			fmt.Println(string(b))
		} else {
			fmt.Println(core.FormatPlanText(plan))
			fmt.Println("\nTip: use --json for structured output.")
		}
	case "clean":
		if !cfg.Confirm {
			fmt.Fprintln(os.Stderr, "Refusing to clean without --confirm. Run `cleaner plan` first.")
			os.Exit(1)
		}
		result := core.ExecuteClean(plan, cfg.AuditLogDir)
		b, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(b))
	default:
		printUsage()
	}
}

func hasArg(target string, args []string) bool {
	for _, arg := range args {
		if arg == target {
			return true
		}
	}
	return false
}

func runDoctor() {
	home, _ := os.UserHomeDir()
	checks := []struct {
		Name  string
		Path  string
		Exist bool
	}{
		{Name: "Home exists", Path: home, Exist: core.PathExists(home)},
		{Name: "Library exists", Path: filepath.Join(home, "Library"), Exist: core.PathExists(filepath.Join(home, "Library"))},
	}
	fmt.Println("=== Cleaner Doctor ===")
	for _, c := range checks {
		status := "WARN"
		if c.Exist {
			status = "OK"
		}
		fmt.Printf("- %s: %s (%s)\n", status, c.Name, c.Path)
	}
}

func printUsage() {
	fmt.Println(`Usage:
  cleaner scan [--profile dev] [--json] [--category cache,logs]
  cleaner plan [--profile dev] [--json] [--category cache,logs]
  cleaner clean --confirm [--profile dev] [--category cache,logs]
  cleaner doctor
  cleaner config init`)
}
