package snippet

import (
	"strings"

	"github.com/BurntSushi/toml"
)

type snippetObject struct {
	Trigger string
}

func parseTOML(bs []byte) ([]Snippet, error) {
	content := map[string]string{}
	_, err := toml.Decode(string(bs), &content)
	if err != nil {
		return nil, err
	}

	snippets := []Snippet{}
	for key, val := range content {
		snippets = append(snippets, Snippet{
			trigger: key,
			snippet: strings.TrimSpace(val),
		})
	}

	return snippets, nil
}
