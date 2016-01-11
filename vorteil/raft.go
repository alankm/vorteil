package vorteil

import "gopkg.in/inconshreveable/log15.v2"

func (v *Vorteil) raft() error {
	v.start = v.raftStart
	v.stop = v.raftStop

	v.log = log15.New("vorteil", "raft")

	return nil
}

func (v *Vorteil) raftStart() error {
	return nil
}

func (v *Vorteil) raftStop() {

}
