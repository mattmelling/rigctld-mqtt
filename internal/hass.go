package internal

import (
	"log"
	"fmt"
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var hsDevice = map[string]interface{}{
	"identifiers": []string{"g4iyt_rigctld_bridge"},
	"name": "Rigctld",
	"model": "Rigctld",
	"manufacturer": "G4IYT",
}

func PublishHassDiscovery(mqttClient mqtt.Client, sensors []RigctldSensor, cfg *Config) {
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
