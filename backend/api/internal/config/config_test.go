package config

import "testing"

func TestConfigFields(t *testing.T) {
	var c Config
	c.Host = "127.0.0.1"
	c.Port = 8889
	c.Auth.AccessSecret = "secret"
	if c.Host != "127.0.0.1" || c.Port != 8889 || c.Auth.AccessSecret != "secret" {
		t.Fatalf("config field assignment failed")
	}
}
