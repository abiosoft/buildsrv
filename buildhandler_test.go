package main

import "testing"

func TestBuildHandler(t *testing.T) {
	// TODO: Hmm, how to test this cleanly.
	// Maybe simple error cases, like missing query string parameters
}

func TestCheckInput(t *testing.T) {
	err := checkInput("linux", "amd64", "", nil)
	if err != nil {
		t.Errorf("Expected no errors when input is good, but got '%v'", err)
	}

	err = checkInput("", "amd64", "", nil)
	if err == nil {
		t.Error("Expected error when os missing")
	}

	err = checkInput("linux", "", "", nil)
	if err == nil {
		t.Error("Expected error when arch missing")
	}

	err = checkInput("bad_os", "amd64", "", nil)
	if err == nil {
		t.Error("Expected error when os is invalid")
	}

	err = checkInput("linux", "bad_arch", "", nil)
	if err == nil {
		t.Error("Expected error when arch is invalid")
	}

	err = checkInput("linux", "amd64", "bad_arm", nil)
	if err == nil {
		t.Error("Expected error when arm is invalid")
	}

	err = checkInput("linux", "amd64", "", []string{"alsdjfkaskldfjsjhfskdjhfskdjfhhkjhsk"})
	if err == nil {
		t.Error("Expected error when a feature is invalid")
	}
}
