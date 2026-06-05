package config

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
	MqttClient string
	Topic string
	PollInterval time.Duration
	HassDiscovery bool
	HassName string
}

func Load() (*Config, error) {
	rigctldAddr := os.Getenv("RIGCTL_ADDR")
	if rigctldAddr == "" {
		return nil, fmt.Errorf("RIGCTL_ADDR not specified")
	}

	mqttAddr := os.Getenv("MQTT_ADDR")
	if mqttAddr == "" {
		return nil, fmt.Errorf("MQTT_ADDR not specified")
	}

	mqttUser := os.Getenv("MQTT_USER")
	mqttPass := os.Getenv("MQTT_PASS")
	mqttClient := os.Getenv("MQTT_CLIENT")
	if mqttClient == "" {
		mqttClient = "rigctld"
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

	return &Config {
		RigctldAddr: rigctldAddr,
		MqttAddr: mqttAddr,
		MqttUser: mqttUser,
		MqttPass: mqttPass,
		PollInterval: interval,
		Topic: mqttTopic,
		HassDiscovery: hassDiscovery,
		HassName: hassName,
	}, nil
}
