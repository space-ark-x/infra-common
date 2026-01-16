package dto

type Response[T any] struct {
	Err  bool   `json:"err"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type ListResponse[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
}
