package snippet

import (
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Load snippets from file system and add to the snippet manager.
// TODO: support "all" snippets that are available no matter file type.
func Load(l *log.Logger, configdir string, m *SnippetManager) error {
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

		if filetype == "global" {
			if err = m.AddGlobalSnippets(bs); err != nil {
				l.Printf("Error adding snippets for filetype %s: %s", filetype, err)
			}
		} else {
			if err = m.AddFiletypeSnippets(filetype, bs); err != nil {
				l.Printf("Error adding snippets for filetype %s: %s", filetype, err)
			}
		}

		return nil
	})

	return nil
}