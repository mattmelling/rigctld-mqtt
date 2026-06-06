package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	. "g4iyt.uk/rigctld-mqtt/internal"
)

var sensors = []RigctldSensor {
	{
		Command: "f",
		Name: "frequency",
		HAName: "Frequency",
		HAUnit: "Hz",
		HAIcon: "mdi:radio-tower",
	},
	{
		Command: "m",
		Name: "mode",
		HAName: "Mode",
		HAUnit: "",
		HAIcon: "mdi:wave-form",
	},
	{
		Command: "t",
		Name: "ptt",
		HAName: "PTT",
		HAIcon: "mdi:radio-tower",
	},
}

func main() {

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	opts := mqtt.NewClientOptions().AddBroker(cfg.MqttAddr)
	if cfg.MqttUser != "" {
		opts.SetUsername(cfg.MqttUser)
	}
	if cfg.MqttPass != "" {
		opts.SetPassword(cfg.MqttPass)
	}
	if cfg.MqttClientId != "" {
		opts.SetClientID(cfg.MqttClientId)
	}
	opts.SetAutoReconnect(true)
	opts.SetConnectTimeout(5 * time.Second)

	mqttClient := mqtt.NewClient(opts)
	log.Printf("Connecting to %s", cfg.MqttAddr)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("MQTT connection error: %v", token.Error())
	}
	defer mqttClient.Disconnect(250)

	if cfg.HassDiscovery {
		PublishHassDiscovery(mqttClient, sensors, cfg)
	}
	
	resultsChan := make(chan RigctldCommandResult, 100)
	go func() {
		for {
			select {
			case res := <-resultsChan:
				// log.Printf("%s = %s", res.Command, res.Response)
				topic := fmt.Sprintf("%s/state/%s", cfg.Topic, res.Name)
				mqttClient.Publish(topic, 1, false, res.Response)
			case <-ctx.Done():
				return
			}
		}
	}()
	
	daemonPoller := NewPoller(cfg, resultsChan, sensors)
	daemonPoller.Start(ctx)
}
	
