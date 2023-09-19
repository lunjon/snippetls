package main

import (
	"github.com/lunjon/snippetls/internal"
	"github.com/lunjon/snippetls/internal/snippet"
	"log"
	"os"
	"path"
)

func main() {
	logger := log.New(os.Stderr, "", 0)

	configDir, err := initConfig()
	if err != nil {
		log.Fatalf("failed to resolve home folder: %s", err)
	}

	m := snippet.NewSnippetManager(logger)
	err = snippet.Load(logger, configDir, m)
	if err != nil {
		log.Fatalf("failed to load snippets: %s", err)
	}

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
