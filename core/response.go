package core

import "encoding/json"

const (
	CodeInternal = 3000 + iota
	CodeDatabase
)

var (
	Success                    = &ResponseWrapper{200, "success", nil}
	ResponsePrivilegesInternal = NewFailResponse(0, "internal error")
	ResponsePrivilegesDatabase = NewFailResponse(0, "unexpected database error")
	ResponseBadLogin           = NewFailResponse(0, "invalid username or password")
	ResponseAccessDenied       = NewFailResponse(0, "access denied")
	ResponseVimagesInternal    = NewFailResponse(0, "internal error")
	ResponseVimagesDatabase    = NewFailResponse(0, "unexpected database error")
	ResponseVorteilInternal    = NewFailResponse(0, "internal error")
	ResponseLeader             = NewFailResponse(0, "server isn't current raft leader")
	ResponseAuthentication     = NewFailResponse(0, "not logged into server")
	ResponseBadLoginBody       = NewFailResponse(0, "bad login body")
	ResponseBadAdminCommand    = NewFailResponse(0, "bad admin command")
	ResponseMissingUsername    = NewFailResponse(0, "missing username header")
	ResponseMissingGroup       = NewFailResponse(0, "missing group header")
	ResponseMissingPassword    = NewFailResponse(0, "missing password header")
	ResponseCHMOD              = NewFailResponse(0, "chmod placeholder")
	ResponseCHOWN              = NewFailResponse(0, "chown placeholder")
	ResponseCHGRP              = NewFailResponse(0, "chgrp placeholder")
)

type ResponseWrapper struct {
	Code    int         `json:"status_code"`
	Msg     string      `json:"status"`
	Payload interface{} `json:"payload,omitempty"`
}

func (r *ResponseWrapper) OK() bool {
	return r.Code == 200
}

func (r *ResponseWrapper) JSON() []byte {
	a, _ := json.Marshal(r)
	return a
}

func NewSuccessResponse(payload interface{}) *ResponseWrapper {
	wrapper := &ResponseWrapper{
		Code:    200,
		Msg:     "success",
		Payload: payload,
	}
	return wrapper
}

type ErrorResponse struct {
	wrapper *ResponseWrapper
	Code    int               `json:"code"`
	Msg     string            `json:"status"`
	Info    map[string]string `json:"info,omitempty"`
}

func NewFailResponse(code int, msg string, info ...string) *ErrorResponse {
	wrapper := &ResponseWrapper{
		Code: 500,
		Msg:  "fail",
	}
	resp := &ErrorResponse{
		wrapper: wrapper,
		Code:    code,
		Msg:     msg,
	}
	wrapper.Payload = resp
	resp.SetInfo(info...)
	return resp
}

func (r *ErrorResponse) AddInfo(info ...string) {
	if len(info)%2 != 0 {
		panic(CodeInternal)
	}
	for i := 0; i < len(info); i = i + 2 {
		r.Info[info[i]] = info[i+1]
	}
}

func (r *ErrorResponse) SetInfo(info ...string) *ErrorResponse {
	if len(info) == 0 {
		r.Info = nil
		return r
	}
	r.Info = make(map[string]string)
	r.AddInfo(info...)
	return r
}

func (r *ErrorResponse) OK() bool {
	return r.Code == 200
}

func (r *ErrorResponse) JSON() []byte {
	a, _ := json.Marshal(r.wrapper)
	return a
}
