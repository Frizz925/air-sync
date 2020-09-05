package util

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientIP(t *testing.T) {
	expectedIP := "10.0.1.0"
	ips := expectedIP + ", 0.0.0.0"
	{
		header := make(http.Header)
		header.Set("X-Real-IP", ips)
		clientIP := GetClientIP(&http.Request{Header: header})
		require.Equal(t, expectedIP, clientIP)
	}
	{
		header := make(http.Header)
		header.Set("X-Forwarded-For", ips)
		clientIP := GetClientIP(&http.Request{Header: header})
		require.Equal(t, expectedIP, clientIP)
	}
	{
		clientIP := GetClientIP(&http.Request{RemoteAddr: expectedIP + ":0"})
		require.Equal(t, expectedIP, clientIP)
	}
}
