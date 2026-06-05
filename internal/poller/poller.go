package poller

import (
	"context"
	"log"
	"net"
	"time"
	"bufio"
	"fmt"

	"g4iyt.uk/rigctld-mqtt/internal/config"
)

type RigctldSensor struct {
	Command string
	Name string
	HAName string
	HAUnit string
	HAIcon string
}

type Poller struct {
	config config.Config
	resultsChan chan<-RigctldCommandResult
	sensors []RigctldSensor
}

type RigctldCommandResult struct {
	Name string
	Command string
	Response string
}

func NewPoller(cfg *config.Config, resultsChan chan<- RigctldCommandResult, sensors []RigctldSensor) *Poller {
	return &Poller {
		config: *cfg,
		resultsChan: resultsChan,
		sensors: sensors,
	}
}

func (p *Poller) Start(ctx context.Context) {
	ticker := time.NewTicker(p.config.PollInterval)
	defer ticker.Stop()

	log.Printf("Polling every %s", p.config.PollInterval)
	p.poll()

	for {
		select {
		case <- ctx.Done():
			log.Println("Poller shutting down")
			return
		case <- ticker.C:
			p.poll()
		}
	}
}

func (p *Poller) command(name string, cmd string, conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	_, err := fmt.Fprintf(conn, "%s\n", cmd)
	if err != nil {
		log.Printf("Failed to send command '%s': %v", cmd, err)
	}

	if scanner.Scan() {
		p.resultsChan <- RigctldCommandResult {
			Name: name,
			Command: cmd,
			Response: scanner.Text(),
		}
	} else {
		if err := scanner.Err(); err != nil {
			log.Printf("Read error for '%s': %v", cmd, err)
		}
		return
	}
}

func (p *Poller) poll() {
	dialer := net.Dialer {
		Timeout: 3 * time.Second,
	}
	conn, err := dialer.Dial("tcp", p.config.RigctldAddr)
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	for _, sensor := range p.sensors {
		p.command(sensor.Name, sensor.Command, conn)
	}
}
