package snippet

import (
	"bytes"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/lunjon/gokdl"
)

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

func (m *SnippetManager) AddConfig(bs []byte) error {
	doc, err := gokdl.Parse(bytes.NewReader(bs))
	if err != nil {
		return err
	}

	for _, node := range doc.Nodes() {
		switch node.Name {
		case "valid-for", "extends":
			if len(node.Args) > 0 {
				return fmt.Errorf("%s takes no arguments", node.Name)
			}
			if len(node.Props) > 0 {
				return fmt.Errorf("%s has no properties", node.Name)
			}

			for _, child := range node.Children {
				snips, ok := m.snippets[child.Name]
				if !ok {
					continue
				}

				for _, arg := range child.Args {
					if v, ok := arg.Value.(string); ok {
						m.snippets[v] = snips
					}
				}
			}
		default:
			return fmt.Errorf("unknown configuration option: %s", node.Name)
		}
	}

	return nil
}

func (m *SnippetManager) AddGlobalSnippets(content []byte) error {
	sn, err := parseKDL(content)
	if err != nil {
		return err
	}

	m.globals = append(m.globals, sn...)
	return nil
}

func (m *SnippetManager) AddFiletypeSnippets(filetype string, content []byte) error {
	sn, err := parseKDL(content)
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
	aliases []string
}

func (sn Snippet) Matches(prefix string) bool {
	if strings.HasPrefix(sn.trigger, prefix) {
		return true
	}

	for _, alias := range sn.aliases {
		if strings.HasPrefix(alias, prefix) {
			return true
		}
	}

	return false
}

func (sn Snippet) Trigger() string {
	return sn.trigger
}

func (sn Snippet) Expand() string {
	return sn.snippet
}
