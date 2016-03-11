package core

import (
	"encoding/json"
	"net/http"

	"github.com/alankm/vorteil/core/privileges"
	"github.com/gorilla/mux"
)

func (c *Core) Login(username, password string) (*privileges.Session, error) {
	return nil, nil
}

func (c *Core) routeLogin(subrouter *mux.Router) {

}

func (c *Core) loginHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	val := make(map[string]string)
	err := dec.Decode(&val)
	if err != nil {
		w.Write(ResponseBadLoginBody.JSON())
		return
	}
	username := val["username"]
	password := val["password"]
	if username != "root" || password != "guest" {
		w.Write(ResponseBadLogin.JSON())
		return
	}
	session, err := c.Login(username, password)
	if err != nil {
		w.Write(ResponseBadLogin.JSON())
		return
	}
	enc, err := c.web.cookie.Encode("vorteil", val)
	if err != nil {
		panic(CodeInternal)
	}
	cookie := &http.Cookie{
		Name:  "vorteil",
		Value: enc,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	var services []string
	services = make([]string, 0)
	for _, val := range c.web.services {
		if session.Access(val) == nil {
			services = append(services, val)
		}
	}
	w.Write(NewSuccessResponse(services).JSON())
}
