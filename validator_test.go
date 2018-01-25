package mapper_test

import (
	"github.com/gopub/mapper"
	"strings"
	"testing"
)

type Image struct {
	Width  int    `mapper:"w,min=100,max=800"`
	Height int    `mapper:"h,min=100,max=800"`
	Link   string `mapper:"pattern=url"`
}

type Topic struct {
	Title      string   `mapper:"min=2,max=30"`
	CoverImage *Image   `mapper:"optional"`
	MoreImages []*Image `mapper:"optional"`
}

func TestValidate(t *testing.T) {
	topic := &Topic{
		Title: "a",
	}
	err := mapper.Validate(topic)
	if err == nil {
		t.FailNow()
	}

	topic.Title = strings.Repeat("a", 31)
	err = mapper.Validate(topic)
	if err == nil {
		t.FailNow()
	}

	topic.Title = strings.Repeat("a", 30)
	err = mapper.Validate(topic)
	if err != nil {
		t.FailNow()
	}

	topic.Title = strings.Repeat("a", 2)
	err = mapper.Validate(topic)
	if err != nil {
		t.FailNow()
	}

	topic.CoverImage = &Image{
		Width:  100,
		Height: 0,
		Link:   "https://www.image.com",
	}
	err = mapper.Validate(topic)
	if err == nil {
		t.FailNow()
	}

	topic.CoverImage = &Image{
		Width:  100,
		Height: 900,
		Link:   "https://www.image.com",
	}
	err = mapper.Validate(topic)
	if err == nil {
		t.FailNow()
	}

	topic.CoverImage = &Image{
		Width:  100,
		Height: 800,
		Link:   "https://www.image.com",
	}
	err = mapper.Validate(topic)
	if err != nil {
		t.FailNow()
	}

	topic.MoreImages = []*Image{{
		Width:  0,
		Height: 800,
		Link:   "https://www.image.com",
	}}
	err = mapper.Validate(topic)
	if err == nil {
		t.FailNow()
	}

	topic.MoreImages = []*Image{{
		Width:  100,
		Height: 800,
		Link:   ":",
	}}
	err = mapper.Validate(topic)
	if err == nil {
		t.FailNow()
	}

	topic.MoreImages = []*Image{{
		Width:  100,
		Height: 800,
		Link:   "https://www.image.com",
	}}
	err = mapper.Validate(topic)
	if err != nil {
		t.FailNow()
	}
}
