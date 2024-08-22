package lib

import (
	"context"
	"testing"
	// "time"
)

func TestStripCtlAndExtFromUnicode(t *testing.T) {
	src := "déjà vu" + // precomposed unicode
		"\n\000\037 \041\176\177\200\377\n" + // various boundary cases
		"as⃝df̅" // unicode combining characters

	str := stripCtlAndExtFromUnicode(context.Background(), src)
	if str != "deja vu asdf" {
		t.Error("Unable to normalize string ", str)
	}

	src = "<script type='text/javascript'>alert(1)</script>" // unicode combining characters

	str = stripCtlAndExtFromUnicode(context.Background(), src)
	if str != "script typetextjavascriptalert1script" {
		t.Error("Unable to normalize string ", str)
	}

	src = "asd}}gigxa'\\/\"<gertu" // unicode combining characters

	src = "abc - tst"
	str = stripCtlAndExtFromUnicode(context.Background(), src)
	t.Log(str)
}
