package zombie

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/renevo/zombieutils/pkg/steam"
	"github.com/sirupsen/logrus"
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

	// install the server (will add TFP mods back in)
	if err := g.Install(s.Steam, installPath); err != nil {
		return errors.Wrapf(err, "failed to install zombie game to %q", s.Path)
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

	// install mods from places, these are pretty non-standard in the zip files provided to download, but most are git repositories
	for _, mod := range s.Mods {
		if err := mod.Install(ctx, modPath); err != nil {
			logrus.Errorf("Failed to install mod %q: %v", mod.Name, err)
		}
	}

	// install modpacks
	for _, modpack := range s.ModPacks {
		if err := modpack.Install(ctx, modPath); err != nil {
			logrus.Errorf("Failed to install modpack %q: %v", modpack.Name, err)
		}
	}

	return nil
}

func (s *Server) installAdminConfig(savePath string) error {
	adminConfig := struct {
		XMLName     xml.Name               `xml:"adminTools"`
		Admins      []ServerAdmin          `xml:"users>user"`
		Permissions []ServerPermission     `xml:"commands>permission"`
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

	saveDir := filepath.Join(savePath, ".local", "share", "7DaysToDie", "Saves")

	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create save directory %q: %w", saveDir, err)
	}

	fileLocation := filepath.Join(saveDir, s.AdminFileName)

	return errors.Wrapf(os.WriteFile(fileLocation, data, os.ModePerm), "failed to write admin file %q", fileLocation)
}
