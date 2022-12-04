package constants

import "time"

const (
	DefaultSecretName string = "uptimerobot-secret"
	SecretConfigKey   string = "apiKey.yaml"

	MonitorDefaultRequeueTime time.Duration = time.Minute * 5
)
