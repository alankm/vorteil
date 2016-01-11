package vorteil

import (
	"os"
	"strings"

	"github.com/alankm/privileges"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/sisatech/multiserver"

	"gopkg.in/inconshreveable/log15.v2"
)

type Vorteil struct {
	config configuration
	cookie *securecookie.SecureCookie
	users  *privileges.Privileges
	http   *multiserver.HTTPServer
	start  func() error
	stop   func()
	log    log15.Logger
}

func New(target string) (*Vorteil, error) {

	v := new(Vorteil)
	err := v.config.load(target)
	if err != nil {
		return v, err
	}

	err = v.setup()
	return v, err

}

func (v *Vorteil) setup() error {

	v.config.Data = strings.TrimRight(v.config.Data, "/")
	err := os.MkdirAll(v.config.Data, 0777)
	if err != nil {
		return err
	}

	v.cookie = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
	v.users, err = privileges.New(v.config.Data + "/vorteil.db")
	if err != nil {
		return err
	}

	v.initHTTP()

	switch v.config.Mode {
	case "standalone":
		return v.standalone()
	case "raft":
		return v.raft()
	default:
		return ErrMode
	}

}

func (v *Vorteil) Start() error {
	return v.start()
}

func (v *Vorteil) Stop() {
	v.stop()
}

func (v *Vorteil) initHTTP() {

	router := mux.NewRouter()

	loginSH := &serviceHandler{
		vorteil: v,
		handle:  v.loginHandler,
	}
	adminSH := &serviceHandler{
		vorteil: v,
		handle:  v.adminHandler,
	}

	router.Handle("/services/login", loginSH).Methods("POST")
	router.Handle("/services/admin", adminSH)

	v.http = multiserver.NewHTTPServer(v.config.Bind, router, nil /* TODO: implement TLS */)

}
