package calendar

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestRun(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, loadTestdata(t, "blog_full.xml"))
		}),
	)
	defer ts.Close()

	var out bytes.Buffer
	err := Run(&out, []string{ts.URL})
	// os.WriteFile("testdata/blog_full.ical", out.Bytes(), 0644)
	assert.NilError(t, err)
	assert.Equal(t, out.String(), loadTestdata(t, "blog_full.ical"))
}

func TestGetEventItems(t *testing.T) {
	source := loadTestdata(t, "blog.xml")
	items, err := getEventPosts(strings.NewReader(source))
	assert.NilError(t, err)
	assert.Assert(t, is.Len(items, 1))
	assert.Equal(t, items[0].Title, "2022年7月のGoイベント一覧")
}

func TestExtractBaseDate(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		title := "2022年7月のGoイベント一覧"
		got, err := extractBaseDate(title)
		assert.NilError(t, err)
		want := newTime(t, 2022, time.July, 1, 0, 0)
		assert.DeepEqual(t, got, want)
	})
	t.Run("invalid", func(t *testing.T) {
		title := "invalid title"
		_, err := extractBaseDate(title)
		assert.ErrorContains(t, err, "cannot parse")
	})
}

func TestExtractEvents(t *testing.T) {
	source := loadTestdata(t, "content.html")
	got, err := extractEvents(source, newTime(t, 2022, time.August, 1, 0, 0))
	assert.NilError(t, err)

	want := []event{
		{
			start:    newTime(t, 2022, time.August, 1, 17, 0),
			end:      newTime(t, 2022, time.August, 1, 20, 0),
			location: "オンライン",
			title:    "Go 2 リリースパーティー",
			url:      "https://gocon.connpass.com/event/1234/",
		},
		{
			start:    newTime(t, 2022, time.August, 31, 7, 0),
			end:      newTime(t, 2022, time.August, 31, 7, 30),
			location: "兵庫県飾磨市",
			title:    "shikamashi.go#1",
			url:      "https://gocon.connpass.com/event/2345/",
		},
	}
	assert.Check(t, is.DeepEqual(got, want, cmp.AllowUnexported(event{})))
}

func TestParsePartialTime(t *testing.T) {
	base := newTime(t, 2022, time.July, 7, 0, 0)
	t.Run("normal", func(t *testing.T) {
		got, err := parsePartialTime("15:04", "9:00", base)
		assert.NilError(t, err)
		want := newTime(t, 2022, time.July, 7, 9, 0)
		assert.DeepEqual(t, got, want)
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := parsePartialTime("15:04", "invalid", base)
		assert.ErrorContains(t, err, "cannot parse")
	})
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
