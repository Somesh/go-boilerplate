package config

import (
	"testing"
	"time"
)

func TestToGraceConfig(t *testing.T) {
	input := GraceCfg{Timeout: "1m1s"}
	result := input.ToGraceConfig()
	if result.HTTPReadTimeout.String() != "10s" {
		t.Errorf("ToGraceConfig: Expected: 10s, Got: %s", result.HTTPReadTimeout.String())
	}
}

func TestGetConfig(t *testing.T) {
	CF.Timeout = time.Duration(1)
	if GetConfig().Timeout != time.Duration(1) {
		t.Errorf("GetConfig: Expected: %+v, Got: %+v", time.Duration(1),GetConfig().Timeout)
	}
}
