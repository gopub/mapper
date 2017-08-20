# goparam

Assign parameters to model and validate values

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
    
    topic := &Topic{
    	Title: "a",
    }
    err := Validate(topic)
    //...