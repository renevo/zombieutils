package zombie

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/xml"
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

	installPath, err := filepath.Abs(filepath.FromSlash(s.Path))
	if err != nil {
		return errors.Wrapf(err, "failed to resolve install directory %q", s.Path)
	}

	if err := os.MkdirAll(installPath, os.ModePerm); err != nil {
		return errors.Wrapf(err, "failed to create install directory: %q", installPath)
	}

	if err := g.Install(s.Steam, installPath); err != nil {
		return errors.Wrapf(err, "failed to install zombie game to %q", s.Path)
	}

	// mod path cleaning and creations
	modPath := filepath.Join(installPath, "Mods")
	if s.CleanMods {
		logrus.Infof("Cleaning mod directory")
		modPath := filepath.Join(installPath, "Mods")
		if err := os.RemoveAll(modPath); err != nil {
			return errors.Wrapf(err, "failed to clean mods folder %q", modPath)
		}
	}

	if err := os.MkdirAll(modPath, os.ModePerm); err != nil {
		return errors.Wrapf(err, "failed to create mod directory: %q", modPath)
	}

	// create the save path
	savePath, err := filepath.Abs(filepath.FromSlash(s.SaveFolder))
	if err != nil {
		return errors.Wrapf(err, "failed to resolve save directory %q", s.SaveFolder)
	}
	if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
		return errors.Wrapf(err, "failed to create save directory: %q", savePath)
	}

	// install admin config
	if err := s.installAdminConfig(savePath); err != nil {
		return errors.Wrapf(err, "failed to install admin.xml to %q", s.SaveFolder)
	}

	// install alloc mods - this are a bit "special"
	if len(s.FixesVersion) > 0 {
		if err := s.installAllocMods(ctx, installPath, savePath); err != nil {
			return errors.Wrapf(err, "failed to install alloc mods: %q", s.FixesVersion)
		}
	}

	// install mods from places, these are pretty non-standard in the zip files provided to download, but most are git repositories
	for _, mod := range s.Mods {
		if err := mod.Install(ctx, modPath); err != nil {
			logrus.Errorf("Failed to install mod %q: %v", mod.Name, err)
		}
	}

	return nil
}

func (s *Server) installAdminConfig(savePath string) error {
	adminConfig := struct {
		XMLName     xml.Name               `xml:"adminTools"`
		Admins      []ServerAdmin          `xml:"admins>user"`
		Permissions []ServerPermission     `xml:"permissions>permission"`
		Whitelist   []ServerWhitelistEntry `xml:"whitelist>user"`
	}{
		Admins:      s.Admins,
		Permissions: s.Permissions,
		Whitelist:   s.Whitelist,
	}

	data, err := xml.MarshalIndent(&adminConfig, "", "  ")
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %q", s.AdminFileName)
	}

	data = append([]byte(xml.Header), data...)

	// gross...
	data = []byte(strings.ReplaceAll(string(data), "></user>", " />"))
	data = []byte(strings.ReplaceAll(string(data), "></permission>", " />"))

	return errors.Wrapf(os.WriteFile(filepath.Join(savePath, s.AdminFileName), data, os.ModePerm), "failed to write %q to %q", s.AdminFileName, savePath)
}

func (s *Server) installAllocMods(ctx context.Context, installPath, savePath string) error {
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

		realPath := filepath.Join(installPath, cur.Name)

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

	// now do the web permissions file
	webConfig := struct {
		XMLName     xml.Name              `xml:"webpermissions"`
		AdminTokens []ServerWebToken      `xml:"admintokens>token"`
		Permissions []ServerWebPermission `xml:"permissions>permission"`
	}{
		Permissions: s.WebPermissions,
		AdminTokens: s.WebTokens,
	}

	data, err := xml.MarshalIndent(&webConfig, "", "  ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal webpermissions.xml")
	}

	data = append([]byte(xml.Header), data...)

	// gross...
	data = []byte(strings.ReplaceAll(string(data), "></token>", " />"))
	data = []byte(strings.ReplaceAll(string(data), "></permission>", " />"))

	return errors.Wrapf(os.WriteFile(filepath.Join(savePath, "webpermissions.xml"), data, os.ModePerm), "failed to write webpermissions.xml to %q", savePath)
}
