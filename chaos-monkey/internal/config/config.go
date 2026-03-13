package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// Victim cluster
	KubeconfigData  string // base64-encoded kubeconfig for Cluster A
	VictimNamespace string
	ClusterAUUID    string // UpCloud UUID of Cluster A


	// Scheduler
	MinInterval time.Duration
	MaxInterval time.Duration

	// Chaos action weights (must sum to 100)
	WeightPodKill       int
	WeightNodeReboot    int
	WeightNodeStop      int
	WeightNetworkDisrupt int
	WeightScaleZero     int

	// Action config
	NodeStopDuration    time.Duration
	ScaleZeroDuration   time.Duration
	NetworkDisruptDelay string

	// Prometheus metrics port
	MetricsPort string
}

func Load() (*Config, error) {
	viper.AutomaticEnv()

	viper.SetDefault("VICTIM_NAMESPACE", "victim")
	viper.SetDefault("MIN_INTERVAL", "30s")
	viper.SetDefault("MAX_INTERVAL", "120s")
	viper.SetDefault("WEIGHT_POD_KILL", 40)
	viper.SetDefault("WEIGHT_NODE_REBOOT", 20)
	viper.SetDefault("WEIGHT_NODE_STOP", 10)
	viper.SetDefault("WEIGHT_NETWORK_DISRUPT", 20)
	viper.SetDefault("WEIGHT_SCALE_ZERO", 10)
	viper.SetDefault("NODE_STOP_DURATION", "60s")
	viper.SetDefault("SCALE_ZERO_DURATION", "30s")
	viper.SetDefault("NETWORK_DISRUPT_DELAY", "200ms")
	viper.SetDefault("METRICS_PORT", "8080")

	kubeconfigData := viper.GetString("KUBECONFIG_DATA")
	if kubeconfigData == "" {
		// Try reading from file path
		kubeconfigPath := viper.GetString("KUBECONFIG_PATH")
		if kubeconfigPath != "" {
			data, err := os.ReadFile(kubeconfigPath)
			if err != nil {
				return nil, fmt.Errorf("reading kubeconfig file: %w", err)
			}
			kubeconfigData = base64.StdEncoding.EncodeToString(data)
		}
	}

	if kubeconfigData == "" {
		return nil, fmt.Errorf("KUBECONFIG_DATA or KUBECONFIG_PATH must be set")
	}

	minInterval, err := time.ParseDuration(viper.GetString("MIN_INTERVAL"))
	if err != nil {
		return nil, fmt.Errorf("parsing MIN_INTERVAL: %w", err)
	}

	maxInterval, err := time.ParseDuration(viper.GetString("MAX_INTERVAL"))
	if err != nil {
		return nil, fmt.Errorf("parsing MAX_INTERVAL: %w", err)
	}

	nodeStopDuration, err := time.ParseDuration(viper.GetString("NODE_STOP_DURATION"))
	if err != nil {
		return nil, fmt.Errorf("parsing NODE_STOP_DURATION: %w", err)
	}

	scaleZeroDuration, err := time.ParseDuration(viper.GetString("SCALE_ZERO_DURATION"))
	if err != nil {
		return nil, fmt.Errorf("parsing SCALE_ZERO_DURATION: %w", err)
	}

	return &Config{
		KubeconfigData:      kubeconfigData,
		VictimNamespace:     viper.GetString("VICTIM_NAMESPACE"),
		ClusterAUUID:        viper.GetString("CLUSTER_A_UUID"),
		MinInterval:         minInterval,
		MaxInterval:         maxInterval,
		WeightPodKill:       viper.GetInt("WEIGHT_POD_KILL"),
		WeightNodeReboot:    viper.GetInt("WEIGHT_NODE_REBOOT"),
		WeightNodeStop:      viper.GetInt("WEIGHT_NODE_STOP"),
		WeightNetworkDisrupt: viper.GetInt("WEIGHT_NETWORK_DISRUPT"),
		WeightScaleZero:     viper.GetInt("WEIGHT_SCALE_ZERO"),
		NodeStopDuration:    nodeStopDuration,
		ScaleZeroDuration:   scaleZeroDuration,
		NetworkDisruptDelay: viper.GetString("NETWORK_DISRUPT_DELAY"),
		MetricsPort:         viper.GetString("METRICS_PORT"),
	}, nil
}
