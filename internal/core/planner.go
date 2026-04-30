package core

import (
	"fmt"
	"strings"

	"cleandev-cli/internal/types"
)

func BuildPlan(candidates []types.Candidate) types.Plan {
	summary := types.PlanSummary{
		TotalItems:       len(candidates),
		ReclaimableBytes: 0,
		ByCategory:       map[string]int64{},
		ByRisk:           map[types.RiskLevel]int64{},
	}
	actions := make([]types.Action, 0, len(candidates))
	for _, c := range candidates {
		summary.ReclaimableBytes += c.Reclaimable
		summary.ByCategory[c.Category] += c.Reclaimable
		summary.ByRisk[c.RiskLevel] += c.Reclaimable
		action := "delete"
		if c.ReportOnly || c.Status != "ready" {
			action = "report_only"
		}
		actions = append(actions, types.Action{
			Path:        c.Path,
			Category:    c.Category,
			RiskLevel:   c.RiskLevel,
			Reclaimable: c.Reclaimable,
			Action:      action,
			Reason:      c.Reason,
		})
	}
	return types.Plan{Summary: summary, Actions: actions}
}

func FormatPlanText(plan types.Plan) string {
	var b strings.Builder
	b.WriteString("=== Cleaner Plan Preview ===\n")
	b.WriteString(fmt.Sprintf("Total candidates: %d\n", plan.Summary.TotalItems))
	b.WriteString(fmt.Sprintf("Potential reclaim: %s\n\n", formatBytes(plan.Summary.ReclaimableBytes)))
	b.WriteString("By category:\n")
	for category, size := range plan.Summary.ByCategory {
		b.WriteString(fmt.Sprintf("- %s: %s\n", category, formatBytes(size)))
	}
	b.WriteString("\nActions:\n")
	for _, action := range plan.Actions {
		b.WriteString(fmt.Sprintf("- [%s] %s (%s, %s, %s)\n", action.Action, action.Path, action.Category, action.RiskLevel, formatBytes(action.Reclaimable)))
	}
	return b.String()
}

func formatBytes(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}
	units := []string{"KB", "MB", "GB", "TB"}
	value := float64(size)
	idx := 0
	for value >= 1024 && idx < len(units)-1 {
		value /= 1024
		idx++
	}
	return fmt.Sprintf("%.1f %s", value, units[idx])
}
