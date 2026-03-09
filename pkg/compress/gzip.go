package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func GzipCompress(data []byte) ([]byte, error) {
	buf := bytes.Buffer{}
	writer := gzip.NewWriter(&buf)

	if _, err := writer.Write(data); err != nil {
		return nil, fmt.Errorf("failed to compress data: %w", err)
	}

	writer.Close()

	return buf.Bytes(), nil
}
