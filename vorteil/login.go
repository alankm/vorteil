package vorteil

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alankm/privileges"
)

func (v *Vorteil) loginHandler(s *privileges.Session, w http.ResponseWriter, r *http.Request) {

	dec := json.NewDecoder(r.Body)
	val := make(map[string]string)
	err := dec.Decode(&val)
	if err != nil {
		w.Write(ErrBadLoginBody.Response(""))
		return
	}

	username := val["username"]
	password := val["password"]

	_, err = v.users.Login(username, password)
	if err != nil {
		w.Write(ErrInvalidLogin.Response(""))
		return
	}

	if enc, err := v.cookie.Encode("vorteil", val); err == nil {
		cookie := &http.Cookie{
			Name:  "vorteil",
			Value: enc,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
		w.Write(Success.Response(""))
		return
	}

	fmt.Println("A")

	w.Write(ErrInternal.Response(""))
	return

}
