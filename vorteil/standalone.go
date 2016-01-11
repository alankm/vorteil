package vorteil

import "gopkg.in/inconshreveable/log15.v2"

func (v *Vorteil) standalone() error {
	v.start = v.standaloneStart
	v.stop = v.standaloneStop

	v.log = log15.New("vorteil", "standalone")

	return nil
}

func (v *Vorteil) standaloneStart() error {
	return v.http.Start()
}

func (v *Vorteil) standaloneStop() {
	v.http.Stop()
}
