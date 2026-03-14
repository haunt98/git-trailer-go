package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
)

var flagSessionID string

func init() {
	flag.StringVar(&flagSessionID, "session", "", "which sessionID to export")
}

func main() {
	flag.Parse()

	ctx := context.Background()

	if flagSessionID == "" {
		var err error
		flagSessionID, err = GetRecentSessionID(ctx)
		if err != nil {
			slog.Error("GetRecentSessionID", "error", err)
			return
		}
	}

	if flagSessionID == "" {
		return
	}

	sessionExportModels, err := GetSessionExportModels(ctx, flagSessionID)
	if err != nil {
		slog.Error("GetSessionExportModels", "error", err)
		return
	}

	if len(sessionExportModels) == 0 {
		return
	}

	for _, model := range sessionExportModels {
		fmt.Printf("Co-Authored-By: OpenCode %s %s <noreply@opencode.ai>\n",
			model.ProviderID,
			model.ModelID,
		)
	}
}
