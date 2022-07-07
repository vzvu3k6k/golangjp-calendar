package calendar

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func loadTestdata(t *testing.T, filename string) string {
	t.Helper()
	f, err := os.Open(filepath.Join("testdata", filename))
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() { f.Close() })

	content, err := io.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	return string(content)
}
