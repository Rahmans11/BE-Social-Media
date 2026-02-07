package dto

type Response struct {
	Msg     string         `json:"msg"`
	Success bool           `json:"success"`
	Data    []any          `json:"data"`
	Error   string         `json:"error,omitempty"`
	Meta    PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Page      int    `json:"page,omitempty"`
	TotalPage int    `json:"total_page,omitempty"`
	NextPage  string `json:"next_page,omitempty"`
	PrevPage  string `json:"prev_page,omitempty"`
}
