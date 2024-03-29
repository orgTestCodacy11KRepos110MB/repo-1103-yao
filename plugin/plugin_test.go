package plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yaoapp/gou"
	"github.com/yaoapp/yao/config"
)

func TestLoad(t *testing.T) {
	gou.Plugins = make(map[string]*gou.Plugin)
	Load(config.Conf)
	LoadFrom("404")
	check(t)
}

func check(t *testing.T) {
	keys := []string{}
	for key := range gou.Plugins {
		keys = append(keys, key)
	}
	assert.Equal(t, 1, len(keys))
}
