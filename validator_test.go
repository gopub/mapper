package goparam

import (
	"strings"
	"testing"
)

type Image struct {
	Width  int    `param:"min=100,max=800"`
	Height int    `param:"min=100,max=800"`
	Link   string `param:"pattern=url"`
}

type Topic struct {
	Title      string   `param:"min=2,max=30"`
	CoverImage *Image   `param:"optional"`
	MoreImages []*Image `param:"optional"`
}

func TestValidate(t *testing.T) {
	topic := &Topic{
		Title: "a",
	}
	err := Validate(topic)
	if err == nil {
		t.FailNow()
	}

	topic.Title = strings.Repeat("a", 31)
	err = Validate(topic)
	if err == nil {
		t.FailNow()
	}

	topic.Title = strings.Repeat("a", 30)
	err = Validate(topic)
	if err != nil {
		t.FailNow()
	}

	topic.Title = strings.Repeat("a", 2)
	err = Validate(topic)
	if err != nil {
		t.FailNow()
	}

	topic.CoverImage = &Image{
		Width:  100,
		Height: 0,
		Link:   "https://www.image.com",
	}
	err = Validate(topic)
	if err == nil {
		t.FailNow()
	}

	topic.CoverImage = &Image{
		Width:  100,
		Height: 900,
		Link:   "https://www.image.com",
	}
	err = Validate(topic)
	if err == nil {
		t.FailNow()
	}

	topic.CoverImage = &Image{
		Width:  100,
		Height: 800,
		Link:   "https://www.image.com",
	}
	err = Validate(topic)
	if err != nil {
		t.FailNow()
	}

	topic.MoreImages = []*Image{{
		Width:  0,
		Height: 800,
		Link:   "https://www.image.com",
	}}
	err = Validate(topic)
	if err == nil {
		t.FailNow()
	}

	topic.MoreImages = []*Image{{
		Width:  100,
		Height: 800,
		Link:   ":",
	}}
	err = Validate(topic)
	if err == nil {
		t.FailNow()
	}

	topic.MoreImages = []*Image{{
		Width:  100,
		Height: 800,
		Link:   "https://www.image.com",
	}}
	err = Validate(topic)
	if err != nil {
		t.FailNow()
	}
}
