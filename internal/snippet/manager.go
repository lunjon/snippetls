package snippet

import (
	"log"
	"path"
	"strings"
)

// TODO: support files without extension
type SnippetManager struct {
	logger   *log.Logger
	globals  []Snippet
	snippets map[string][]Snippet
}

func NewSnippetManager(l *log.Logger) *SnippetManager {
	return &SnippetManager{
		logger:   l,
		globals:  []Snippet{},
		snippets: map[string][]Snippet{},
	}
}

func (m *SnippetManager) AddGlobalSnippets(content []byte) error {
	sn, err := parseTOML(content)
	if err != nil {
		return err
	}

	m.globals = append(m.globals, sn...)
	return nil
}

func (m *SnippetManager) AddFiletypeSnippets(filetype string, content []byte) error {
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
	ext := path.Ext(uri)
	if ext == "" {
		return nil
	}

	ext = strings.TrimPrefix(ext, ".")

	snips := m.snippets[ext]
	snips = append(snips, m.globals...)

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
