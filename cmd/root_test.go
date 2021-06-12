package cmd

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	viper.AddConfigPath("..")
	initConfig()

	testConfigProperties := []struct {
		prop string
		exp  interface{}
	}{
		{"consumer.kafka.address", []interface{}{"kafka:9092"}},
		{"consumer.kafka.topic", "events"},
		{"consumer.kafka.partition", int64(0)},
	}

	for _, tc := range testConfigProperties {
		t.Run(tc.prop, func(tt *testing.T) {
			ass := assert.New(tt)
			cfg := viper.Get(tc.prop)
			ass.Equal(tc.exp, cfg)
		})
	}
}

func TestRootCommand(t *testing.T) {
	ass := assert.New(t)
	output, err := executeCommand(rootCmd, "")
	ass.NoError(err)
	ass.Contains(output, rootCmd.Long, "Unexpected output.")
	ass.Contains(output, rootCmd.UsageString(), "Unexpected output.")
}
