package types

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#serverCapabilities
type ServerCapabilities struct {
	CompletionProvider CompletionOptions       `json:"completionProvider"`
	TextDocumentSync   TextDocumentSyncOptions `json:"textDocumentSync"`
}

type ServerInfo struct {
	Name string `json:"name"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

func Initialize() InitializeResult {
	return InitializeResult{
		Capabilities: ServerCapabilities{
			CompletionProvider: CompletionOptions{ResolveProvider: true},
			TextDocumentSync: TextDocumentSyncOptions{
				OpenClose: true,
				Change:    FullTextDocumentSyncKind,
			},
		},
		ServerInfo: ServerInfo{
			Name: "snippetls",
		},
	}
}
