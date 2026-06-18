package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgentCnf_String(t *testing.T) {
	t.Parallel()
	cnf := &AgentCnf{Address: "localhost:8080"}
	assert.NotEmpty(t, cnf.String())
}

func TestAgentCnf_Set(t *testing.T) {
	t.Parallel()
	cnf := &AgentCnf{}
	err := cnf.Set("localhost:9090")
	require.NoError(t, err)
	assert.Equal(t, "localhost:9090", cnf.Address)
}

func TestAgentCnf_Type(t *testing.T) {
	t.Parallel()
	assert.NotEmpty(t, (&AgentCnf{}).Type())
}

func TestServerCnf_String(t *testing.T) {
	t.Parallel()
	cnf := &ServerCnf{Address: "localhost:8080"}
	assert.Contains(t, cnf.String(), "localhost:8080")
}

func TestServerCnf_Set(t *testing.T) {
	t.Parallel()
	cnf := &ServerCnf{}
	err := cnf.Set("localhost:9090")
	require.NoError(t, err)
	assert.Equal(t, "localhost:9090", cnf.Address)
}

func TestServerCnf_Type(t *testing.T) {
	t.Parallel()
	assert.NotEmpty(t, (&ServerCnf{}).Type())
}

func TestNewServerCnf(t *testing.T) {
	t.Parallel()
	cnf := NewServerCnf()
	assert.Equal(t, "localhost:8080", cnf.Address)
	assert.Equal(t, 300, cnf.StoreInterval)
	assert.Equal(t, "store.json", cnf.FileStoragePath)
	assert.True(t, cnf.Restore)
}

func TestNewAgentCnf(t *testing.T) {
	t.Parallel()
	cnf := NewAgentCnf()
	assert.Equal(t, "localhost:8080", cnf.Address)
	assert.Equal(t, 10, cnf.ReportInterval)
	assert.Equal(t, 2, cnf.PollInterval)
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "cnf_*.json")
	require.NoError(t, err)
	t.Cleanup(func() { os.Remove(f.Name()) })
	_, err = f.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	return f.Name()
}
