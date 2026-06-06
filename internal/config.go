package internal

import (
	"fmt"
	"os"
	"time"
	"strings"
)

type Config struct {
	RigctldAddr string
	MqttAddr string
	MqttUser string
	MqttPass string
	MqttClientId string
	Topic string
	PollInterval time.Duration
	HassDiscovery bool
	HassName string
	HassDisplayName string
}

func LoadConfig() (*Config, error) {
	rigctldAddr := os.Getenv("RIGCTLD_ADDR")
	if rigctldAddr == "" {
		return nil, fmt.Errorf("RIGCTLD_ADDR not specified")
	}

	mqttAddr := os.Getenv("MQTT_ADDR")
	if mqttAddr == "" {
		return nil, fmt.Errorf("MQTT_ADDR not specified")
	}

	mqttUser := os.Getenv("MQTT_USER")
	mqttPass := os.Getenv("MQTT_PASS")
	mqttClientId := os.Getenv("MQTT_CLIENT_ID")
	if mqttClientId == "" {
		mqttClientId = "rigctld"
	}

	mqttTopic := os.Getenv("MQTT_TOPIC")
	if mqttTopic == "" {
		mqttTopic = "rigctld"
	}
	
	pollIntervalStr := os.Getenv("POLL_INTERVAL")
	if pollIntervalStr == "" {
		pollIntervalStr = "10s"
	}

	interval, err := time.ParseDuration(pollIntervalStr)
	if err != nil {
		return nil, fmt.Errorf("Invalid POLL_INTERVAL: %w", err)
	}

	hassDiscoveryStr := os.Getenv("HASS_DISCOVERY")
	hassDiscovery := false
	if hassDiscoveryStr != "" {
		hassDiscovery = true
	}

	hassName := ""
	if hassDiscovery {
		hassName = strings.Replace(mqttTopic, "/", "_", -1)
	}

	hassDisplayName := os.Getenv("HASS_NAME")
	if hassDisplayName == "" {
		hassDisplayName = "Rigctld"
	}

	return &Config {
		RigctldAddr: rigctldAddr,
		MqttAddr: mqttAddr,
		MqttUser: mqttUser,
		MqttPass: mqttPass,
		MqttClientId: mqttClientId,
		PollInterval: interval,
		Topic: mqttTopic,
		HassDiscovery: hassDiscovery,
		HassName: hassName,
		HassDisplayName: hassDisplayName,
	}, nil
}
