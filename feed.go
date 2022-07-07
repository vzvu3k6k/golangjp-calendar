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

func (f *Feed) GetEventPosts() []*Post {
	var items []*Post
	for _, item := range f.Items {
		if titlePattern.MatchString(item.Title) {
			items = append(items, &Post{item})
		}
	}
	return items
}
