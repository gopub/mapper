# goparam

Assign parameters to model and validate values

```
type Image struct {
	Width  int    `param:"w,min=100,max=800"`
	Height int    `param:"h,min=100,max=800"`
	Link   string `param:"pattern=url"`
}

type Topic struct {
	Title      string   `param:"min=2,max=30"`
	CoverImage *Image   `param:"optional"`
	MoreImages []*Image `param:"optional"`
}
```
    
## Validate
    topic := &Topic{
    	Title: "a",
    }
    err := goparam.Validate(topic)
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
	err := goparam.Assign(&topic, params)
	//...