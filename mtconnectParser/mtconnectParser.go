// getXMLsendMQTT
package mtconnectParser

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/ArtemVladimirov/goMTConnect2MQTT/environment"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

//MTConnect XML Structure
type Response struct {
	Header       []XMLHeaderStruct    `xml:"Header"`
	DeviceStream []DeviceStreamStruct `xml:"Streams>DeviceStream"`
}

type XMLHeaderStruct struct {
	LastSequence int64 `xml:"lastSequence,attr"`
}

type DeviceStreamStruct struct {
	Name    string          `xml:"name,attr"`
	Events  []EventsStruct  `xml:"ComponentStream>Events"`
	Samples []SamplesStruct `xml:"ComponentStream>Samples"`
}

type EventsStruct struct {
	ControllerMode string `xml:"ControllerMode"`
	Program        string `xml:"Program"`
	Execution      string `xml:"Execution"`
	ProgramComment string `xml:"ProgramComment"`
	PartCount      string `xml:"PartCount"`
}
type SamplesStruct struct {
	Load []LoadStruct `xml:"Load"`
}

type LoadStruct struct {
	LoadNameAttr string `xml:"name,attr"`
	Value        string `xml:",chardata"`
}

//MQTT message structure
type CollectedCNCData struct {
	Name           string
	ControllerMode string
	Program        string
	ProgramComment string
	Execution      string
	LoadS1         int
	PartCount      string
	Time           time.Time
}

//Start MQTT Client and Start Parsing MTConnect
func StartMqttClient(cfg environment.Config) {

	c := make(chan string, 1)

	//Configurate MQTT Client
	opts := MQTT.NewClientOptions().AddBroker("tcp://" + cfg.MQTT_HOST).SetUsername(cfg.MQTT_USERNAME).SetPassword(cfg.MQTT_PASSWORD).SetClientID("MTCONNECT2MQTT")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetAutoReconnect(true)
	opts.OnConnect = func(c MQTT.Client) {}

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		token.Error()
	} else {
		fmt.Println("Connected to MQTT " + cfg.MQTT_HOST)
	}

	go mqttSender(client, cfg)
	<-c
}

func mqttSender(client MQTT.Client, cfg environment.Config) {
	//Create a new slice, where we store previous data, to compare it with new information
	var devicesInformationOld []CollectedCNCData

	url := "http://" + cfg.MTCONNECT_HOST + "/current"

	for {

		if xmlBytes, err := getXML(url); err != nil {
			time.Sleep(1 * time.Second)
			continue
		} else {
			var res Response
			err := xml.Unmarshal(xmlBytes, &res)
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			// If we get no information, lets start loop again
			if len(res.DeviceStream) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}

			//Create a new slice, where we store new data
			var devicesInformation []CollectedCNCData

			//Decode XML
			for i, _ := range res.DeviceStream {
				var device CollectedCNCData

				if len(res.DeviceStream[i].Name) > 0 {
					device.Name = res.DeviceStream[i].Name
				}
				if len(res.DeviceStream[i].Events) > 0 {
					for m, _ := range res.DeviceStream[i].Events {
						if len(res.DeviceStream[i].Events[m].ControllerMode) > 0 {
							device.ControllerMode = res.DeviceStream[i].Events[m].ControllerMode
						}
						if len(res.DeviceStream[i].Events[m].Program) > 0 {
							device.Program = res.DeviceStream[i].Events[m].Program
						}
						if len(res.DeviceStream[i].Events[m].Execution) > 0 {
							device.Execution = res.DeviceStream[i].Events[m].Execution
						}
						if len(res.DeviceStream[i].Events[m].ProgramComment) > 0 {
							device.ProgramComment = res.DeviceStream[i].Events[m].ProgramComment
						}
						if len(res.DeviceStream[i].Events[m].PartCount) > 0 {
							device.PartCount = res.DeviceStream[i].Events[m].PartCount
						}
					}
				}

				if len(res.DeviceStream[i].Samples) > 0 {
					for m, _ := range res.DeviceStream[i].Samples {
						if len(res.DeviceStream[i].Samples[m].Load) > 0 {
							for s, _ := range res.DeviceStream[i].Samples[m].Load {
								if res.DeviceStream[i].Samples[m].Load[s].LoadNameAttr == "S1load" {
									load, err := strconv.Atoi(res.DeviceStream[i].Samples[m].Load[s].Value)
									if err != nil {
										device.LoadS1 = 0
									} else {
										device.LoadS1 = load
									}
								}
							}
						}
					}
				}

				//
				if len(devicesInformation) == 0 {
					devicesInformation = append(devicesInformation, device)
				} else {
					repeatValue := false
					for m, _ := range devicesInformation {
						if devicesInformation[m].Name == device.Name {
							repeatValue = true
						}
					}
					if repeatValue == false {
						devicesInformation = append(devicesInformation, device)
					}
				}
			}

			//MQTT Sender
			for i, _ := range devicesInformation {
				//Check that we will send a new data
				equalFlag := false
				for m, _ := range devicesInformationOld {
					if devicesInformation[i] == devicesInformationOld[m] {
						equalFlag = true
						break
					}
				}
				if equalFlag == false {
					devicesInformation[i].Time = time.Now()
					message, err := json.Marshal(devicesInformation[i])
					if err != nil {
						fmt.Println(err)
					} else {
						sensorTopic := "factory/sensor/" + devicesInformation[i].Name + "/config"
						client.Publish(sensorTopic, 0, true, message)
					}
				}
			}
			devicesInformationOld = devicesInformation
		}
		time.Sleep(1 * time.Second)
	}
}

func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}
