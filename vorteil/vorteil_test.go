package vorteil

import "testing"

func TestConfig(t *testing.T) {
	config := new(configuration)
	err := config.load("./../config")
	if err != nil {
		t.Error(nil)
	}
}
