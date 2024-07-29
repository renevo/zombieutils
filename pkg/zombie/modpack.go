package zombie

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
		return errors.Wrap(err, "failed to download mod")
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read mod")
	}

	zw, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return errors.Wrap(err, "mod not a valid zip file")
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
				return errors.Wrapf(err, "failed to create directory %q", fullPath)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
			return errors.Wrapf(err, "failed to create directory %q", filepath.Dir(fullPath))
		}

		logrus.Infof("Installing Modpack File To %q", fullPath)
		data, err := readZipFile(file)
		if err != nil {
			return errors.Wrapf(err, "failed to read zipped file %q", file.Name)
		}

		if err := os.WriteFile(fullPath, data, os.ModePerm); err != nil {
			return errors.Wrapf(err, "failed to write zipped file %q", fullPath)
		}
	}

	return nil
}
