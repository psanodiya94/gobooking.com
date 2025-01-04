package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var handler testHandler
	h := NoSurf(&handler)
	switch v := h.(type) {
	case http.Handler:
		// do nothing; test passed
	default:
		t.Errorf("NoSurf did not return an http.Handler, is %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var handler testHandler
	h := SessionLoad(&handler)
	switch v := h.(type) {
	case http.Handler:
		// do nothing; test passed
	default:
		t.Errorf("NoSurf did not return an http.Handler, is %T", v)
	}
}
