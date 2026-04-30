package core

import (
	"os"
	"path/filepath"
	"testing"

	"cleandev-cli/internal/types"
)

func TestProfileFilter(t *testing.T) {
	rules := DefaultRules()
	foundNPM := false
	for _, r := range rules {
		if r.ID == "npm-cache" {
			foundNPM = true
		}
	}
	if !foundNPM {
		t.Fatalf("expected npm-cache rule")
	}
}

func TestProtectedPathSkipped(t *testing.T) {
	tmp := t.TempDir()
	protected := filepath.Join(tmp, "protected")
	if err := os.MkdirAll(protected, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(protected, "junk.log"), []byte("abc"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg := types.Config{
		Profile: types.ProfileDev,
		Rules: []types.Rule{
			{ID: "custom", Name: "Custom", Category: "cache", Paths: []string{protected}, RiskLevel: types.RiskLow, Profiles: []types.Profile{types.ProfileDev}, DeleteStrategy: "remove_dir_if_safe"},
		},
		Categories:     map[string]struct{}{},
		ProtectedPaths: []string{protected},
	}
	results := ScanCandidates(cfg)
	if len(results) != 1 || results[0].Status != "skipped_protected" {
		t.Fatalf("expected skipped protected, got %+v", results)
	}
}

func TestCleanRemoveAndReportOnly(t *testing.T) {
	tmp := t.TempDir()
	deletable := filepath.Join(tmp, "delete-me")
	reportOnly := filepath.Join(tmp, "report-only")
	_ = os.MkdirAll(deletable, 0o755)
	_ = os.MkdirAll(reportOnly, 0o755)
	_ = os.WriteFile(filepath.Join(deletable, "f.txt"), []byte("x"), 0o644)
	_ = os.WriteFile(filepath.Join(reportOnly, "f.txt"), []byte("y"), 0o644)

	plan := BuildPlan([]types.Candidate{
		{Category: "cache", Path: deletable, RiskLevel: types.RiskLow, Reclaimable: 1, ReportOnly: false, Status: "ready"},
		{Category: "dev-cache", Path: reportOnly, RiskLevel: types.RiskHigh, Reclaimable: 1, ReportOnly: true, Status: "ready"},
	})
	res := ExecuteClean(plan, filepath.Join(tmp, "logs"))
	if len(res.Deleted) != 1 || len(res.Skipped) != 1 {
		t.Fatalf("unexpected clean result: %+v", res)
	}
	if _, err := os.Stat(deletable); !os.IsNotExist(err) {
		t.Fatalf("expected deletable removed")
	}
	if _, err := os.Stat(reportOnly); err != nil {
		t.Fatalf("expected report-only path kept")
	}
}

func TestConfigPriority(t *testing.T) {
	tmp := t.TempDir()
	prev, _ := os.Getwd()
	defer os.Chdir(prev)
	_ = os.Chdir(tmp)

	_ = os.WriteFile(".cleandevrc.json", []byte(`{"profile":"safe","include_paths":["~/foo"]}`), 0o644)
	explicit := filepath.Join(tmp, "custom-config.json")
	_ = os.WriteFile(explicit, []byte(`{"profile":"aggressive","include_paths":["~/bar"]}`), 0o644)

	cfg := ResolveConfig(types.Config{Command: "scan", Profile: types.ProfileDev, ConfigPath: explicit})
	if cfg.Profile != types.ProfileDev {
		t.Fatalf("expected cli profile dev, got %s", cfg.Profile)
	}
	if len(cfg.IncludePaths) < 2 {
		t.Fatalf("expected merged include paths")
	}
}
