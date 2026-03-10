package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	ctx := context.Background()

	sessionID, err := getRecentSessionID(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	sessionExportModel, err := getRecentSessionExportModel(ctx, sessionID)
	if err != nil {
		log.Fatalln(err)
	}

	if sessionExportModel.ProviderID == "" ||
		sessionExportModel.ModelID == "" {
		fmt.Println("No recent session")
		return
	}

	fmt.Printf("Co-Authored-By: opencode %s %s <noreply@opencode.ai>\n",
		sessionExportModel.ProviderID,
		sessionExportModel.ModelID,
	)
}

func getRecentSessionID(ctx context.Context) (string, error) {
	args := []string{"session", "list", "-n", "1", "--format", "json"}
	output, err := exec.CommandContext(ctx, "opencode", args...).Output()
	if err != nil {
		return "", fmt.Errorf("opencode: failed: %w", err)
	}

	var sessionListItems []SessionListItem
	if err := json.Unmarshal(output, &sessionListItems); err != nil {
		return "", fmt.Errorf("json: failed to unmarshal: %w", err)
	}

	if len(sessionListItems) == 0 ||
		sessionListItems[0].ID == "" {
		return "", nil
	}

	return sessionListItems[0].ID, nil
}

func getRecentSessionExportModel(ctx context.Context, sessionID string) (SessionExportModel, error) {
	args := []string{"export", sessionID}
	output, err := exec.CommandContext(ctx, "opencode", args...).Output()
	if err != nil {
		return SessionExportModel{}, fmt.Errorf("session: failed: %w", err)
	}

	var data SessionExportData
	if err := json.Unmarshal(output, &data); err != nil {
		return SessionExportModel{}, fmt.Errorf("json: failed to unmarshal: %w", err)
	}

	for _, message := range data.Messages {
		if message.Info.Model.ProviderID != "" &&
			message.Info.Model.ModelID != "" {
			return message.Info.Model, nil
		}
	}

	return SessionExportModel{}, nil
}
