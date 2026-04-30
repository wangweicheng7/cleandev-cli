package core

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

func WriteAuditLog(logDir string, payload map[string]any) error {
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return err
	}
	record := map[string]any{"time": time.Now().UTC().Format(time.RFC3339)}
	for k, v := range payload {
		record[k] = v
	}
	body, _ := json.Marshal(record)
	logPath := filepath.Join(logDir, "audit.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(body, '\n'))
	return err
}
