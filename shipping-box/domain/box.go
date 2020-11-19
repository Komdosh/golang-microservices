package domain

type Box struct {
	Length int
	Width  int
	Height int
}

func (f Box) GetLength() int {
	return f.Length
}

func (f Box) GetWidth() int {
	return f.Width
}

func (f Box) GetHeight() int {
	return f.Height
}
