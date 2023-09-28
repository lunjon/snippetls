# SnippetLS

A rudimentary implementation of a language server for the [language server protocol](https://microsoft.github.io/language-server-protocol/) (LSP).

## Installation

Using `go`:
```sh
go get github.com/lunjon/snippetls
```

## Usage

This language server was only developed in order to inject snippets into an editor with LSP snippet support.

Snippets are loaded per language as a file name `<ext>.kdl` put into the `~/.config/snippetls/`.
So, for instance, to load snippets for the rust programming language you would creates such a file named `rs.kdl`.

### Snippet definitions
The snippets are definied in [kdl](https://kdl.dev/) files, as they are very simple.

A snippet is created using with `key = "<snippet>"`, for instance:

```kdl
# go.kdl

iferr "
if err != nil {
    $1 
}"

// The node name is the trigger
map {
    // Snippet can be specified using the "snippet" node name
    snippet "map[string]$1"
    // Defined one or more aliases for the trigger
    aliases "m" "mp"
}
```

### Global snippets
Global snippets are support via the `global.kdl` file. That is,
create a file called `global.kdl` and these snippets will be available no matter
the filetype.
