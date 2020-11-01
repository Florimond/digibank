package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Config_HappyPath(t *testing.T) {
	cfg, err := Load("sample_config.json")

	assert.NoError(t, err, "Error while reading config.")
	assert.NotEmpty(t, cfg, "Config is empty.")
}

func Test_Config_FileDoesNotExist(t *testing.T) {
	cfg, err := Load("")

	assert.Error(t, err, "Error while reading config.")
	assert.Empty(t, cfg, "Config is empty.")
}

func ExampleLoad() {
	cfg, err := Load("sample_config.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg.Rest.Endpoint)

	// Output: localhost:8181
}
