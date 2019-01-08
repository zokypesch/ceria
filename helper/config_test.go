package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigByName(t *testing.T) {
	// Test null return
	help := NewReadConfigService()
	err := help.Init()

	assert.Equal(t, "", help.GetByName(""))
	assert.NoError(t, err)
	assert.NotEmpty(t, help.GetByName("db.host"))
	// test all data its also get empty string because db is root
	assert.Empty(t, help.GetByName("db"))
}
