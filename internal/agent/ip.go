package agent

import (
	"fmt"
	"net"
)

// OutboundIP определяет исходящий IP-адрес хоста агента,
// который используется для соединения с адресом address.
// Соединение по UDP не отправляет пакетов, а лишь выбирает локальный интерфейс.
func OutboundIP(address string) (net.IP, error) {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return nil, fmt.Errorf("agent.OutboundIP Dial: %w", err)
	}
	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP, nil
}
