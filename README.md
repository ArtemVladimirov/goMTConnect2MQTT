# goMTConnect2MQTT
Open source project that converts XML data from MTConnect Agent to MQTT. You can install this converter by one string using Docker to track a lot of CNC from different manufactures, such as Fanuc, Okuma, Haas.
The project is written in pure Go and the weight of the docker container is only <8MB

## Requirements
Firstly, check that a Docker, a MQTT broker and a MTConnect Agent are correctly installed.

## Installation
```
    docker run  \
    --init -d --name="goMTConnect2MQTT" \
    -e MQTT_HOST='192.168.1.2:1883' \
    -e MQTT_USERNAME='user' \
    -e MQTT_PASSWORD='password' \
    -e MTCONNECT_HOST='192.168.1.2:5001' \
    --restart always -t -i artvladimirov/gomtconnect2mqtt:latest
```
That's all. After that you should subscribe to MQTT topic (such as factory/sensor/YOUR_CNC_NAME/config).
Structure of MQTT topic will be in JSON.
Example of the topic:
```json
    {"Name":"YOUR_CNC_NAME","ControllerMode":"AUTOMATIC","Program":"60.9000","ProgramComment":"UNAVAILABLE","Execution":"ACTIVE","LoadS1":0,"PartCount":"12893","Time":"2022-03-22T10:31:28.3204848Z"}
```
## Limitations
This project does not convert all the data that can be obtained from MTConnect. In particular, coordinates are not transmitted. This is done because this data is rarely used in real projects to track the status of CNC. 


