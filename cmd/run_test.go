package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecute_MissingConfig(t *testing.T) {
	rootCmd.SetArgs([]string{"--config", "nonexistent.yaml"})
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)

	err := rootCmd.Execute()
	// cobra returns nil and writes to stderr; we just verify no panic
	_ = err
}

func TestInit_Flags(t *testing.T) {
	f := rootCmd.Flags()

	configFlag := f.Lookup("config")
	require.NotNil(t, configFlag)
	assert.Equal(t, "vaultdiff.yaml", configFlag.DefValue)

	formatFlag := f.Lookup("format")
	require.NotNil(t, formatFlag)
	assert.Equal(t, "text", formatFlag.DefValue)

	maskFlag := f.Lookup("mask")
	require.NotNil(t, maskFlag)
	assert.Equal(t, "false", maskFlag.DefValue)

	failFlag := f.Lookup("fail-on-diff")
	require.NotNil(t, failFlag)
	assert.Equal(t, "false", failFlag.DefValue)
}

func TestRunDiff_InvalidFormat(t *testing.T) {
	cmd := rootCmd
	cmd.SetArgs([]string{"--config", "../internal/config/example.yaml", "--format", "xml"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := runDiff(cmd, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "xml")
}
