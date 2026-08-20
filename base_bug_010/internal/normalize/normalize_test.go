package normalize

import "testing"

func TestMessage(t *testing.T) {
	if Message("token=abc 123") != "token=<redacted> <n>" {
		t.Fatal("normalize failed")
	}
}
