package client

import (
	"reflect"
	"testing"
)

func TestClient_NewClient(t *testing.T) {
	c := NewClient(nil, nil, nil)

	if c == nil {
		t.Error("Client init error")
	}

	if reflect.TypeOf(c).String() != "*client.Client" {
		t.Error("Wrong type error ", reflect.TypeOf(c).String(), "*client.Client")
	}
}
