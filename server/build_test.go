package server

import (
	"strings"
	"testing"
)

func TestBuildHash(t *testing.T) {
	hash := buildHash("os", "arch", "arm", "a,b,c,d")

	if !strings.Contains(hash, "os") {
		t.Error("Expected hash to contain 'os', but it didn't")
	}
	if !strings.Contains(hash, "arch") {
		t.Error("Expected hash to contain 'arch', but it didn't")
	}
	if !strings.Contains(hash, "arm") {
		t.Error("Expected hash to contain 'arm', but it didn't")
	}
	if !strings.Contains(hash, "a,b,c,d") {
		t.Error("Expected hash to contain 'a,b,c,d', but it didn't")
	}
}
