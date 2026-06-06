package internal

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

type RigctldSensor struct {
	Command string
	Name    string
	HAName  string
	HAUnit  string
	HAIcon  string
}

type Poller struct {
	config      Config
	resultsChan chan<- RigctldCommandResult
	sensors     []RigctldSensor
}

type RigctldCommandResult struct {
	Name     string
	Command  string
	Response string
}

func NewPoller(cfg *Config, resultsChan chan<- RigctldCommandResult, sensors []RigctldSensor) *Poller {
	return &Poller{
		config:      *cfg,
		resultsChan: resultsChan,
		sensors:     sensors,
	}
}

func (p *Poller) Start(ctx context.Context) {
	log.Printf("Starting poller, interval: %s", p.config.PollInterval)

	ticker := time.NewTicker(p.config.PollInterval)
	defer ticker.Stop()

	var conn net.Conn
	var err error

	for {
		if conn == nil {
			conn, err = p.connect(ctx)
			if err != nil {
				log.Printf("Connection failed: %v, retrying in 5s", err)
				select {
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Second):
					continue
				}
			}
		}

		if err := p.poll(conn); err != nil {
			log.Printf("Poll failed: %v, closing connection", err)
			conn.Close()
			conn = nil
		}

		select {
		case <-ctx.Done():
			if conn != nil {
				conn.Close()
			}
			log.Println("Poller shutting down")
			return
		case <-ticker.C:
			// Next poll
		}
	}
}

func (p *Poller) connect(ctx context.Context) (net.Conn, error) {
	dialer := net.Dialer{
		Timeout: 3 * time.Second,
	}
	conn, err := dialer.DialContext(ctx, "tcp", p.config.RigctldAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (p *Poller) poll(conn net.Conn) error {
	for _, sensor := range p.sensors {
		if err := p.command(sensor.Name, sensor.Command, conn); err != nil {
			return err
		}
	}
	return nil
}

func (p *Poller) command(name string, cmd string, conn net.Conn) error {
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	_, err := fmt.Fprintf(conn, "%s\n", cmd)
	if err != nil {
		return fmt.Errorf("failed to send command '%s': %w", cmd, err)
	}

	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		p.resultsChan <- RigctldCommandResult{
			Name:     name,
			Command:  cmd,
			Response: scanner.Text(),
		}
		return nil
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read error for '%s': %w", cmd, err)
	}
	return fmt.Errorf("no response for '%s'", cmd)
}
