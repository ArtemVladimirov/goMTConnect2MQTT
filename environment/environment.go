// environment
package environment

import "github.com/caarlos0/env"

// Environment variables
type Config struct {
	MQTT_HOST      string `env:"MQTT_HOST" envDefault:"192.168.1.37:1883"`
	MQTT_USERNAME  string `env:"MQTT_USERNAME" envDefault:""`
	MQTT_PASSWORD  string `env:"MQTT_PASSWORD" envDefault:""`
	MTCONNECT_HOST string `env:"MTCONNECT_HOST" envDefault:"192.168.1.37:5001"`
}

//Getting Environment variables
func GetEnvVars() (data Config, err error) {
	data = Config{}
	if err := env.Parse(&data); err != nil {
		return data, err
	} else {
		return data, nil
	}
}
