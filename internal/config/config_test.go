package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	value := Default()
	if value.DBPath == "" || value.Listen == "" {
		t.Fatal("empty config")
	}
}
