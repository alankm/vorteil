package vorteil

import (
	"encoding/hex"
	"net/http"

	"github.com/alankm/privileges"
)

type AdminHeaders struct {
	grp, usr, pwd []string
	gok, uok, pok bool
}

var adminHandlers = map[string]func(*privileges.Session, *AdminHeaders) Response{
	"NewGroup":            adminNewGroup,
	"NewUser":             adminNewUser,
	"SetPassword":         adminSetPassword,
	"GetGid":              adminGetGid,
	"SetGid":              adminSetGid,
	"UserIsInGroup":       adminUserIsInGroup,
	"AddUserToGroup":      adminAddUserToGroup,
	"RemoveUserFromGroup": adminRemoveUserFromGroup,
	"ListUsers":           adminListUsers,
	"ListGroups":          adminListGroups,
	"ListUsersGroups":     adminListUsersGroups,
	"ListUsersWithGid":    adminListUsersWithGid,
	"DeleteUser":          adminDeleteUser,
	"DeleteGroup":         adminDeleteGroup,
}

func (v *Vorteil) adminHandler(s *privileges.Session, w http.ResponseWriter, r *http.Request) {

	leader := ""

	err := r.ParseForm()
	if err != nil {
		w.Write(ErrInternal.Response(leader))
		return
	}

	var h AdminHeaders
	h.grp, h.gok = r.Header["Group"]
	h.usr, h.uok = r.Header["Username"]
	h.pwd, h.pok = r.Header["Password"]

	handler := adminHandlers[r.FormValue("command")]
	if handler == nil {
		w.Write(ErrBadCommand.Response(leader))
		return
	}
	w.Write(handler(s, &h).Response(leader))

}

func adminNewGroup(s *privileges.Session, h *AdminHeaders) Response {

	if !h.gok {
		return ErrMissingGroupHeader
	}
	err := s.NewGroup(h.grp[0])
	if err == nil {
		return Success
	}

	if err.Error() == "access denied" {
		return ErrAccessDenied
	}

	return ErrInternal

}

func adminNewUser(s *privileges.Session, h *AdminHeaders) Response {

	if !h.uok {
		return ErrMissingUsernameHeader
	}
	if !h.pok {
		return ErrMissingPasswordHeader
	}

	salt := hex.EncodeToString(privileges.GenerateSalt64())
	hash, _ := privileges.Hash(salt, h.pwd[0])
	err := s.NewUser(h.usr[0], salt, hash)
	if err == nil {
		return Success
	}

	if err.Error() == "access denied" {
		return ErrAccessDenied
	}

	return ErrInternal

}

func adminSetPassword(s *privileges.Session, h *AdminHeaders) Response {

	if !h.uok {
		return ErrMissingUsernameHeader
	}
	if !h.pok {
		return ErrMissingPasswordHeader
	}

	salt := hex.EncodeToString(privileges.GenerateSalt64())
	hash, _ := privileges.Hash(salt, h.pwd[0])
	err := s.ChangePassword(h.usr[0], salt, hash)
	if err == nil {
		return Success
	}

	if err.Error() == "access denied" {
		return ErrAccessDenied
	}

	return ErrInternal

}

func adminRemoveUserFromGroup(s *privileges.Session, h *AdminHeaders) Response {

	if !h.uok {
		return ErrMissingUsernameHeader
	}
	if !h.gok {
		return ErrMissingGroupHeader
	}
	err := s.UserRemoveGroup(h.usr[0], h.grp[0])
	if err == nil {
		return Success
	}

	if err.Error() == "access denied" {
		return ErrAccessDenied
	}

	return ErrInternal

}

func adminAddUserToGroup(s *privileges.Session, h *AdminHeaders) Response {

	if !h.uok {
		return ErrMissingUsernameHeader
	}
	if !h.gok {
		return ErrMissingGroupHeader
	}
	err := s.UserAddGroup(h.usr[0], h.grp[0])
	if err == nil {
		return Success
	}

	if err.Error() == "access denied" {
		return ErrAccessDenied
	}

	return ErrInternal

}

func adminGetGid(s *privileges.Session, h *AdminHeaders) Response {

	if !h.uok {
		return ErrMissingUsernameHeader
	}
	group, _ := s.Gid(h.usr[0], "")
	return NewGroupResponse(group)

}

func adminSetGid(s *privileges.Session, h *AdminHeaders) Response {

	if !h.uok {
		return ErrMissingUsernameHeader
	}
	if !h.gok {
		return ErrMissingGroupHeader
	}
	_, err := s.Gid(h.usr[0], h.grp[0])
	if err == nil {
		return Success
	}

	if err.Error() == "access denied" {
		return ErrAccessDenied
	}

	return ErrInternal

}

func adminUserIsInGroup(s *privileges.Session, h *AdminHeaders) Response {

	if !h.uok {
		return ErrMissingUsernameHeader
	}
	if !h.gok {
		return ErrMissingGroupHeader
	}
	b, _ := s.UserInGroup(h.usr[0], h.grp[0])
	return NewUserInGroupResponse(b)

}

func adminListGroups(s *privileges.Session, h *AdminHeaders) Response {

	groups, _ := s.ListGroups()
	return NewListGroupsResponse(groups)

}

func adminListUsers(s *privileges.Session, h *AdminHeaders) Response {

	users, _ := s.ListUsers()
	return NewListUsersResponse(users)

}

func adminListUsersGroups(s *privileges.Session, h *AdminHeaders) Response {

	if !h.uok {
		return ErrMissingUsernameHeader
	}
	groups, _ := s.UserListGroups(h.usr[0])
	return NewListGroupsResponse(groups)

}

func adminListUsersWithGid(s *privileges.Session, h *AdminHeaders) Response {

	if !h.gok {
		return ErrMissingGroupHeader
	}
	users, _ := s.GroupListUsersGids(h.grp[0])
	return NewListUsersResponse(users)

}

func adminDeleteUser(s *privileges.Session, h *AdminHeaders) Response {

	if !h.uok {
		return ErrMissingUsernameHeader
	}
	err := s.DeleteUser(h.usr[0])
	if err == nil {
		return Success
	}

	if err.Error() == "access denied" {
		return ErrAccessDenied
	}

	return ErrInternal

}

func adminDeleteGroup(s *privileges.Session, h *AdminHeaders) Response {

	if !h.gok {
		return ErrMissingGroupHeader
	}
	err := s.DeleteGroup(h.grp[0])
	if err == nil {
		return Success
	}

	if err.Error() == "access denied" {
		return ErrAccessDenied
	}

	return ErrInternal

}
