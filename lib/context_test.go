package lib

import (
	"context"
	"net/http"
	"testing"

)

func TestGetChannelID(t *testing.T) {
	expectedValue := "1"

	ctx := context.WithValue(context.Background(), channelIDKey, expectedValue)

	if devID := GetChannelID(ctx); devID != expectedValue {
		t.Errorf("unexpected result = \"%s\", want \"%s\"", devID, expectedValue)
	}
}

func TestGetIPAddress(t *testing.T) {
	expectedValue := "192.168.1.1"

	ctx := context.WithValue(context.Background(), ipAddressKey, expectedValue)

	if ip := GetIPAddress(ctx); ip != expectedValue {
		t.Errorf("unexpected result = \"%s\", want \"%s\"", ip, expectedValue)
	}
}

func TestGetUserAgent(t *testing.T) {
	expectedValue := "user-agent-foo"

	ctx := context.WithValue(context.Background(), userAgentKey, expectedValue)

	if ua := GetUserAgent(ctx); ua != expectedValue {
		t.Errorf("unexpected result = \"%s\", want \"%s\"", ua, expectedValue)
	}
}



func TestParseIPAddress(t *testing.T) {
	ctx := context.Background()
	req := &http.Request{RemoteAddr: "127.0.0.1:8080"}

	ctx = parseIPAddress(ctx, req)
	if ip := GetIPAddress(ctx); ip != "127.0.0.1" {
		t.Errorf("unexpected result = \"%s\", want \"%s\"", ip, "127.0.0.1")
	}
}

func TestParseUserAgent(t *testing.T) {
	ctx := context.Background()
	req := &http.Request{Header: http.Header{}}

	expectedValue := "foo"

	req.Header.Set("User-Agent", expectedValue)
	ctx = parseUserAgent(ctx, req)

	if ua := GetUserAgent(ctx); ua != expectedValue {
		t.Errorf("unexpected result = \"%s\", want \"%s\"", ua, expectedValue)
	}
}

func TestParseAccessToken(t *testing.T) {
	ctx := context.Background()
	req := &http.Request{Header: http.Header{}}

	expectedValue := "foo"

	req.Header.Set("Authorization", "Bearer "+expectedValue)
	ctx = parseAccessToken(ctx, req)

	if token := GetAccessToken(ctx); token != expectedValue {
		t.Errorf("unexpected result = \"%s\", want \"%s\"", token, expectedValue)
	}
}

func TestParseRequestContext(t *testing.T) {
	req, _ := http.NewRequest("", "", nil)
	ParseRequestContext(req)
}

