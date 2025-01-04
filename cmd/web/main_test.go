package main

import (
	"testing"
)

func TestRun(t *testing.T) {
	err := run()
	if err != nil {
		t.Errorf("run() error = %v", err)
		return
	}
}
