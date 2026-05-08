package audit

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileAudit struct {
	file    *os.File
	encoder *json.Encoder
}

func NewFileAudit(filename string) (*FileAudit, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open audit file: %w", err)
	}

	encoder := json.NewEncoder(file)

	return &FileAudit{
		file:    file,
		encoder: encoder,
	}, nil
}

func (a *FileAudit) Observe(e EventMessage) error {
	return a.encoder.Encode(e)
}

func (a *FileAudit) Close() error {
	return a.file.Close()
}
