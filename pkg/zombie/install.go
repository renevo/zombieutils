package zombie

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/renevo/zombieutils/pkg/steam"
	"github.com/sirupsen/logrus"

	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

func (s *Server) Install(ctx context.Context) error {
	// steam install game
	g := &steam.Game{ID: 294420}
	if s.Experimental {
		g.Beta = "latest_experimental"
	}

	if err := g.Install(s.Steam, s.Path); err != nil {
		return errors.Wrapf(err, "failed to install zombie game to %q", s.Path)
	}

	// install alloc mods
	if err := s.installAllocMods(ctx); err != nil {
		return errors.Wrapf(err, "failed to install alloc mods: %q", s.FixesVersion)
	}

	// TODO: other mods from places, these are pretty non-standard in the zip files provided to download, but most are git repositories

	return nil
}

func (s *Server) installAllocMods(ctx context.Context) error {
	client := cleanhttp.DefaultClient()
	fixesURL := fmt.Sprintf("http://illy.bz/fi/7dtd/server_fixes_v%s.tar.gz", strings.ReplaceAll(s.FixesVersion, ".", "_"))

	req, _ := http.NewRequest(http.MethodGet, fixesURL, nil)
	req = req.WithContext(ctx)

	resp, err := client.Do(req)

	if err != nil {
		return errors.Wrapf(err, "failed to download server fixes from %q", fixesURL)
	}
	defer resp.Body.Close()

	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "failed to open new gzip stream from %q", fixesURL)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		cur, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return errors.Wrapf(err, "failed to read tar file from %q", fixesURL)
		}

		if cur.Typeflag != tar.TypeReg {
			continue
		}

		realPath, _ := filepath.Abs(filepath.FromSlash(s.Path))
		realPath = filepath.Join(realPath, cur.Name)

		data, err := io.ReadAll(tr)
		if err != nil {
			return errors.Wrapf(err, "failed to read file %q from tar file %q", cur.Name, fixesURL)
		}

		dir, _ := filepath.Split(realPath)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return errors.Wrapf(err, "failed to create directory %q from tar file %q", dir, fixesURL)
		}

		if err := os.WriteFile(realPath, data, os.ModePerm); err != nil {
			return errors.Wrapf(err, "failed to write file %q from tar file %q", realPath, fixesURL)
		}

		logrus.Debugf("Installed File: %q; Size: %d\n", realPath, cur.Size)
	}

	return nil
}
