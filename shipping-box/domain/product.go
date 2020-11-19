package domain

type Product struct {
	Name   string
	Length int
	Width  int
	Height int
}

func (f Product) GetLength() int {
	return f.Length
}

func (f Product) GetWidth() int {
	return f.Width
}

func (f Product) GetHeight() int {
	return f.Height
}
