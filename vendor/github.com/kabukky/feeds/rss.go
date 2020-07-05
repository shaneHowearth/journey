// Package feeds -
package feeds

// rss support
// validation done according to spec here:
//    http://cyber.law.harvard.edu/rss/rss.html

import (
	"encoding/xml"
	"fmt"
	"time"
)

// private wrapper around the RssFeed which gives us the <rss>..</rss> xml
type rssFeedXML struct {
	XMLName    xml.Name `xml:"rss"`
	Version    string   `xml:"version,attr"`
	XMLnsMedia string   `xml:"xmlns:media,attr"`
	XMLnsDc    string   `xml:"xmlns:dc,attr"`
	XMLnsAtom  string   `xml:"xmlns:atom,attr"`
	Channel    *RssFeed
}

// RssAtomLink -
type RssAtomLink struct {
	XMLName xml.Name `xml:"atom:link"`
	Href    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr"`
	Type    string   `xml:"type,attr"`
}

// RssGUID -
type RssGUID struct {
	XMLName     xml.Name `xml:"guid"`
	IsPermaLink bool     `xml:"isPermaLink,attr"`
	Value       string   `xml:",innerxml"`
}

// RssImage -
type RssImage struct {
	XMLName xml.Name `xml:"image"`
	URL     string   `xml:"url"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Width   int      `xml:"width,omitempty"`
	Height  int      `xml:"height,omitempty"`
}

// RssMedia -
type RssMedia struct {
	XMLName xml.Name `xml:"media:content"`
	Medium  string   `xml:"medium,attr"`
	URL     string   `xml:"url,attr"`
	Width   int      `xml:"width,attr,omitempty"`
	Height  int      `xml:"height,attr,omitempty"`
}

// RssTextInput -
type RssTextInput struct {
	XMLName     xml.Name `xml:"textInput"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Name        string   `xml:"name"`
	Link        string   `xml:"link"`
}

// RssFeed -
type RssFeed struct {
	XMLName        xml.Name `xml:"channel"`
	AtomLink       *RssAtomLink
	Title          string `xml:"title"`       // required
	Link           string `xml:"link"`        // required
	Description    string `xml:"description"` // required
	Language       string `xml:"language,omitempty"`
	Copyright      string `xml:"copyright,omitempty"`
	ManagingEditor string `xml:"managingEditor,omitempty"` // Author used
	WebMaster      string `xml:"webMaster,omitempty"`
	PubDate        string `xml:"pubDate,omitempty"`       // created or updated
	LastBuildDate  string `xml:"lastBuildDate,omitempty"` // updated used
	Category       string `xml:"category,omitempty"`
	Generator      string `xml:"generator,omitempty"`
	Docs           string `xml:"docs,omitempty"`
	Cloud          string `xml:"cloud,omitempty"`
	TTL            int    `xml:"ttl,omitempty"`
	Rating         string `xml:"rating,omitempty"`
	SkipHours      string `xml:"skipHours,omitempty"`
	SkipDays       string `xml:"skipDays,omitempty"`
	Image          *RssImage
	TextInput      *RssTextInput
	Items          []*RssItem
}

// RssItem -
type RssItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`       // required
	Link        string   `xml:"link"`        // required
	Description string   `xml:"description"` // required
	Author      string   `xml:"dc:creator,omitempty"`
	Category    string   `xml:"category,omitempty"`
	Comments    string   `xml:"comments,omitempty"`
	Enclosure   *RssEnclosure
	GUID        *RssGUID // ID used
	PubDate     string   `xml:"pubDate,omitempty"` // created or updated
	Source      string   `xml:"source,omitempty"`
	Image       *RssMedia
}

// RssEnclosure -
type RssEnclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	URL     string   `xml:"url,attr"`
	Length  string   `xml:"length,attr"`
	Type    string   `xml:"type,attr"`
}

// Rss -
type Rss struct {
	*Feed
}

// create a new RssImage with a generic Image struct's data
func newRssImage(i *Image) *RssImage {
	image := &RssImage{
		URL:    i.URL,
		Title:  i.Title,
		Link:   i.Link,
		Width:  i.Width,
		Height: i.Height,
	}
	return image
}

// create a new RssImage with a generic Image struct's data
func newRssMediaImage(i *Image) *RssMedia {
	image := &RssMedia{
		Medium: "image",
		URL:    i.URL,
		Width:  i.Width,
		Height: i.Height,
	}
	return image
}

// create a new RssItem with a generic Item struct's data
func newRssItem(i *Item) *RssItem {
	item := &RssItem{
		Title:       i.Title,
		Link:        i.Link.Href,
		Description: i.Description,
		GUID: &RssGUID{
			IsPermaLink: false,
			Value:       i.ID,
		},
		PubDate: anyTimeFormat(time.RFC1123, i.Created, i.Updated),
	}
	if i.Author != nil {
		item.Author = i.Author.Name
	}
	if i.Image != nil {
		item.Image = newRssMediaImage(i.Image)
	}
	return item
}

// RssFeed - create a new RssFeed with a generic Feed struct's data
func (r *Rss) RssFeed() *RssFeed {
	pub := anyTimeFormat(time.RFC1123, r.Created)
	build := anyTimeFormat(time.RFC1123, r.Updated)
	author := ""
	if r.Author != nil {
		author = r.Author.Email
		if len(r.Author.Name) > 0 {
			author = fmt.Sprintf("%s (%s)", r.Author.Email, r.Author.Name)
		}
	}

	channel := &RssFeed{
		AtomLink: &RssAtomLink{
			Href: r.URL,
			Rel:  "self",
			Type: "application/rss+xml",
		},
		Title:          r.Title,
		Link:           r.Link.Href,
		Description:    r.Description,
		ManagingEditor: author,
		PubDate:        pub,
		LastBuildDate:  build,
		Copyright:      r.Copyright,
	}
	if r.Image != nil {
		channel.Image = newRssImage(r.Image)
	}
	for _, i := range r.Items {
		channel.Items = append(channel.Items, newRssItem(i))
	}
	return channel
}

// FeedXML - return an XML-Ready object for an Rss object
func (r *Rss) FeedXML() interface{} {
	// only generate version 2.0 feeds for now
	return r.RssFeed().FeedXML()

}

// FeedXML - return an XML-ready object for an RssFeed object
func (r *RssFeed) FeedXML() interface{} {
	return &rssFeedXML{Version: "2.0", XMLnsMedia: "http://search.yahoo.com/mrss/", XMLnsDc: "http://purl.org/dc/elements/1.1/", XMLnsAtom: "http://www.w3.org/2005/Atom", Channel: r}
}
