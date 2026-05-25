package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	models "github.com/vrnvgasu/metrics/internal/model"
)

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open producer file: %w", err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return &Producer{
		file:    file,
		encoder: encoder,
	}, nil
}
func (p *Producer) WriteMetrics(l *models.MetricsList) error {
	if _, err := p.file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	if err := p.file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}

	return p.encoder.Encode(l)
}
func (p *Producer) Close() error {
	return p.file.Close()
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open consumer file: %w", err)
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}
func (c *Consumer) ReadMetrics() (models.MetricsList, error) {
	l := make(models.MetricsList, 0)
	if err := c.decoder.Decode(&l); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	return l, nil
}
func (c *Consumer) Close() error {
	return c.file.Close()
}
