package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"cleandev-cli/internal/types"
)

const defaultConfigName = ".cleandevrc.json"

func ParseArgs(args []string) types.Config {
	cfg := types.Config{
		Command:    "plan",
		Profile:    types.ProfileDev,
		Categories: map[string]struct{}{},
	}
	if len(args) > 1 {
		cfg.Command = args[1]
	}
	for i := 2; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--json":
			cfg.JSONOutput = true
		case "--confirm":
			cfg.Confirm = true
		case "--profile":
			if i+1 < len(args) {
				cfg.Profile = types.Profile(args[i+1])
				i++
			}
		case "--category":
			if i+1 < len(args) {
				for _, c := range strings.Split(args[i+1], ",") {
					c = strings.TrimSpace(c)
					if c != "" {
						cfg.Categories[c] = struct{}{}
					}
				}
				i++
			}
		case "--config":
			if i+1 < len(args) {
				cfg.ConfigPath = args[i+1]
				i++
			}
		}
	}
	return cfg
}

func ResolveConfig(cliCfg types.Config) types.Config {
	home, _ := os.UserHomeDir()
	defaultCfg := types.Config{
		Command:        cliCfg.Command,
		Profile:        types.ProfileDev,
		JSONOutput:     cliCfg.JSONOutput,
		Confirm:        cliCfg.Confirm,
		Categories:     cliCfg.Categories,
		IncludePaths:   []string{},
		ExcludePaths:   []string{},
		ProtectedPaths: DefaultProtectedPaths(),
		Rules:          DefaultRules(),
		AuditLogDir:    filepath.Join(home, ".cleandev", "logs"),
	}

	localCfg := readConfigFile(filepath.Join(".", defaultConfigName))
	explicitCfg := types.Config{}
	if cliCfg.ConfigPath != "" {
		explicitCfg = readConfigFile(ExpandHome(cliCfg.ConfigPath))
	}

	merged := defaultCfg
	merged = mergeConfig(merged, localCfg)
	merged = mergeConfig(merged, explicitCfg)
	merged = mergeConfig(merged, cliCfg)
	merged.IncludePaths = append(append(defaultCfg.IncludePaths, localCfg.IncludePaths...), explicitCfg.IncludePaths...)
	merged.ExcludePaths = append(append(defaultCfg.ExcludePaths, localCfg.ExcludePaths...), explicitCfg.ExcludePaths...)
	merged.ProtectedPaths = append(append(defaultCfg.ProtectedPaths, localCfg.ProtectedPaths...), explicitCfg.ProtectedPaths...)
	return merged
}

func InitConfigTemplate() (string, error) {
	target := filepath.Join(".", defaultConfigName)
	if PathExists(target) {
		return target, nil
	}
	template := map[string]any{
		"profile":         string(types.ProfileDev),
		"include_paths":   []string{},
		"exclude_paths":   []string{},
		"protected_paths": []string{},
		"rules":           []types.Rule{},
	}
	body, _ := json.MarshalIndent(template, "", "  ")
	err := os.WriteFile(target, append(body, '\n'), 0o644)
	return target, err
}

func readConfigFile(path string) types.Config {
	if path == "" || !PathExists(path) {
		return types.Config{}
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return types.Config{}
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(content, &raw); err != nil {
		return types.Config{}
	}
	cfg := types.Config{}
	if v, ok := raw["profile"]; ok {
		var s string
		_ = json.Unmarshal(v, &s)
		cfg.Profile = types.Profile(s)
	}
	_ = json.Unmarshal(raw["include_paths"], &cfg.IncludePaths)
	_ = json.Unmarshal(raw["exclude_paths"], &cfg.ExcludePaths)
	_ = json.Unmarshal(raw["protected_paths"], &cfg.ProtectedPaths)
	_ = json.Unmarshal(raw["rules"], &cfg.Rules)
	return cfg
}

func mergeConfig(base, override types.Config) types.Config {
	if override.Profile != "" {
		base.Profile = override.Profile
	}
	if override.ConfigPath != "" {
		base.ConfigPath = override.ConfigPath
	}
	if override.AuditLogDir != "" {
		base.AuditLogDir = override.AuditLogDir
	}
	if len(override.Rules) > 0 {
		base.Rules = override.Rules
	}
	if len(override.Categories) > 0 {
		base.Categories = override.Categories
	}
	base.JSONOutput = base.JSONOutput || override.JSONOutput
	base.Confirm = base.Confirm || override.Confirm
	return base
}
