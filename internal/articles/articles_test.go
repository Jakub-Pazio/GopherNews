package articles

import (
	"testing"
)

func TestStrictContains(t *testing.T) {
	var positiveTests = []struct {
		articleName string
		key  string
	}{
			{"Learn C++", "C++"},
			{"C++ and more", "C++"},
	}
	var negativeTests = []struct {
		articleName string
		key  string
	}{
			{"Learn C", "C++"},
			{"Is Glo bad?", "Go"},
	}
	for _, test := range positiveTests {
		testName := test.articleName + " contains " + test.key
		t.Run(testName, func(t *testing.T) {
			if !StrictContains(test.articleName, test.key) {
				t.Errorf("strictContains(\"%s\", \"%s\") should be true", test.articleName, test.key)
			}
		})
	}
	for _, test := range negativeTests {
		testName := test.articleName + " does not contain " + test.key
		t.Run(testName, func(t *testing.T) {
			if StrictContains(test.articleName, test.key) {
				t.Errorf("strictContains(\"%s\", \"%s\") should be false", test.articleName, test.key)
			}
		})
	}
}
