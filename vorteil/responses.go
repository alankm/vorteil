package vorteil

import "encoding/json"

const (
	StatusMsgSuccess  = "success"
	StatusCodeSuccess = 200

	StatusMsgFail  = "error"
	StatusCodeFail = 500
)

var (
	Success = &GeneralResponse{
		Status: StatusMsgSuccess,
		Code:   StatusCodeSuccess,
	}
)

var (
	// Vorteil
	ErrNotLeader    = NewErrResponse(1000, "server is not current leader")
	ErrInternal     = NewErrResponse(1001, "an internal error occurred")
	ErrBadLoginBody = NewErrResponse(1001, "that body of the login request was invalid")

	// Privileges
	ErrInvalidLogin          = NewErrResponse(2000, "username or password is invalid")
	ErrAccessDenied          = NewErrResponse(2001, "user does not have permissions to perform this action")
	ErrBadCommand            = NewErrResponse(2002, "command parameter not set to valid command")
	ErrMissingGroupHeader    = NewErrResponse(2010, "command requires a group header argument")
	ErrMissingUsernameHeader = NewErrResponse(2011, "command requires a username header argument")
	ErrMissingPasswordHeader = NewErrResponse(2012, "command requires a password header argument")

	// Vimages
	ErrNoFileAtTarget = NewErrResponse(3000, "no image or folder exists at the specified target")
)

type Response interface {
	Response(string) []byte
}

type GeneralResponse struct {
	Status  string   `json:"status"`
	Code    int      `json:"status_code"`
	Leader  string   `json:"leader_info,omitempty"`
	Payload Response `json:"payload,omitempty"`
}

func NewGeneralResponse(code int, msg string, payload Response) *GeneralResponse {
	r := new(GeneralResponse)
	r.Code = code
	r.Status = msg
	r.Payload = payload
	return r
}

func (r *GeneralResponse) Response(leader string) []byte {
	r.Leader = leader
	a, _ := json.Marshal(r)
	return a
}

type ErrResponse struct {
	Code    int    `json:"error_code"`
	Message string `json:"error_message"`
	super   GeneralResponse
}

func NewErrResponse(code int, msg string) *ErrResponse {
	r := new(ErrResponse)
	r.super.Payload = r
	r.super.Status = StatusMsgFail
	r.super.Code = StatusCodeFail
	r.Code = code
	r.Message = msg
	return r
}

func (r *ErrResponse) Response(leader string) []byte {
	r.super.Leader = leader
	a, _ := json.Marshal(r.super)
	return a
}

type ListUsersResponse struct {
	Users []string `json:"users"`
	super GeneralResponse
}

func NewListUsersResponse(users []string) *ListUsersResponse {
	r := new(ListUsersResponse)
	r.super.Code = StatusCodeSuccess
	r.super.Status = StatusMsgSuccess
	r.super.Payload = r
	r.Users = users
	return r
}

func (r *ListUsersResponse) Response(leader string) []byte {
	r.super.Leader = leader
	a, _ := json.Marshal(r.super)
	return a
}

type ListGroupsResponse struct {
	Groups []string `json:"groups"`
	super  GeneralResponse
}

func NewListGroupsResponse(groups []string) *ListGroupsResponse {
	r := new(ListGroupsResponse)
	r.super.Code = StatusCodeSuccess
	r.super.Status = StatusMsgSuccess
	r.super.Payload = r
	r.Groups = groups
	return r
}

func (r *ListGroupsResponse) Response(leader string) []byte {
	r.super.Leader = leader
	a, _ := json.Marshal(r.super)
	return a
}

type GroupResponse struct {
	Group string `json:"group"`
	super GeneralResponse
}

func NewGroupResponse(group string) *GroupResponse {
	r := new(GroupResponse)
	r.super.Code = StatusCodeSuccess
	r.super.Status = StatusMsgSuccess
	r.super.Payload = r
	r.Group = group
	return r
}

func (r *GroupResponse) Response(leader string) []byte {
	r.super.Leader = leader
	a, _ := json.Marshal(r.super)
	return a
}

type UserInGroupResponse struct {
	Bool  bool `json:"user_is_in_group"`
	super GeneralResponse
}

func NewUserInGroupResponse(b bool) *UserInGroupResponse {
	r := new(UserInGroupResponse)
	r.super.Code = StatusCodeSuccess
	r.super.Status = StatusMsgSuccess
	r.super.Payload = r
	r.Bool = b
	return r
}

func (r *UserInGroupResponse) Response(leader string) []byte {
	r.super.Leader = leader
	a, _ := json.Marshal(r.super)
	return a
}
