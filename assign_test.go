package param_test

import (
	"github.com/natande/goparam"
	"testing"
)

func TestAssign(t *testing.T) {
	params := map[string]interface{}{
		"title": "this is title",
		"cover_image": map[string]interface{}{
			"w":    100,
			"h":    200,
			"link": "https://www.image.com",
		},
		"more_images": []map[string]interface{}{
			{
				"w":    100,
				"h":    200,
				"link": "https://www.image.com",
			},
		},
	}

	var topic *Topic
	err := param.Assign(&topic, params)
	if err != nil {
		t.FailNow()
	}
}

func TestAssignSlice(t *testing.T) {
	params := map[string]interface{}{
		"title": "this is title",
		"cover_image": map[string]interface{}{
			"w":    100,
			"h":    200,
			"link": "https://www.image.com",
		},
		"more_images": []map[string]interface{}{
			{
				"w":    100,
				"h":    200,
				"link": "https://www.image.com",
			},
		},
	}

	values := []interface{}{params}
	var topics []*Topic
	err := param.Assign(&topics, values)
	if err != nil || len(topics) == 0 {
		t.FailNow()
	}
}
