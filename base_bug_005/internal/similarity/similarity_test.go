package similarity

import "testing"

func TestScoreUsesUniqueTokens(t *testing.T) {
	if Similar("db timeout", "db db db") {
		t.Fatal("duplicate tokens inflated score")
	}
	if Score("", " ") != 1 {
		t.Fatal("empty messages should match")
	}
}
