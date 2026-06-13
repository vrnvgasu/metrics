package config

import (
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
