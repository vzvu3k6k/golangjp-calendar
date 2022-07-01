package calendar

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func loadTestdata(t *testing.T, filename string) string {
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

func TestGetEventItems(t *testing.T) {
	source := loadTestdata(t, "blog.xml")
	items, err := getEventPosts(strings.NewReader(source))
	if err != nil {
		t.Error(err)
	}

	if len(items) != 1 {
		t.Errorf("expected 1 item, but got %d", len(items))
	}
	if items[0].Title != "2022年7月のGoイベント一覧" {
		t.Errorf("expected \"2022年7月のGoイベント一覧\", but got %s", items[0].Title)
	}
}

func TestExtractEvents(t *testing.T) {
	source := loadTestdata(t, "content.html")
	events := extractEvents(source)
	if len(events) != 2 {
		t.Errorf("expected 2 events, but got %d", len(events))
	}
}
