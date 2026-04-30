package core

import (
	"os"
	"path/filepath"

	"cleandev-cli/internal/types"
)

func DefaultRules() []types.Rule {
	home, _ := os.UserHomeDir()
	hp := func(parts ...string) string {
		return filepath.Join(append([]string{home}, parts...)...)
	}
	return []types.Rule{
		{ID: "user-library-cache", Name: "User Library Caches", Category: "cache", Paths: []string{hp("Library", "Caches")}, RiskLevel: types.RiskLow, Profiles: []types.Profile{types.ProfileSafe, types.ProfileDev, types.ProfileAggressive}, DeleteStrategy: "remove_dir_if_safe"},
		{ID: "user-library-logs", Name: "User Library Logs", Category: "logs", Paths: []string{hp("Library", "Logs")}, RiskLevel: types.RiskLow, Profiles: []types.Profile{types.ProfileSafe, types.ProfileDev, types.ProfileAggressive}, DeleteStrategy: "remove_dir_if_safe"},
		{ID: "npm-cache", Name: "npm cache", Category: "dev-cache", Paths: []string{hp(".npm")}, RiskLevel: types.RiskLow, Profiles: []types.Profile{types.ProfileDev, types.ProfileAggressive}, DeleteStrategy: "remove_dir_if_safe"},
		{ID: "pnpm-cache", Name: "pnpm store", Category: "dev-cache", Paths: []string{hp("Library", "pnpm", "store")}, RiskLevel: types.RiskLow, Profiles: []types.Profile{types.ProfileDev, types.ProfileAggressive}, DeleteStrategy: "remove_dir_if_safe"},
		{ID: "yarn-cache", Name: "yarn cache", Category: "dev-cache", Paths: []string{hp("Library", "Caches", "Yarn"), hp(".cache", "yarn")}, RiskLevel: types.RiskLow, Profiles: []types.Profile{types.ProfileDev, types.ProfileAggressive}, DeleteStrategy: "remove_dir_if_safe"},
		{ID: "pip-cache", Name: "pip cache", Category: "dev-cache", Paths: []string{hp("Library", "Caches", "pip"), hp(".cache", "pip")}, RiskLevel: types.RiskLow, Profiles: []types.Profile{types.ProfileDev, types.ProfileAggressive}, DeleteStrategy: "remove_dir_if_safe"},
		{ID: "xcode-derived-data", Name: "Xcode DerivedData", Category: "dev-cache", Paths: []string{hp("Library", "Developer", "Xcode", "DerivedData")}, RiskLevel: types.RiskMedium, Profiles: []types.Profile{types.ProfileDev, types.ProfileAggressive}, DeleteStrategy: "remove_dir_if_safe"},
		{ID: "homebrew-cache", Name: "Homebrew cache", Category: "dev-cache", Paths: []string{hp("Library", "Caches", "Homebrew")}, RiskLevel: types.RiskMedium, Profiles: []types.Profile{types.ProfileAggressive}, DeleteStrategy: "remove_dir_if_safe"},
		{ID: "docker-data", Name: "Docker data", Category: "dev-cache", Paths: []string{hp("Library", "Containers", "com.docker.docker")}, RiskLevel: types.RiskHigh, Profiles: []types.Profile{types.ProfileAggressive}, DeleteStrategy: "report_only"},
	}
}

func DefaultProtectedPaths() []string {
	home, _ := os.UserHomeDir()
	return []string{
		"/",
		"/System",
		"/usr",
		"/Applications",
		filepath.Join(home, "Documents"),
		filepath.Join(home, "Desktop"),
		filepath.Join(home, "development"),
		filepath.Join(home, "Projects"),
		filepath.Join(home, "WeChatProjects"),
		filepath.Join(home, "go", "src"),
	}
}
