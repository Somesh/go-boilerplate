package lib

import (
	"context"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Advanced Unicode normalization and filtering,
// see http://blog.golang.org/normalization and
// http://godoc.org/golang.org/x/text/unicode/norm for more
// details.
func stripCtlAndExtFromUnicode(ctx context.Context, str string) string {
	isOk := func(r rune) bool {
		return r < 32 || (r >= 33 && r <= 44) || r == 46 || r == 47 || (r >= 58 && r <= 64) || (r >= 91 && r <= 96) || r >= 123
	}
	// The isOk filter is such that there is no need to chain to norm.NFC
	t := transform.Chain(norm.NFKD, transform.RemoveFunc(isOk))
	// This Transformer could also trivially be applied as an io.Reader
	// or io.Writer filter to automatically do such filtering when reading
	// or writing data anywhere.
	str, _, _ = transform.String(t, str)
	return str
}
