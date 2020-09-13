package messages

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Tag    string `json:"tag,omitempty"`
	Castle string `json:"castle"`
	Level  int    `json:"level"`
	Health int    `json:"hp"`
}
