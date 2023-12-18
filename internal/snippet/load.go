package snippet

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Load snippets from file system and add to the snippet manager.
func Load(l *log.Logger, configdir string) (*SnippetManager, error) {
	// Reserved for future use
	configFileName := "config.kdl"
	m := NewSnippetManager(l)

	filepath.WalkDir(configdir, func(fullpath string, d fs.DirEntry, err error) error {
		if err != nil {
			l.Printf("Unexpected error: %s", err)
			return err
		}

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
		if ext != "kdl" {
			return nil
		}

		bs, err := os.ReadFile(fullpath)
		if err != nil {
			l.Printf("Error reading file %s: %s", fullpath, err)
			return nil
		}

		if filetype == "global" {
			if err = m.AddGlobalSnippets(bs); err != nil {
				l.Printf("Error adding global snippets: %s", err)
			}
		} else {
			if err = m.AddFiletypeSnippets(filetype, bs); err != nil {
				l.Printf("Error adding snippets for filetype %s: %s", filetype, err)
			}
		}

		return nil
	})

	configFilepath := path.Join(configdir, configFileName)
	stat, err := os.Stat(configFilepath)
	if err == nil && !stat.IsDir() {
		bs, err := os.ReadFile(configFilepath)
		if err != nil {
			return m, fmt.Errorf("error reading config file: %s", err)
		}

		m.AddConfig(bs)
	}

	return m, nil
}
