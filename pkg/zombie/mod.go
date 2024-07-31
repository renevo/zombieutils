package zombie

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

type ServerMod struct {
	Name   string `hcl:"name,label"`
	URL    string `hcl:"url"`
	Filter string `hcl:"path_filter,optional"`
}

func (sm *ServerMod) Install(ctx context.Context, modPath string) error {
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

	trimDir := ""
	hasModInfo := false

	// first pass, find the mod info
	for _, file := range zw.File {
		if len(sm.Filter) > 0 && !strings.Contains(file.Name, sm.Filter) {
			continue
		}

		if dir, fileName := filepath.Split(file.Name); fileName == "ModInfo.xml" {
			trimDir = dir
			hasModInfo = true
		}
	}

	if !hasModInfo {
		return errors.New("mod does not contain a ModInfo.xml")
	}

	// we want whatever the ModInfo.xml root path is, to be the root path for the installation
	for _, file := range zw.File {
		if len(sm.Filter) > 0 && !strings.Contains(file.Name, sm.Filter) {
			continue
		}

		fileName := strings.TrimPrefix(file.Name, trimDir)
		if fileName == "" {
			continue
		}

		fullPath := filepath.Join(modPath, sm.Name, fileName)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %q: %w", fullPath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %q: %w", filepath.Dir(fullPath), err)
		}

		slog.Info("Installing Mod File", "path", fullPath)
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

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}
