package internal

import (
	"log"
	"fmt"
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func generateDeviceIdentifier(cfg *Config) map[string]interface{} {
	return map[string]interface{}{
		"identifiers": []string{
			fmt.Sprintf("g4iyt_rigctld_bridge_%s", cfg.HassName),
		},
		"name": cfg.HassDisplayName,
		"model": "Rigctld",
		"manufacturer": "G4IYT",
	}
}

func PublishHassDiscovery(mqttClient mqtt.Client, sensors []RigctldSensor, cfg *Config) {
	for _, sensor := range sensors {
		topic := fmt.Sprintf("homeassistant/sensor/%s_%s/config", cfg.HassName, sensor.Name)
		sensorConfig := map[string]interface{}{
			"name": sensor.HAName,
			"unique_id": fmt.Sprintf("%s_%s", cfg.HassName, sensor.Name),
			"state_topic": fmt.Sprintf("%s/state/%s", cfg.Topic, sensor.Name),
			"icon": sensor.HAIcon,
			"device": generateDeviceIdentifier(cfg),
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
