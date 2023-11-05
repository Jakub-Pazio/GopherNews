package articles

import (
	"testing"
)

func TestStrictContains(t *testing.T) {
	t.Run("strict contains positive", func(t *testing.T) {
		if !StrictContains("Learn C++", "C++") {
			t.Error("strictContains(\"Learn C++\", \"C++\") should be true")
		}

		if !StrictContains("Learn C++ and more", "C++") {
			t.Error("strictContains(\"Learn C++ and more\", \"C++\") should be true")
		}
	})
}
