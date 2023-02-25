package shared

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
