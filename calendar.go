package calendar

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	ics "github.com/arran4/golang-ical"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
)

func Run(out io.Writer, args []string) error {
	feedURL := args[0]

	feed, err := http.Get(feedURL)
	if err != nil {
		return err
	}
	defer feed.Body.Close()

	posts, err := extractEventPosts(feed.Body)
	if err != nil {
		return err
	}

	var events []event
	for _, post := range posts[:1] {
		baseDate, err := extractBaseDate(post.Title)
		if err != nil {
			return err
		}

		_events, err := extractEvents(post.Content, baseDate)
		if err != nil {
			return err
		}
		events = append(events, _events...)
	}

	calendar := buildCalendar(events)
	fmt.Fprintln(out, calendar)

	return nil
}

func buildCalendar(events []event) string {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	for _, event := range events {
		e := cal.AddEvent(fmt.Sprintf("id@domain-%d", event.start.Unix())) // TODO: まともなIDを生成する
		e.SetStartAt(event.start)
		e.SetEndAt(event.end)
		e.SetSummary(event.title)
		e.SetLocation(event.location)
		e.SetURL(event.url)
	}
	return cal.Serialize()
}

var titlePattern = regexp.MustCompile(`^\d{4}年\d{1,2}月のGoイベント一覧$`)

func extractEventPosts(source io.Reader) ([]*gofeed.Item, error) {
	feed, err := gofeed.NewParser().Parse(source)
	if err != nil {
		return nil, err
	}

	var items []*gofeed.Item
	for _, item := range feed.Items {
		if titlePattern.MatchString(item.Title) {
			items = append(items, item)
		}
	}
	return items, nil
}

func extractEvents(content string, baseDate time.Time) ([]event, error) {
	doc, err := htmlquery.Parse(strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	nodes, err := htmlquery.QueryAll(doc, "//div/ul/li")
	if err != nil {
		return nil, err
	}

	var events []event
	for _, node := range nodes {
		e, err := parseEvent(node, baseDate)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

var eventTextRegexp = regexp.MustCompile(`(\d{1,2}/\d{1,2})\([日月火水木金土]\)\p{Zs}*(\d{1,2}:\d{2})〜(\d{1,2}:\d{2})\p{Zs}*\[(.+?)\]`)

func parseEvent(node *html.Node, baseDate time.Time) (event, error) {
	var e event

	link := htmlquery.FindOne(node, "./a")
	title, url, err := parseTitleAndURL(link)
	if err != nil {
		return event{}, err
	}
	e.title = title
	e.url = url

	text := htmlquery.InnerText(node)
	matches := eventTextRegexp.FindStringSubmatch(text)
	e.location = matches[4]

	start, end, err := parseStartAndEnd(text, baseDate)
	if err != nil {
		return event{}, err
	}
	e.start = start
	e.end = end

	return e, nil
}

func parseTitleAndURL(link *html.Node) (string, string, error) {
	title := htmlquery.InnerText(link)
	url := htmlquery.SelectAttr(link, "href")
	return title, url, nil
}

func parseStartAndEnd(text string, baseDate time.Time) (time.Time, time.Time, error) {
	matches := eventTextRegexp.FindStringSubmatch(text)
	date, err := parsePartialDate("1/2", matches[1], baseDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	start, err := parsePartialTime("15:04", matches[2], date)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err := parsePartialTime("15:04", matches[3], date)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return start, end, nil
}

func parsePartialDate(layout, value string, base time.Time) (time.Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(base.Year(), t.Month(), t.Day(), 0, 0, 0, 0, base.Location()), nil
}

func parsePartialTime(layout, value string, base time.Time) (time.Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(base.Year(), base.Month(), base.Day(), t.Hour(), t.Minute(), 0, 0, base.Location()), nil
}

func extractBaseDate(postTitle string) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return time.Time{}, err
	}
	return time.ParseInLocation("2006年1月のGoイベント一覧", postTitle, loc)
}
