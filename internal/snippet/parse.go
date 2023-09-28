package snippet

import (
	"fmt"
	"strings"

	"github.com/lunjon/gokdl"
)

var emptySnippet = Snippet{}

func parseKDL(bs []byte) ([]Snippet, error) {
	doc, err := gokdl.Parse(bs)
	if err != nil {
		return nil, err
	}

	snippets := []Snippet{}
	for _, node := range doc.Nodes() {
		sn, err := parseNode(node)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, sn)
	}

	return snippets, nil
}

func parseNode(node gokdl.Node) (Snippet, error) {
	var aliases []string
	var argSnippet string
	var childSnippet string

	errReturn := func(msg ...string) (Snippet, error) {
		m := strings.Join(msg, " ")
		return Snippet{}, fmt.Errorf("invalid snippet: %s: %s", node.Name, m)
	}

	if len(node.Args) == 1 {
		val, ok := node.Args[0].Value.(string)
		if !ok {
			return errReturn("argument must be string")
		}

		argSnippet = val
	}

	if len(node.Args) == 0 {
		// No args given => check children
		if len(node.Children) == 0 {
			return errReturn("no snippet body found")
		}

		options := map[string]bool{}
		for _, child := range node.Children {
			key := child.Name
			if options[key] {
				return errReturn("duplicate option:", key)
			}

			if key == "snippet" {
				// Check that there's a string arg, else invalid
				if len(child.Args) != 1 {
					return errReturn("expected snippet key to have one string argument")
				}

				val, ok := child.Args[0].Value.(string)
				if !ok {
					return errReturn("argument must be string")
				}

				childSnippet = val
				options[key] = true
			} else if key == "aliases" {
				if len(child.Children) > 0 {
					return errReturn("invalid aliases: expected no children")
				}

				if len(child.Args) == 0 {
					return errReturn("invalid aliases: no arguments")
				}

				for _, arg := range child.Args {
					val, ok := arg.Value.(string)
					if !ok {
						return errReturn("alias must have type string")
					}

					val = strings.TrimSpace(val)
					if val == "" {
						return errReturn("invalid alias: empty string")
					}

					aliases = append(aliases, val)
				}
			} else {
				return errReturn("unknown option:", key)
			}
		}
	}

	argSnippet = strings.TrimSpace(argSnippet)
	childSnippet = strings.TrimSpace(childSnippet)

	var snippet string
	if argSnippet != "" && childSnippet != "" {
		return emptySnippet, fmt.Errorf("invalid snippet: %s: must not set snippet in arg and child", node.Name)
	} else if argSnippet == "" && childSnippet != "" {
		return emptySnippet, fmt.Errorf("invalid snippet: %s: no snippet body found", node.Name)
	} else if argSnippet != "" {
		snippet = argSnippet
	} else {
		snippet = childSnippet
	}

	return Snippet{
		trigger: node.Name,
		snippet: snippet,
		aliases: aliases,
	}, nil
}
