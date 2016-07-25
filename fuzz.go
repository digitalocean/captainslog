// +build gofuzz

package captainslog

// Fuzz is for use with gofozz.
func Fuzz(data []byte) int {
	// data = data.append('\n')
	p := NewParser()
	_, err := p.ParseBytes(data)
	if err != nil {
		return 0
	}
	return 1
}
