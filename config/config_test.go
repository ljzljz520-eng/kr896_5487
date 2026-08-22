package config

import "testing"

func TestConfigDefaults(t *testing.T) {
	c := Defaults()
	if c.Address == "" || c.DBPath == "" || c.StaticDir == "" {
		t.Fatal("defaults incomplete")
	}
	if c.Validate() != nil {
		t.Fatal("defaults invalid")
	}
}
