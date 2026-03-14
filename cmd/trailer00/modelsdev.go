package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bytedance/sonic"
)

const (
	modelsDevURL       = "https://models.dev/api.json"
	modelsDevCacheFile = "models-dev.json"
)

var ErrHTTPStatusNotOK = errors.New("http status code not ok")

func LoadModelsDevData(ctx context.Context) (ModelsDevData, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("os: failed to get user cache dir: %w", err)
	}

	cacheFilePath := filepath.Join(userCacheDir, "trailer00", modelsDevCacheFile)

	if _, err := os.Stat(cacheFilePath); err == nil {
		data, err := os.ReadFile(cacheFilePath)
		if err == nil {
			var api ModelsDevData
			if err := sonic.Unmarshal(data, &api); err == nil {
				return api, nil
			}

			// Invalid cache
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, modelsDevURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("http: failed to new request: %w", err)
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http: failed to get models.dev: %w", err)
	}
	defer func() {
		if err := rsp.Body.Close(); err != nil {
			slog.Error("rsp.Body.Close", "error", err)
			return
		}
	}()

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http: unexpected status code: %d: %w", rsp.StatusCode, ErrHTTPStatusNotOK)
	}

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("io: failed to read all: %w", err)
	}

	var data ModelsDevData
	if err := sonic.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("json: failed to unmarshal: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(cacheFilePath), 0o700); err != nil {
		return nil, fmt.Errorf("os: failed to mkdir: %w", err)
	}

	if err := os.WriteFile(cacheFilePath, body, 0o600); err != nil {
		return nil, fmt.Errorf("os: failed to write file: %w", err)
	}

	return data, nil
}
