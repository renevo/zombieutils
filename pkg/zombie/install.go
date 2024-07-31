package zombie

import (
	"context"
	"encoding/xml"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/renevo/zombieutils/pkg/steam"
)

func (s *Server) Install(ctx context.Context) error {
	// steam install game
	g := &steam.Game{ID: 294420}
	if s.Experimental {
		g.Beta = "latest_experimental"
	}

	installPath, err := filepath.Abs(filepath.FromSlash(s.Path))
	if err != nil {
		return fmt.Errorf("failed to resolve install directory %q: %w", s.Path, err)
	}

	if err := os.MkdirAll(installPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create install directory: %q: %w", installPath, err)
	}

	// mod path cleaning and creations
	modPath := filepath.Join(installPath, "Mods")
	if s.CleanMods {
		slog.Info("Cleaning mod directory")
		modPath := filepath.Join(installPath, "Mods")
		if err := os.RemoveAll(modPath); err != nil {
			return fmt.Errorf("failed to clean mods folder %q: %w", modPath, err)
		}
	}

	if err := os.MkdirAll(modPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create mod directory: %q: %w", modPath, err)
	}

	// install the server (will add TFP mods back in)
	if err := g.Install(s.Steam, installPath); err != nil {
		return fmt.Errorf("failed to install zombie game to %q: %w", s.Path, err)
	}

	// create the save path
	savePath, err := filepath.Abs(filepath.FromSlash(s.SaveFolder))
	if err != nil {
		return fmt.Errorf("failed to resolve save directory %q: %w", s.SaveFolder, err)
	}
	if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create save directory: %q: %w", savePath, err)
	}

	// install admin config
	if err := s.installAdminConfig(savePath); err != nil {
		return fmt.Errorf("failed to install admin.xml to %q: %w", s.SaveFolder, err)
	}

	// install mods from places, these are pretty non-standard in the zip files provided to download, but most are git repositories
	for _, mod := range s.Mods {
		if err := mod.Install(ctx, modPath); err != nil {
			slog.Error("Failed to install mod", "mod", mod.Name, "err", err)
		}
	}

	// install modpacks
	for _, modpack := range s.ModPacks {
		if err := modpack.Install(ctx, modPath); err != nil {
			slog.Error("Failed to install modpack", "modpack", modpack.Name, "err", err)
		}
	}

	return nil
}

func (s *Server) installAdminConfig(savePath string) error {
	// TODO: should probably read the file for the existing webusers, as they have to be created with `createwebuser` command in game
	/*
	   INFO[0028] 2024-07-29T19:49:24 1.606 WRN [Web] [Perms] Ignoring user-entry because of missing 'platform' or 'userid' attribute: <user name="Dante" userid="76561197969618392" pass="password" platform="Stream" crossplatform="EOS" crossuserid="000256d97ada456e870e75495a3ee51e" />
	   INFO[0028] 2024-07-29T19:49:24 1.611 WRN [Web] [Perms] Ignoring apitoken-entry because of missing 'name' attribute: <token token="admin" secret="s3cr3t" permission_level="0" />
	*/
	adminConfig := struct {
		XMLName     xml.Name               `xml:"adminTools"`
		Admins      []ServerAdmin          `xml:"users>user"`
		Permissions []ServerPermission     `xml:"commands>permission"`
		Whitelist   []ServerWhitelistEntry `xml:"whitelist>user"`
		WebUsers    []WebUser              `xml:"webusers>user"`
		WebModules  []WebModule            `xml:"webmodules>module"`
		APITokens   []APIToken             `xml:"apitokens>token"`
	}{
		Admins:      s.Admins,
		Permissions: s.Permissions,
		Whitelist:   s.Whitelist,
		WebUsers:    s.WebUsers,
		WebModules:  s.WebModules,
		APITokens:   s.APITokens,
	}

	data, err := xml.MarshalIndent(&adminConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal %q: %w", s.AdminFileName, err)
	}

	data = append([]byte(xml.Header), data...)

	// gross...
	data = []byte(strings.ReplaceAll(string(data), "></user>", " />"))
	data = []byte(strings.ReplaceAll(string(data), "></permission>", " />"))
	data = []byte(strings.ReplaceAll(string(data), "></module>", " />"))
	data = []byte(strings.ReplaceAll(string(data), "></token>", " />"))

	// have to save this in multiple locations for some reason
	saveDirs := []string{
		filepath.Join(savePath, ".local", "share", "7DaysToDie", "Saves"),
		filepath.Join(savePath, "Saves"),
	}

	for _, saveDir := range saveDirs {
		if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create save directory %q: %w", saveDir, err)
		}

		fileLocation := filepath.Join(saveDir, s.AdminFileName)

		if err := os.WriteFile(fileLocation, data, os.ModePerm); err != nil {
			return fmt.Errorf("failed to write admin file %q: %w", fileLocation, err)
		}
	}

	return nil
}
