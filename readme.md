# Assign map to struct
```
type Image struct {
	Width  int    
	Height int    
	Link   string 
}

type Topic struct {
	Title      string  
	CoverImage *Image  
	MoreImages []*Image 
}

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
	
```

# Assign struct to struct
``` 
type UserRow struct {
    ID int64,
    Name string,
    CreatedAt int64
}

type User struct {
    ID int64,
    Name string
}

row, err := getUserRow(userID)
var user *User
if err == nil {
    mapper.Assign(user, row)
}
```
    
## Validate model
``` 
type Image struct {
	Width  int    `mapper:"min=100,max=800"`
	Height int    `mapper:"min=100,max=800"`
	Link   string `mapper:"pattern=url"`
	Format string `mapper:"optional"`
}

var img := &Image{
    ...
}

err := mapper.Validate(img)

```
    
