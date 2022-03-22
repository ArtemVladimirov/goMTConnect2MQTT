// goMTConnect2MQTT project main.go
package main

import (
	"fmt"

	"github.com/ArtemVladimirov/goMTConnect2MQTT/environment"
	"github.com/ArtemVladimirov/goMTConnect2MQTT/mtconnectParser"
)

func main() {
	//Get environment vars
	cfg, err := environment.GetEnvVars()
	if err != nil {
		fmt.Println(err)
		return
	}
	//Starting MQTT Parser
	mtconnectParser.StartMqttClient(cfg)
}
