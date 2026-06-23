package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseDurationSec(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input   string
		wantSec int
		wantOk  bool
	}{
		{"1s", 1, true},
		{"5s", 5, true},
		{"1m", 60, true},
		{"2m30s", 150, true},
		{"300s", 300, true},
		{"invalid", 0, false},
		{"300", 0, false},
		{"", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			got, ok := parseDurationSec(tt.input)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantSec, got)
		})
	}
}

func TestServerCnf_ApplyFile(t *testing.T) {
	t.Parallel()
	restoreTrue := true
	restoreFalse := false
	tests := []struct {
		name  string
		input ServerFileCnf
		check func(t *testing.T, cnf *ServerCnf)
	}{
		{
			name: "all fields",
			input: ServerFileCnf{
				Address: "myserver:9090", Restore: &restoreFalse,
				StoreInterval: "5s", FileStoragePath: "/tmp/store.db",
				DatabaseDSN: "host=db", CryptoKey: "/tmp/key.pem",
			},
			check: func(t *testing.T, cnf *ServerCnf) {
				assert.Equal(t, "myserver:9090", cnf.Address)
				assert.False(t, cnf.Restore)
				assert.Equal(t, 5, cnf.StoreInterval)
				assert.Equal(t, "/tmp/store.db", cnf.FileStoragePath)
				assert.Equal(t, "host=db", cnf.DatabaseDSN)
				assert.Equal(t, "/tmp/key.pem", cnf.CryptoKey)
			},
		},
		{
			name:  "empty file keeps defaults",
			input: ServerFileCnf{},
			check: func(t *testing.T, cnf *ServerCnf) {
				assert.Equal(t, "localhost:8080", cnf.Address)
				assert.Equal(t, 300, cnf.StoreInterval)
				assert.True(t, cnf.Restore)
			},
		},
		{
			name:  "restore nil keeps default true",
			input: ServerFileCnf{Restore: nil},
			check: func(t *testing.T, cnf *ServerCnf) { assert.True(t, cnf.Restore) },
		},
		{
			name:  "restore true pointer",
			input: ServerFileCnf{Restore: &restoreTrue},
			check: func(t *testing.T, cnf *ServerCnf) { assert.True(t, cnf.Restore) },
		},
		{
			name:  "invalid duration ignored",
			input: ServerFileCnf{StoreInterval: "invalid"},
			check: func(t *testing.T, cnf *ServerCnf) { assert.Equal(t, 300, cnf.StoreInterval) },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cnf := NewServerCnf()
			cnf.ApplyFile(&tt.input)
			tt.check(t, cnf)
		})
	}
}

func TestAgentCnf_ApplyFile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input AgentFileCnf
		check func(t *testing.T, cnf *AgentCnf)
	}{
		{
			name: "all fields",
			input: AgentFileCnf{
				Address: "agentserver:8081", ReportInterval: "30s",
				PollInterval: "5s", CryptoKey: "/tmp/pub.pem",
			},
			check: func(t *testing.T, cnf *AgentCnf) {
				assert.Equal(t, "agentserver:8081", cnf.Address)
				assert.Equal(t, 30, cnf.ReportInterval)
				assert.Equal(t, 5, cnf.PollInterval)
				assert.Equal(t, "/tmp/pub.pem", cnf.CryptoKey)
			},
		},
		{
			name:  "empty file keeps defaults",
			input: AgentFileCnf{},
			check: func(t *testing.T, cnf *AgentCnf) {
				assert.Equal(t, "localhost:8080", cnf.Address)
				assert.Equal(t, 10, cnf.ReportInterval)
				assert.Equal(t, 2, cnf.PollInterval)
			},
		},
		{
			name:  "invalid duration ignored",
			input: AgentFileCnf{ReportInterval: "bad", PollInterval: "also-bad"},
			check: func(t *testing.T, cnf *AgentCnf) {
				assert.Equal(t, 10, cnf.ReportInterval)
				assert.Equal(t, 2, cnf.PollInterval)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cnf := NewAgentCnf()
			cnf.ApplyFile(&tt.input)
			tt.check(t, cnf)
		})
	}
}

func TestLoadFileCnf(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		content string
		loadFn  func(string) error
		wantErr bool
	}{
		{
			name:    "server valid json",
			content: `{"address":"localhost:9090","store_interval":"10s"}`,
			loadFn:  func(f string) error { _, err := LoadServerFileCnf(f); return err },
		},
		{
			name:    "server invalid json",
			content: `{invalid}`,
			loadFn:  func(f string) error { _, err := LoadServerFileCnf(f); return err },
			wantErr: true,
		},
		{
			name:    "agent valid json",
			content: `{"address":"localhost:9091","report_interval":"10s"}`,
			loadFn:  func(f string) error { _, err := LoadAgentFileCnf(f); return err },
		},
		{
			name:    "agent invalid json",
			content: `{invalid}`,
			loadFn:  func(f string) error { _, err := LoadAgentFileCnf(f); return err },
			wantErr: true,
		},
		{
			name:    "server file not found",
			content: "",
			loadFn:  func(_ string) error { _, err := LoadServerFileCnf("/nonexistent"); return err },
			wantErr: true,
		},
		{
			name:    "agent file not found",
			content: "",
			loadFn:  func(_ string) error { _, err := LoadAgentFileCnf("/nonexistent"); return err },
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var path string
			if tt.content != "" {
				path = writeTempFile(t, tt.content)
			}
			err := tt.loadFn(path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
