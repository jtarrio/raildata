package raildata_test

import (
	"testing"

	"github.com/jtarrio/raildata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestColor(t *testing.T) {
	{
		color, err := raildata.ParseHtmlColor("#123def")
		require.NoError(t, err)
		assert.Equal(t, "#123def", color.Html())
		r, g, b := color.RGB()
		assert.Equal(t, 0x12, r)
		assert.Equal(t, 0x3d, g)
		assert.Equal(t, 0xef, b)
	}
	{
		color, err := raildata.ParseHtmlColor("#F71")
		require.NoError(t, err)
		assert.Equal(t, "#ff7711", color.Html())
		r, g, b := color.RGB()
		assert.Equal(t, 0xff, r)
		assert.Equal(t, 0x77, g)
		assert.Equal(t, 0x11, b)
	}
	{
		_, err := raildata.ParseHtmlColor("#F711")
		assert.Error(t, err)
	}
	{
		_, err := raildata.ParseHtmlColor("123def")
		assert.Error(t, err)
	}
	{
		_, err := raildata.ParseHtmlColor("123d")
		assert.Error(t, err)
	}
}
