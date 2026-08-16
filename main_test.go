package main

import "testing"

func TestValidateName(t *testing.T) {
	cases := map[string]bool{
		"my-server":          true,
		"server_2":           true,
		"Server123":          true,
		"":                   false,
		"server with spaces": false,
		"server/../etc":      false,
		"café":               false,
	}

	for input, wantOK := range cases {
		err := validateName(input)
		gotOK := err == nil
		if gotOK != wantOK {
			t.Errorf("validateName(%q) = err %v, expected ok=%v", input, err, wantOK)
		}
	}
}

func TestValidatePositiveInt(t *testing.T) {
	cases := map[string]bool{
		"20":  true,
		"1":   true,
		"0":   false,
		"-5":  false,
		"abc": false,
		"":    false,
	}

	for input, wantOK := range cases {
		err := validatePositiveInt(input)
		gotOK := err == nil
		if gotOK != wantOK {
			t.Errorf("validatePositiveInt(%q) = err %v, expected ok=%v", input, err, wantOK)
		}
	}
}
