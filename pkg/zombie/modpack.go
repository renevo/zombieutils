package zombie

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

type ServerModPack struct {
	Name string `hcl:"name,label"`
	URL  string `hcl:"url"`
}

func (sm *ServerModPack) Install(ctx context.Context, modPath string) error {
	client := cleanhttp.DefaultClient()

	// google drive workaround junk
	client.CheckRedirect = func(r *http.Request, via []*http.Request) error {
		r.URL.Opaque = r.URL.Path
		return nil
	}

	req, _ := http.NewRequest(http.MethodGet, sm.URL, nil)
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download mod: %w", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read mod: %w", err)
	}

	zw, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("mod not a valid zip file: %w", err)
	}

	// extract all the files
	for _, file := range zw.File {
		// ignore git files
		if strings.Contains(file.Name, "/.git/") {
			continue
		}

		fullPath := filepath.Join(modPath, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %q: %w", fullPath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %q: %w", filepath.Dir(fullPath), err)
		}

		slog.Info("Installing Modpack File", "path", fullPath)
		data, err := readZipFile(file)
		if err != nil {
			return fmt.Errorf("failed to read zipped file %q: %w", file.Name, err)
		}

		if err := os.WriteFile(fullPath, data, os.ModePerm); err != nil {
			return fmt.Errorf("failed to write zipped file %q: %w", fullPath, err)
		}
	}

	return nil
}
