package internal

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

type Cache struct {
	logger    *log.Logger
	documents map[string]textDocument
	docsLock  *sync.RWMutex
}

func newCache(l *log.Logger) *Cache {
	return &Cache{
		logger:    l,
		documents: map[string]textDocument{},
		docsLock:  &sync.RWMutex{},
	}
}

func (c *Cache) add(uri, text string) {
	c.docsLock.Lock()
	defer c.docsLock.Unlock()
	doc := newTextDocument(uri, text)
	c.documents[uri] = doc
}

func (c *Cache) lookup(uri string) (textDocument, bool) {
	c.docsLock.RLock()
	defer c.docsLock.RUnlock()
	doc, found := c.documents[uri]
	return doc, found
}

func (c *Cache) update(uri, text string) error {
	c.docsLock.Lock()
	defer c.docsLock.Unlock()

	doc, found := c.documents[uri]
	if !found {
		return fmt.Errorf("document not found: %s", uri)
	}

	c.documents[uri] = doc.update(text)
	return nil
}

func (c *Cache) remove(uri string) {
	c.docsLock.Lock()
	defer c.docsLock.Unlock()
	delete(c.documents, uri)
}

type textDocument struct {
	uri   string
	lines []string
}

func newTextDocument(uri, text string) textDocument {
	return textDocument{
		uri:   uri,
		lines: stringLines(text),
	}
}

func (d textDocument) update(text string) textDocument {
	d.lines = stringLines(text)
	return d
}

func (d textDocument) getLine(linenum int) (string, bool) {
	if len(d.lines) <= linenum {
		return "", false
	}
	return d.lines[linenum], true
}

func stringLines(s string) []string {
	return strings.Split(s, "\n")
}
