# params

Assign parameters to model and validate values

```
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
```
    
## Validate
    topic := &Topic{
    	Title: "a",
    }
    err := params.Validate(topic)
    //...
    
## Assign
Assign will validate result automatically

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
	err := params.Assign(&topic, params)
	//...