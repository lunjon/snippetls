package snippet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseSimpleNodes(t *testing.T) {
	tests := []struct {
		testname        string
		bs              string
		expectedTrigger string
		expectedContent string
	}{
		{"simple 1", `todo "TODO"`, "todo", "TODO"},
	}

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			snippets, err := parseKDL([]byte(test.bs))
			require.NoError(t, err)
			require.Len(t, snippets, 1)

			sn := snippets[0]
			require.Equal(t, test.expectedTrigger, sn.trigger)
			require.Equal(t, test.expectedContent, sn.snippet)
		})
	}
}

func TestParseSimpleNodesInvalid(t *testing.T) {
	tests := []struct {
		testname string
		bs       string
	}{
		{"invalid type - int", `todo 123`},
		{"invalid value - null", `todo null`},
		{"invalid value - empty string", `todo ""`},
		{"invalid value - whitespace only", `todo "  "`},
	}

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			snippets, err := parseKDL([]byte(test.bs))
			require.Error(t, err)
			require.Nil(t, snippets)
		})
	}
}

func TestParseWithChildren(t *testing.T) {
	tests := []struct {
		testname        string
		bs              string
		expectedTrigger string
		expectedContent string
		expectedAliases []string
	}{
		{
			"snippet in arg",
			`todo "TODO" {}`,
			"todo",
			"TODO",
			nil,
		},
		{
			"snippet in children",
			`todo {
				snippet "TODO"
			}`,
			"todo",
			"TODO",
			nil,
		},
		{
			"aliases",
			`todo {
				snippet "TODO"
				aliases "td" "tod"
			}`,
			"todo",
			"TODO",
			[]string{"td", "tod"},
		},
	}

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			snippets, err := parseKDL([]byte(test.bs))
			require.NoError(t, err)
			require.Len(t, snippets, 1)

			sn := snippets[0]
			require.Equal(t, test.expectedTrigger, sn.trigger)
			require.Equal(t, test.expectedContent, sn.snippet)
			require.Equal(t, test.expectedAliases, sn.aliases)
		})
	}
}

func TestParseWithChildrenInvalid(t *testing.T) {
	tests := []struct {
		testname string
		bs       string
	}{
		// {
		// 	"no snippet",
		// 	`todo {}`,
		// },
		// {
		// 	"unknown key",
		// 	`todo {
		// 		invalidKey
		// 	}`,
		// },
		// {
		// 	"invalid type (arg)",
		// 	`todo 123 {
		// 	}`,
		// },
		// {
		// 	"invalid type (child)",
		// 	`todo {
		// 		snippet null
		// 	}`,
		// },
		// {
		// 	"missing snippet arg",
		// 	`todo {
		// 		snippet
		// 	}`,
		// },
		{
			"invalid alias - missing arguments",
			`todo "Todo" {
				aliases
			}`,
		},
		// {
		// 	"invalid alias - children",
		// 	`todo "Todo" {
		// 		aliases { hehe; }
		// 	}`,
		// },
		// {
		// 	"invalid alias - invalid type",
		// 	`todo "Todo" {
		// 		aliases 123 null
		// 	}`,
		// },
	}

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			snippets, err := parseKDL([]byte(test.bs))
			require.Error(t, err)
			require.Nil(t, snippets)
		})
	}
}
