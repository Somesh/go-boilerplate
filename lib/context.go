
package lib

import (
	"context"
	"net/http"
	"net/url"
	"strings"

)

const (
	ipAddressKey           int = iota //0
	userAgentKey                      //1
	queryParams                       //2
	args                              //3
	accessTokenKey                    //4
)

func ParseRequestContext(req *http.Request) context.Context {
	ctx := req.Context()

	ctx = parseIPAddress(ctx, req)
	ctx = parseUserAgent(ctx, req)
	ctx = parseQueryParams(ctx, req)
	ctx = parseArguments(ctx, req)
	ctx = parseAccessToken(ctx, req)

	return ctx
}


func GetIPAddress(ctx context.Context) string {
	val, _ := ctx.Value(ipAddressKey).(string)
	return val
}

func GetUserAgent(ctx context.Context) string {
	val, _ := ctx.Value(userAgentKey).(string)
	return val
}

func GetQueryParams(ctx context.Context) string {
	val, _ := ctx.Value(queryParams).(string)
	return val
}

func GetArguments(ctx context.Context) string {
	val, _ := ctx.Value(args).(string)
	return val
}

func GetAccessToken(ctx context.Context) string {
	val, _ := ctx.Value(accessTokenKey).(string)
	return val
}

func parseIPAddress(ctx context.Context, req *http.Request) context.Context {
	forwardedIP := req.Header.Get("X-Forwarded-For")
	if forwardedIP == "" {
		forwardedIP = "127.0.0.1"
	}

	return context.WithValue(ctx, ipAddressKey, forwardedIP)
}

func parseUserAgent(ctx context.Context, req *http.Request) context.Context {
	userAgent := req.UserAgent()
	if userAgent == "" {
		userAgent = "http/curl" // SS: TODO Move this to constant
	}
	return context.WithValue(ctx, userAgentKey, userAgent)
}

func parseAccessToken(ctx context.Context, req *http.Request) context.Context {
	var token, bearerToken string
	bearerToken = req.Header.Get("Accounts-Authorization")
	if bearerToken == "" {
		bearerToken = req.Header.Get("Authorization")
	}
	bearer := strings.Split(bearerToken, "Bearer ")
	if len(bearer) > 1 {
		token = bearer[1]
	}
	return context.WithValue(ctx, accessTokenKey, token)
}

func parseQueryParams(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, queryParams, req.URL.RawQuery)
}


func parseArguments(ctx context.Context, req *http.Request) context.Context {
	queryparams, _ := url.ParseQuery(req.URL.RawQuery)
	return context.WithValue(ctx, args, queryparams.Encode())
}
