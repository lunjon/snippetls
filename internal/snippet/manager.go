package snippet

import (
	"log"
	"path"
	"strings"
)

type SnippetManager struct {
	logger   *log.Logger
	snippets map[string][]Snippet
}

func NewSnippetManager(l *log.Logger) *SnippetManager {
	return &SnippetManager{
		logger:   l,
		snippets: map[string][]Snippet{},
	}
}

func (m *SnippetManager) AddSnippets(filetype string, content []byte) error {
	sn, err := parseTOML(content)
	if err != nil {
		return err
	}

	ft := strings.TrimPrefix(filetype, ".")
	existing := m.snippets[ft]
	m.snippets[ft] = append(existing, sn...)
	return nil
}

func (m *SnippetManager) Search(uri string, prefix string) []Snippet {
	// TODO: support files without extension
	ext := path.Ext(uri)
	if ext == "" {
		return nil
	}

	ext = strings.TrimPrefix(ext, ".")

	snips, ok := m.snippets[ext]
	if !ok {
		return nil
	}

	found := []Snippet{}
	for _, sn := range snips {
		if sn.Matches(prefix) {
			found = append(found, sn)
		}
	}

	return found
}

type Snippet struct {
	trigger string
	snippet string
}

func (sn Snippet) Matches(prefix string) bool {
	return strings.HasPrefix(sn.trigger, prefix)
}

func (sn Snippet) Trigger() string {
	return sn.trigger
}

func (sn Snippet) Expand() string {
	return sn.snippet
}
