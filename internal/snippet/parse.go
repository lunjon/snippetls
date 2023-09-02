package snippet

import "github.com/BurntSushi/toml"

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
			snippet: val,
		})
	}

	return snippets, nil
}
