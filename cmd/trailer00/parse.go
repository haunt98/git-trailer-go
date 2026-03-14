package main

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"slices"

	"github.com/bytedance/sonic"
)

func GetRecentSessionID(ctx context.Context) (string, error) {
	args := []string{"session", "list", "-n", "1", "--format", "json"}
	output, err := exec.CommandContext(ctx, "opencode", args...).Output()
	if err != nil {
		return "", fmt.Errorf("opencode: failed: %w", err)
	}

	var sessionListItems []SessionListItem
	if err := sonic.Unmarshal(output, &sessionListItems); err != nil {
		return "", fmt.Errorf("json: failed to unmarshal: %w", err)
	}

	if len(sessionListItems) == 0 ||
		sessionListItems[0].ID == "" {
		return "", nil
	}

	return sessionListItems[0].ID, nil
}

func GetSessionExportModels(ctx context.Context, sessionID string) ([]SessionExportModel, error) {
	tempFile, err := os.CreateTemp("", "opencode-export-*.json")
	if err != nil {
		return nil, fmt.Errorf("os: failed to create temp: %w", err)
	}
	defer func() {
		if err := tempFile.Close(); err != nil {
			slog.Error("tempFile.Close", "error", err)
			return
		}

		if err := os.Remove(tempFile.Name()); err != nil {
			slog.Error("os.Remove", "error", err)
			return
		}
	}()

	args := []string{"export", sessionID}
	cmd := exec.CommandContext(ctx, "opencode", args...)
	cmd.Stdout = tempFile

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("opencode: failed: %w", err)
	}

	output, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("os: failed to read file: %w", err)
	}

	var data SessionExportData
	if err := sonic.Unmarshal(output, &data); err != nil {
		return nil, fmt.Errorf("json: failed to unmarshal: %w", err)
	}

	models := make([]SessionExportModel, 0, len(data.Messages))
	for _, message := range data.Messages {
		if message.Info.Model.ProviderID != "" &&
			message.Info.Model.ModelID != "" {
			models = append(models, message.Info.Model)
		}
	}

	slices.SortFunc(models, func(a, b SessionExportModel) int {
		return cmp.Or(
			cmp.Compare(a.ProviderID, b.ProviderID),
			cmp.Compare(a.ModelID, b.ModelID),
		)
	})

	models = slices.CompactFunc(models, func(a, b SessionExportModel) bool {
		return a.ProviderID == b.ProviderID &&
			a.ModelID == b.ModelID
	})

	return models, nil
}
