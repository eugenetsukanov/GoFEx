package main


type MyReader struct{}

func (m MyReader) Read(b []byte) (int, error) {
	for x := range b {
		b[x] = 'A'
	}

	return len(b), nil
}

func main() {
	reader.Validate(MyReader{})
}
