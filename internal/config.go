package internal

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	RigctldAddr     string
	MqttAddr        string
	MqttUser        string
	MqttPass        string
	MqttClientId    string
	Topic           string
	PollInterval    time.Duration
	HassDiscovery   bool
	HassName        string
	HassDisplayName string
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	value = strings.ToLower(value)
	return value == "true" || value == "1" || value == "yes" || value == "on"
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

	mqttTopic := getEnv("MQTT_TOPIC", "rigctld")
	hassDiscovery := getEnvBool("HASS_DISCOVERY", false)
	
	pollIntervalStr := getEnv("POLL_INTERVAL", "10s")
	interval, err := time.ParseDuration(pollIntervalStr)
	if err != nil {
		return nil, fmt.Errorf("invalid POLL_INTERVAL: %w", err)
	}

	return &Config{
		RigctldAddr:     rigctldAddr,
		MqttAddr:        mqttAddr,
		MqttUser:        os.Getenv("MQTT_USER"),
		MqttPass:        os.Getenv("MQTT_PASS"),
		MqttClientId:    getEnv("MQTT_CLIENT_ID", "rigctld"),
		PollInterval:    interval,
		Topic:           mqttTopic,
		HassDiscovery:   hassDiscovery,
		HassName:        strings.ReplaceAll(mqttTopic, "/", "_"),
		HassDisplayName: getEnv("HASS_NAME", "Rigctld"),
	}, nil
}
