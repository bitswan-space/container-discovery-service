package config

import (
	"os"

	"bitswan.space/container-discovery-service/internal/logger"
	"gopkg.in/yaml.v2"
)

type Configuration struct {
	MQTTBrokerUrl        string   `yaml:"mqtt-broker-url"`
	MQTTContainersPub    string   `yaml:"mqtt-containers-pub"`
	MQTTTopologyTopics   []string `yaml:"mqtt-topology-topics"`
	MQTTNavigationPub    string   `yaml:"mqtt-navigation-topic"`
	MQTTNavigationSet    string   `yaml:"mqtt-navigation-set"`
	NavigationFile       string   `yaml:"navigation-file"`
	NavigationSchemaFile string   `yaml:"navigation-schema-file"`
	DashboardEntriesFile string   `yaml:"dashboard-entries-file"`
}

var cfg *Configuration

func LoadConfig(filename string) error {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(buf, &cfg); err != nil {
		return err
	}
	logger.Info.Printf("Successfuly loaded configuration")
	return nil
}

func GetConfig() *Configuration {
	return cfg
}
