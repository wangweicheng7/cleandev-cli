package core

import (
	"os"

	"cleandev-cli/internal/types"
)

type CleanResult struct {
	Deleted            []types.Action   `json:"deleted"`
	Skipped            []types.Action   `json:"skipped"`
	Errors             []map[string]any `json:"errors"`
	TotalReclaimedByte int64            `json:"total_reclaimed_bytes"`
}

func ExecuteClean(plan types.Plan, auditLogDir string) CleanResult {
	result := CleanResult{
		Deleted: make([]types.Action, 0),
		Skipped: make([]types.Action, 0),
		Errors:  make([]map[string]any, 0),
	}
	for _, action := range plan.Actions {
		if action.Action != "delete" {
			result.Skipped = append(result.Skipped, action)
			continue
		}
		err := os.RemoveAll(action.Path)
		if err != nil {
			result.Errors = append(result.Errors, map[string]any{
				"path":  action.Path,
				"error": err.Error(),
			})
			_ = WriteAuditLog(auditLogDir, map[string]any{
				"operation": "delete",
				"path":      action.Path,
				"result":    "error",
				"error":     err.Error(),
			})
			continue
		}
		result.Deleted = append(result.Deleted, action)
		result.TotalReclaimedByte += action.Reclaimable
		_ = WriteAuditLog(auditLogDir, map[string]any{
			"operation":        "delete",
			"path":             action.Path,
			"reclaimable_size": action.Reclaimable,
			"result":           "success",
		})
	}
	return result
}
