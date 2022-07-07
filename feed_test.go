package calendar

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestGetEventItems(t *testing.T) {
	feed := prepareFeed(t, loadTestdata(t, "blog.xml"))
	items := feed.GetEventPosts()
	assert.Assert(t, is.Len(items, 1))
	assert.Equal(t, items[0].Title, "2022年7月のGoイベント一覧")
}

func prepareFeed(t *testing.T, source string) *Feed {
	feed, err := NewFeed(strings.NewReader(source))
	assert.NilError(t, err)
	return feed
}
