package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"fmt"
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"g4iyt.uk/rigctld-mqtt/internal/config"
	"g4iyt.uk/rigctld-mqtt/internal/poller"
)

var sensors = []poller.RigctldSensor {
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

var hsDevice = map[string]interface{}{
	"identifiers": []string{"g4iyt_rigctld_bridge"},
	"name": "Rigctld",
	"model": "Rigctld",
	"manufacturer": "G4IYT",
}

func publishHassDiscovery(mqttClient mqtt.Client, cfg *config.Config) {
	for _, sensor := range sensors {
		topic := fmt.Sprintf("homeassistant/sensor/%s_%s/config", cfg.HassName, sensor.Name)
		sensorConfig := map[string]interface{}{
			"name": sensor.HAName,
			"unique_id": fmt.Sprintf("%s_%s", cfg.HassName, sensor.Name),
			"state_topic": fmt.Sprintf("%s/state/%s", cfg.Topic, sensor.Name),
			"icon": sensor.HAIcon,
			"device": hsDevice,
		}

		if sensor.HAUnit != "" {
			sensorConfig["unit_of_measurement"] = sensor.HAUnit
		}

		payload, err := json.Marshal(sensorConfig)
		if err != nil {
			log.Printf("Failed to marshal sensor config: %v", err)
			return
		}
		
		token := mqttClient.Publish(topic, 1, true, payload)
		token.Wait()
	}
}

func main() {

	cfg, err := config.Load()
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
	opts.SetAutoReconnect(true)
	opts.SetConnectTimeout(5 * time.Second)

	mqttClient := mqtt.NewClient(opts)
	log.Printf("Connecting to %s", cfg.MqttAddr)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("MQTT connection error: %v", token.Error())
	}
	defer mqttClient.Disconnect(250)

	if cfg.HassDiscovery {
		publishHassDiscovery(mqttClient, cfg)
	}
	
	resultsChan := make(chan poller.RigctldCommandResult, 100)
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
	
	daemonPoller := poller.NewPoller(cfg, resultsChan, sensors)
	daemonPoller.Start(ctx)
}
	
