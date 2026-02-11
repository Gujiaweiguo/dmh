package svc

import "testing"

func TestServiceContextZeroValue(t *testing.T) {
	var s ServiceContext
	if s.DB != nil {
		t.Fatalf("expected nil db on zero value")
	}
}
