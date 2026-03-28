package response

type DeleteHistory struct {
	OK bool `json:"ok"`
}

func NewDeleteHistory(ok bool) *DeleteHistory {
	return &DeleteHistory{OK: ok}
}
