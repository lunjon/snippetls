package internal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lunjon/snippetls/internal/snippet"
	"github.com/lunjon/snippetls/internal/types"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Server struct {
	logger      *log.Logger
	cache       *Cache
	snippets    *snippet.SnippetManager
	sendChannel chan any // sendChannel is used to synchronize writes back to the client
	writer      io.WriteCloser
	reader      *bufio.Reader
}

func NewServer(logger *log.Logger, sm *snippet.SnippetManager) *Server {
	return &Server{
		logger:      logger,
		cache:       newCache(logger),
		snippets:    sm,
		sendChannel: make(chan any, 5),
		writer:      os.Stdout,
		reader:      bufio.NewReader(os.Stdin),
	}
}

func (s *Server) Start() {
	s.logger.Println("Starting SnippetLS")

	done := make(chan bool)
	go s.handleSend(done)

	for {
		req, err := s.readMessage()
		if err != nil {
			s.logger.Printf("Error reading message from stream: %s", err)
			continue
		}

		if req.Method == "exit" {
			s.logger.Println("Exit request received: exiting")
			done <- true
			return
		}

		go func(s *Server) {
			err = s.handleRequest(req)
			if err != nil {
				s.logger.Printf("Error handling request: %s", err)
			}
		}(s)
	}
}

func (s *Server) readMessage() (types.RequestMessage, error) {
	// Read first line: Content-Length: N\r\n
	contentLine, err := s.reader.ReadBytes('\n')
	if err != nil {
		return types.RequestMessage{}, err
	}

	index := bytes.IndexAny(contentLine, " ")
	if index == -1 {
		return types.RequestMessage{}, fmt.Errorf("failed to read from stream")
	}

	numstring := bytes.TrimSpace(contentLine[index:])
	contentLength, err := strconv.Atoi(string(numstring))
	if err != nil {
		return types.RequestMessage{}, fmt.Errorf("failed to read from stream: %w", err)
	}

	// Skip next line (empty)
	_, _ = s.reader.ReadBytes('\n')

	// Read rest of the message
	buf := make([]byte, contentLength)
	_, err = io.ReadFull(s.reader, buf)
	if err != nil {
		return types.RequestMessage{}, fmt.Errorf("failed to read from stream: %w", err)
	}

	return s.parseMessage(buf)
}

func (s *Server) sendRequest(method string, params any) {
	m := types.NewRequest(method, params)
	s.sendChannel <- m
}

func (s *Server) sendResponse(id int, result any) {
	m := types.ResponseMessage{
		JsonRpc: "2.0",
		ID:      id,
		Result:  result,
	}
	s.sendChannel <- m
}

// This is run in its own go routine.
func (s *Server) handleSend(done chan bool) {
	for {
		select {
		case <-done:
			s.logger.Printf("Received done: returning from handleSend")
			return
		case data := <-s.sendChannel:
			bs, err := json.Marshal(data)
			if err != nil {
				s.logger.Printf("Error trying to marshal data: %s", err)
				continue
			}

			buf := bytes.NewBuffer(nil)

			header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(bs))
			_, err = buf.WriteString(header)
			if err != nil {
				s.logger.Printf("Error when creating header: %s", err)
				continue
			}

			_, err = buf.Write(bs)
			if err != nil {
				s.logger.Printf("Error when creating body: %s", err)
				continue
			}

			mbytes := buf.Bytes()

			_, err = s.writer.Write(mbytes)
			if err != nil {
				s.logger.Printf("Error when writing to stdout: %s", err)
			}
		}
	}
}

// Each request from the client is handled with this method,
// and is started a new go routine for each one.
func (s *Server) handleRequest(m types.RequestMessage) error {
	switch strings.TrimSpace(m.Method) {
	case "shutdown":
		s.logger.Println("Shutdown request received")
		s.sendResponse(m.ID, nil)

	case "initialize":
		initialize := types.Initialize()
		s.sendResponse(m.ID, initialize)
		s.sendRequest("initialized", map[string]any{})

	case "textDocument/didOpen":
		params, err := types.ParamsAs[types.TextDocumentDidOpenParams](m.Params)
		if err != nil {
			return err
		}
		s.cache.add(params.TextDocument.URI, params.TextDocument.Text)

	case "textDocument/didClose":
		params, err := types.ParamsAs[types.TextDocumentDidCloseParams](m.Params)
		if err != nil {
			return err
		}
		s.cache.remove(params.TextDocument.URI)

	case "textDocument/didChange":
		params, err := types.ParamsAs[types.TextDocumentDidChangeParams](m.Params)
		if err != nil {
			return err
		}

		for _, d := range params.ContentChanges {
			s.cache.update(params.TextDocument.URI, d.Text)
		}

	case "textDocument/completion":
		params, err := types.ParamsAs[types.TextDocumentCompletionParams](m.Params)
		if err != nil {
			return err
		}

		doc, found := s.cache.lookup(params.TextDocument.URI)
		if !found {
			return nil
		}

		line, ok := doc.getLine(params.Position.Line)
		if !ok {
			s.sendResponse(m.ID, nil)
			return nil
		}

		line = strings.TrimSpace(line)
		snippets := s.snippets.Search(params.TextDocument.URI, line)

		items := []types.CompletionItem{}
		for _, sn := range snippets {
			items = append(items, types.CompletionItem{
				Kind:       types.SnippetCompletionItem,
				Label:      sn.Trigger(),
				InsertText: sn.Expand(),
			})
		}

		result := types.CompletionList{
			IsIncomplete: len(items) == 0,
			Items:        items,
		}

		s.sendResponse(m.ID, result)

	case "completionItem/resolve":
		params, err := types.ParamsAs[types.CompletionItemResolveParams](m.Params)
		if err != nil {
			return err
		}

		result := types.CompletionItem{
			Kind:       15,
			Label:      params.Label,
			InsertText: params.InsertText,
		}

		s.sendResponse(m.ID, result)
	}

	return nil
}

func (s *Server) parseMessage(bs []byte) (types.RequestMessage, error) {
	var m types.RequestMessage
	err := json.Unmarshal(bs, &m)
	return m, err
}
