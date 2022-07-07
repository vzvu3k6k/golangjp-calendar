package calendar

import (
	"io"
	"regexp"

	"github.com/mmcdole/gofeed"
)

type Feed struct {
	*gofeed.Feed
}

func NewFeed(source io.Reader) (*Feed, error) {
	feed, err := gofeed.NewParser().Parse(source)
	if err != nil {
		return nil, err
	}
	return &Feed{feed}, nil
}

var titlePattern = regexp.MustCompile(`^\d{4}年\d{1,2}月のGoイベント一覧$`)

func (f *Feed) GetEventPosts() []*gofeed.Item {
	var items []*gofeed.Item
	for _, item := range f.Items {
		if titlePattern.MatchString(item.Title) {
			items = append(items, item)
		}
	}
	return items
}
