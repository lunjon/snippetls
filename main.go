package main

import (
	"github.com/lunjon/snippetls/internal"
	"github.com/lunjon/snippetls/internal/snippet"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	logger := log.New(os.Stderr, "", 0)

	configDir, err := initConfig()
	if err != nil {
		log.Fatalf("failed to resolve home folder: %s", err)
	}

	m := snippet.NewSnippetManager(logger)
	err = loadSnippets(logger, configDir, m)

	server := internal.NewServer(logger, m)
	server.Start()
}

func initConfig() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := path.Join(homedir, ".config", "snippetls")
	err = os.MkdirAll(configPath, 0700)
	return configPath, err
}

func loadSnippets(l *log.Logger, configdir string, m *snippet.SnippetManager) error {
	// Reserved for future use
	configFileName := "config.toml"

	filepath.WalkDir(configdir, func(fullpath string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if strings.HasSuffix(fullpath, configFileName) {
			return nil
		}

		basename := path.Base(fullpath)
		parts := strings.Split(basename, ".")
		if len(parts) != 2 {
			return nil
		}

		filetype, ext := parts[0], parts[1]
		if ext != "toml" {
			return nil
		}

		bs, err := os.ReadFile(fullpath)
		if err != nil {
			l.Printf("Error reading file %s: %s", fullpath, err)
			return nil
		}

		if err = m.AddSnippets(filetype, bs); err != nil {
			l.Printf("Error adding snippets for filetype %s: %s", filetype, err)
		}

		return nil
	})
	return nil
}
