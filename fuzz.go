// +build gofuzz

package captainslog

func Fuzz(data []byte) int {
	// data = data.append('\n')
	msg := NewSyslogMsg()
	if err := Unmarshal(data, &msg); err != nil {
		return 0
	}
	return 1
}
