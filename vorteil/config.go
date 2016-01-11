package vorteil

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type configuration struct {
	Mode      string               `yaml:"mode"`
	Data      string               `yaml:"data"`
	Bind      string               `yaml:"bind"`
	Advertise string               `yaml:"advertise"`
	Storage   storageConfiguration `yaml:"storage"`
	Raft      raftConfig           `yaml:"raft"`
}

type raftConfig struct {
}

func (c *configuration) load(target string) error {

	src, err := ioutil.ReadFile(target)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(src, c)
	if err != nil {
		return err
	}

	return nil

}
