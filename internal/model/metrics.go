package models

import (
	"fmt"
	"strconv"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.
type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

func NewMetricsFromStrings(mType, ID, value string) (Metrics, error) {

	m := &Metrics{
		ID:   ID,
		Hash: "",
	}
	if err := m.setTypeFromString(mType); err != nil {
		return *m, fmt.Errorf("model: unable to set metrics type: %w", err)
	}
	if err := m.setValueFromString(value); err != nil {
		return *m, fmt.Errorf("model: unable to set metrics value: %w", err)
	}

	return *m, nil
}

func (m *Metrics) setTypeFromString(v string) error {
	switch v {
	case Counter:
		m.MType = Counter
		return nil
	case Gauge:
		m.MType = Gauge
		return nil
	default:
		return fmt.Errorf("invalid metric type: %s", v)
	}
}

func (m *Metrics) setValueFromString(v string) error {
	switch m.MType {
	case Counter:
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("invalid metric value: %w", err)
		}
		m.Value = &value
		return nil
	case Gauge:
		delta, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid metric delta: %w", err)
		}
		m.Delta = &delta
		return nil
	default:
		return nil
	}
}
