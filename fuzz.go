// +build gofuzz

package captainslog

func Fuzz(data []byte) int {
	msg := NewSyslogMsg()
	if err := Unmarshal(data, &msg); err != nil {
		return 0
	}
	return 1
}
