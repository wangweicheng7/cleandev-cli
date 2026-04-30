package core

import (
	"os"
	"time"

	"cleandev-cli/internal/types"
)

func ScanCandidates(cfg types.Config) []types.Candidate {
	results := make([]types.Candidate, 0)
	for _, rule := range cfg.Rules {
		if !containsProfile(rule.Profiles, cfg.Profile) {
			continue
		}
		if len(cfg.Categories) > 0 {
			if _, ok := cfg.Categories[rule.Category]; !ok {
				continue
			}
		}
		for _, p := range rule.Paths {
			target := ExpandHome(p)
			if !PathExists(target) {
				continue
			}
			if isProtected(target, cfg.ProtectedPaths) {
				results = append(results, buildCandidate(rule, target, 0, "skipped_protected", true, "Path is protected by hard rule"))
				continue
			}
			if isProtected(target, cfg.ExcludePaths) {
				results = append(results, buildCandidate(rule, target, 0, "skipped_excluded", true, "Path excluded by config"))
				continue
			}
			fi, err := os.Stat(target)
			if err != nil {
				continue
			}
			size := fi.Size()
			if fi.IsDir() {
				size = DirSize(target)
			}
			reportOnly := rule.DeleteStrategy == "report_only" || rule.RiskLevel == types.RiskHigh
			results = append(results, buildCandidate(rule, target, size, "ready", reportOnly, ""))
		}
	}

	for _, p := range cfg.IncludePaths {
		target := ExpandHome(p)
		if !PathExists(target) || isProtected(target, cfg.ProtectedPaths) {
			continue
		}
		fi, err := os.Stat(target)
		if err != nil {
			continue
		}
		size := fi.Size()
		if fi.IsDir() {
			size = DirSize(target)
		}
		results = append(results, types.Candidate{
			RuleID:          "custom-include",
			Name:            "Custom include path",
			Category:        "custom",
			Path:            target,
			Size:            size,
			Reclaimable:     size,
			RiskLevel:       types.RiskMedium,
			Status:          "ready",
			ReportOnly:      false,
			Reason:          "",
			DeleteStrategy:  "remove_dir_if_safe",
			LastAccessOrMod: time.Now().UTC().Format(time.RFC3339),
		})
	}
	return results
}

func containsProfile(profiles []types.Profile, profile types.Profile) bool {
	for _, p := range profiles {
		if p == profile {
			return true
		}
	}
	return false
}

func isProtected(target string, protected []string) bool {
	for _, p := range protected {
		if IsSubPath(target, ExpandHome(p)) {
			return true
		}
	}
	return false
}

func buildCandidate(rule types.Rule, target string, size int64, status string, reportOnly bool, reason string) types.Candidate {
	return types.Candidate{
		RuleID:          rule.ID,
		Name:            rule.Name,
		Category:        rule.Category,
		Path:            target,
		Size:            size,
		Reclaimable:     size,
		RiskLevel:       rule.RiskLevel,
		Status:          status,
		ReportOnly:      reportOnly,
		Reason:          reason,
		DeleteStrategy:  rule.DeleteStrategy,
		LastAccessOrMod: time.Now().UTC().Format(time.RFC3339),
	}
}
