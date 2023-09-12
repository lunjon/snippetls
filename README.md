# SnippetLS

A rudimentary implementation of a language server for the [language server protocol](https://microsoft.github.io/language-server-protocol/) (LSP).

## Installation

Using `go`:
```sh
go get github.com/lunjon/snippetls
```

## Usage

This language server was only developed in order to inject snippets into an editor with LSP snippet support.

Snippets are loaded per language as a file name `<ext>.toml` put into the `~/.config/snippetls/`.
So, for instance, to load snippets for the rust programming language you would creates such a file named `rs.toml`.

### Snippet definitions
The snippets are definied in TOML files as they are very simple.

A snippet is created using with `key = "<snippet>"`, for instance:

```toml
# go.toml

iferr = """
if err != nil {
    $1 
}"""

printf = "fmt.Printf(\"$1\")"
```

### Global snippets
Global snippets are support via the `global.toml` file. That is,
create a file called `global.toml` and these snippets will be available no matter
the filetype.
