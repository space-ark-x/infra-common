package dto

type Response struct {
	Err  bool   `json:"err"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type ListResponse struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
}
