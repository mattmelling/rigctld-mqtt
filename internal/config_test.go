package internal

import (
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	t.Run("Default values", func(t *testing.T) {
		// Set required env vars and override existing ones for this test
		t.Setenv("RIGCTLD_ADDR", "localhost:4532")
		t.Setenv("MQTT_ADDR", "tcp://localhost:1883")
		t.Setenv("MQTT_TOPIC", "")
		t.Setenv("POLL_INTERVAL", "")
		t.Setenv("HASS_DISCOVERY", "")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		if cfg.RigctldAddr != "localhost:4532" {
			t.Errorf("Expected RigctldAddr localhost:4532, got %s", cfg.RigctldAddr)
		}
		if cfg.MqttClientId != "rigctld" {
			t.Errorf("Expected default MqttClientId rigctld, got %s", cfg.MqttClientId)
		}
		if cfg.PollInterval != 10*time.Second {
			t.Errorf("Expected default PollInterval 10s, got %v", cfg.PollInterval)
		}
		if cfg.HassDiscovery != false {
			t.Errorf("Expected default HassDiscovery false, got %v", cfg.HassDiscovery)
		}
	})

	t.Run("Custom values", func(t *testing.T) {
		t.Setenv("RIGCTLD_ADDR", "localhost:4532")
		t.Setenv("MQTT_ADDR", "tcp://localhost:1883")
		t.Setenv("MQTT_TOPIC", "ham/radio")
		t.Setenv("POLL_INTERVAL", "5s")
		t.Setenv("HASS_DISCOVERY", "true")

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}

		if cfg.Topic != "ham/radio" {
			t.Errorf("Expected Topic ham/radio, got %s", cfg.Topic)
		}
		if cfg.PollInterval != 5*time.Second {
			t.Errorf("Expected PollInterval 5s, got %v", cfg.PollInterval)
		}
		if cfg.HassDiscovery != true {
			t.Errorf("Expected HassDiscovery true, got %v", cfg.HassDiscovery)
		}
		if cfg.HassName != "ham_radio" {
			t.Errorf("Expected HassName ham_radio, got %s", cfg.HassName)
		}
	})
}
