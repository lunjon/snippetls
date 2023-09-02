package types

type CompletionOptions struct {
	ResolveProvider bool `json:"resolveProvider"`
}

type TextDocumentSyncKind int
type CompletionItemKind int

const (
	FullTextDocumentSyncKind = 1
	SnippetCompletionItem    = 15
)

type TextDocumentSyncOptions struct {
	OpenClose bool                 `json:"openClose"`
	Change    TextDocumentSyncKind `json:"change"`
}

type CompletionItem struct {
	Kind       int    `json:"kind"`
	Label      string `json:"label"`
	InsertText string `json:"insertText"`
}

type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}

type CompletionItemResolveParams struct {
	InsertText string `json:"insertText"`
	Label      string `json:"label"`
}

type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

type TextDocumentItem struct {
	URI     string `json:"uri"`
	Version int    `json:"version"`
	Text    string `json:"text"`
}

type VersionedTextDocumentIdentifier struct {
	Version int `json:"version"`
}

type TextDocumentContentChangeEvent struct {
	Text string `json:"text"`
}

type TextDocumentDidOpenParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

type TextDocumentDidCloseParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type TextDocumentDidChangeParams struct {
	TextDocument   TextDocumentIdentifier           `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type TextDocumentPosition struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type TextDocumentCompletionParams struct {
	ID           int                    `json:"id"`
	Position     TextDocumentPosition   `json:"position"`
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}
