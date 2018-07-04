package controllers

// Status 返回状态
const (
	StatusSuccess  = "success"
	StatusBadReq   = "bad_req"
	StatusNull     = "null"
	StatusNotLogin = "not_login"
	StatusError    = "error"
	StatusNotValid = "not_valid"
	StatusNotAllow = "not_allow"
	StatusNotFound = "not_found"
)

// CommonRes 一般返回值
type CommonRes struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}
