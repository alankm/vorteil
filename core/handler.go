package core

import (
	"fmt"
	"net/http"

	"github.com/alankm/vorteil/core/privileges"
)

type GeneralHandler struct {
	handler func(http.ResponseWriter, *http.Request)
}

func (h *GeneralHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("A")
	h.handler(w, r)
}

type ProtectedHandler struct {
	core    *Core
	handler func(*privileges.Session, http.ResponseWriter, *http.Request)
}

func (h *ProtectedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if x := recover(); x != nil {
			h.core.logs.Crit("Panic response", "trace", x)
		}
	}()

	if h.core.raft.server.Leader() != h.core.raft.config.Advertise && h.core.raft.server.Leader() != h.core.raft.config.Bind {
		h.core.logs.Debug("rerouting to leader")
		w.Write(ResponseLeader.SetInfo("leader", h.core.raft.server.Leader()).JSON())
		return
	}

	// login pass-through
	if r.URL.Path == "/services/login" {
		h.handler(nil, w, r)
		return
	}

	h.core.logs.Debug("Restricted request recieved", "method", r.Method, "target", r.URL.Path)

	// check if logged in
	s := h.core.HandlerLogin(r)
	if s == nil {
		h.core.logs.Debug("non-logged in access attempt")
		w.Write(ResponseAuthentication.JSON())
		return
	}

	h.core.logs.Debug("User credentials validated", "username", s.User())

	h.handler(s, w, r)

}

func (c *Core) HandlerLogin(r *http.Request) *privileges.Session {
	if cookie, err := r.Cookie("vorteil"); err == nil {
		value := make(map[string]string)
		err = c.web.cookie.Decode("vorteil", cookie.Value, &value)
		if err == nil {
			username := value["username"]
			password := value["password"]
			session, err := c.Login(username, password)
			if err == nil {
				return session
			}
		}
	}
	return nil
}
