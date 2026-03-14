package main

import (
	"context"
	"fmt"
	"log"
)

func main() {
	ctx := context.Background()

	sessionID, err := GetRecentSessionID(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	if sessionID == "" {
		return
	}

	sessionExportModels, err := GetRecentSessionExportModels(ctx, sessionID)
	if err != nil {
		log.Fatalln(err)
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
