package calendar

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

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
	events, err := extractEvents(source, time.Date(2020, time.July, 1, 0, 0, 0, 0, getJST(t)))
	if err != nil {
		t.Error(err)
	}
	if len(events) != 2 {
		t.Errorf("expected 2 events, but got %d", len(events))
	}
	for i, expected := range []event{
		{
			start:    newTime(t, 2022, time.August, 1, 17, 30),
			location: "オンライン",
			title:    "Go 2 リリースパーティー",
			url:      "https://gocon.connpass.com/event/1234/",
		},
		{
			start:    newTime(t, 2022, time.August, 31, 7, 0),
			location: "兵庫県飾磨市",
			title:    "shikamashi.go#1",
			url:      "https://gocon.connpass.com/event/2345/",
		},
	} {
		if events[i] != expected {
			t.Errorf("expected %v, but got %v", expected, events[i])
		}
	}
}

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

func newTime(t *testing.T, year int, month time.Month, day, hour, min int) time.Time {
	t.Helper()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Error(err)
	}
	return time.Date(year, month, day, hour, min, 0, 0, loc)
}

func getJST(t *testing.T) *time.Location {
	t.Helper()
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Error(err)
	}
	return loc
}
