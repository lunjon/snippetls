package snippet

import (
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

var logger = log.New(io.Discard, "", 0)

func TestLoad(t *testing.T) {
	_, err := Load(logger, "testdata")
	require.NoError(t, err)
}
