package cmd

import (
	"reflect"
	"testing"
)

func TestIntrospectUnencrypted(t *testing.T) {
	if _, err := Introspect([]byte("data:\n  key: dmFsdWU=\nkind: Secret\n")); err == nil {
		t.Fail()
	} else {
		actual := err.Error()
		expected := "Not encrypted"
		if actual != expected {
			t.Fatalf("actual: %#v != expected: %#v", actual, expected)
		}
	}
}

func TestIntrospect(t *testing.T) {
	expected := []string{
		"5E0832A41DCD0FB11C07EF2F9E58C64809FA6916",
		"72ECF46A56B4AD39C907BBB71646B01B86E50310",
	}
	encrypted, err := EncryptWithContext([]byte("data:\n  key: dmFsdWU=\nkind: Secret\n"), EncryptionContext{
		Keys: Keys{
			KeyWithDEK{Key{KTPGP, expected[0]}, nil},
			KeyWithDEK{Key{KTPGP, expected[1]}, nil},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := reconstructEncryptionContext(encrypted, false, false)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := keyIds(ctx, KTPGP)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("actual: %#v != expected: %#v", actual, expected)
	}
}
