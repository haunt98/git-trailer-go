package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
)

var (
	flagSessionID          string
	flagSkipModelsDevCache bool
)

func init() {
	flag.StringVar(&flagSessionID, "session", "", "which opencode sessionID to export")
	flag.BoolVar(&flagSkipModelsDevCache, "skip-models-dev-cache", false, "skip models.dev cache")
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

	modelsDevData, err := LoadModelsDevData(ctx, flagSkipModelsDevCache)
	if err != nil {
		slog.Error("LoadModelsDevData", "error", err)
		return
	}

	for _, model := range sessionExportModels {
		providerName, modelName := modelsDevData.LookupName(model.ProviderID, model.ModelID)
		fmt.Printf("Co-Authored-By: OpenCode - %s - %s <noreply@opencode.ai>\n",
			providerName,
			modelName,
		)
	}
}
