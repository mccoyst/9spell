// © 2012 Steve McCoy. Available under the MIT license.

package main

import (
	"testing"
)

func TestIsWordSep(t *testing.T) {
	tests := []struct{
		r rune
		ok bool
	}{
		{ '.', false },
		{ 'a', false },
		{ '(', true },
		{ ']', true },
		{ '$', true },
		{ '©', true },
		{ '*', true },
		{ 'ß', false },
	}

	for _, test := range tests {
		is := isWordSep(test.r)
		if is != test.ok {
			t.Error("isWordSep(", test.r, ") should be ", test.ok, "but is", is)
		}
	}
}
