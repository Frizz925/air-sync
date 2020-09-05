package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseTimeDuration(t *testing.T) {
	{
		_, err := ParseTimeDuration("1")
		require.NotNil(t, err)
	}
	{
		d, err := ParseTimeDuration("1s")
		require.Nil(t, err)
		require.Equal(t, 1*time.Second, d)
	}
	{
		d, err := ParseTimeDuration("10m")
		require.Nil(t, err)
		require.Equal(t, 10*time.Minute, d)
	}
	{
		d, err := ParseTimeDuration("24h")
		require.Nil(t, err)
		require.Equal(t, 24*time.Hour, d)
	}
}
