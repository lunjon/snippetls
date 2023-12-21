# SnippetLS

A pure Go, although rudimentary, implementation of a language server for the [language server protocol](https://microsoft.github.io/language-server-protocol/) (LSP) that only provides snippet support.

## Installation

Using `go`:
```sh
go get github.com/lunjon/snippetls
```

## Usage

This language server was only developed in order to inject snippets into an editor with LSP snippet support.

Snippets are loaded per language as a file name `<ext>.kdl` put into the `~/.config/snippetls/`.
So, for instance, to load snippets for the rust programming language you would create such a file named `rs.kdl`.

### Snippet definitions
The snippets are defined in [kdl](https://kdl.dev/) files, since they are typically very simple.

#### Example
```kdl
# go.kdl

// You can use raw strings in order to define snippets
// containing " easily:
log r#"log.Printf("$1")"#

iferr "if err != nil {
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
_Global_ are snippets that are included for all file types.
They can be defined in a file called `global.kdl`. That is,
create a file called `global.kdl` and these snippets will be
available for all file types.

### Configuration
You can create a file called `config.kdl` (in the same directory as the snippets).

This is used to configure `snippetls`:

```kdl
// This section can also be called "extends"
valid-for {
    ex "exs"            // Include all snippets in ex.kdl for exs files as well
    ts "tsx" "js" "jsx" // The ts.kdl snippets are valid for JS and react files
    md "txt"            // The markdown snippets are also valid for .txt files
}
```
