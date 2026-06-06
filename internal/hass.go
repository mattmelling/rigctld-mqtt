package internal

import (
	"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type HassDevice struct {
	Identifiers  []string `json:"identifiers"`
	Name         string   `json:"name"`
	Model        string   `json:"model"`
	Manufacturer string   `json:"manufacturer"`
}

type HassSensorConfig struct {
	Name              string      `json:"name"`
	UniqueID          string      `json:"unique_id"`
	ObjectID          string      `json:"object_id"`
	StateTopic        string      `json:"state_topic"`
	Icon              string      `json:"icon,omitempty"`
	UnitOfMeasurement string      `json:"unit_of_measurement,omitempty"`
	Device            *HassDevice `json:"device,omitempty"`
}

func PublishHassDiscovery(mqttClient mqtt.Client, sensors []RigctldSensor, cfg *Config) {
	device := &HassDevice{
		Identifiers: []string{
			fmt.Sprintf("g4iyt_rigctld_bridge_%s", cfg.HassName),
		},
		Name:         cfg.HassDisplayName,
		Model:        "Rigctld",
		Manufacturer: "G4IYT",
	}

	for _, sensor := range sensors {
		topic := fmt.Sprintf("homeassistant/sensor/%s_%s/config", cfg.HassName, sensor.Name)
		sensorConfig := HassSensorConfig{
			Name:       sensor.HAName,
			UniqueID:   fmt.Sprintf("%s_%s", cfg.HassName, sensor.Name),
			ObjectID:   fmt.Sprintf("%s_%s", cfg.HassName, sensor.Name),
			StateTopic: fmt.Sprintf("%s/state/%s", cfg.Topic, sensor.Name),
			Icon:       sensor.HAIcon,
			Device:     device,
		}

		if sensor.HAUnit != "" {
			sensorConfig.UnitOfMeasurement = sensor.HAUnit
		}

		payload, err := json.Marshal(sensorConfig)
		if err != nil {
			log.Printf("Failed to marshal sensor config for %s: %v", sensor.Name, err)
			continue
		}

		token := mqttClient.Publish(topic, 1, true, payload)
		if token.Wait() && token.Error() != nil {
			log.Printf("Failed to publish discovery for %s: %v", sensor.Name, token.Error())
		}
	}
}
