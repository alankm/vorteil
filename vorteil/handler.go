package vorteil

import (
	"net/http"
	"strings"

	"github.com/alankm/privileges"
)

type serviceHandler struct {
	vorteil *Vorteil
	handle  func(*privileges.Session, http.ResponseWriter, *http.Request)
}

func (sh serviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// default login
	key, err := sh.vorteil.users.Login("guest", "")
	if err != nil {
		w.Write(ErrInternal.Response(""))
		return
	}

	// default access to webpages
	if !strings.HasPrefix(r.URL.Path, "/services/admin") && !strings.HasPrefix(r.URL.Path, "/services/images") {
		sh.handle(key, w, r)
		return
	}

	// check cookie
	if cookie, err := r.Cookie("vorteil"); err == nil {
		value := make(map[string]string)
		err = sh.vorteil.cookie.Decode("vorteil", cookie.Value, &value)
		if err == nil {
			username := value["username"]
			password := value["password"]
			key, err := sh.vorteil.users.Login(username, password)
			if err == nil {
				sh.handle(key, w, r)
				return
			}
		}
	}

	w.Write(ErrAccessDenied.Response(""))
	return

}
