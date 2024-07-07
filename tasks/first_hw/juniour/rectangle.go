package juniour


type Rectangle struct {
	Width int
	Height int
}

func (r Rectangle) Area() int {
	return r.Width * r.Height
}