package mqtt

import (
	"encoding/json"
	"log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"math/rand"
	"time"
	"strconv"
)

func StartPublisher() {
	opts := mqtt.NewClientOptions()
	broker := os.Getenv("MQTT_BROKER_URL")
	log.Printf("URL MQTT: %s", os.Getenv("MQTT_BROKER_URL"))
	if broker == "" {
		broker = "tcp://localhost:1883"
	}
	opts.AddBroker(broker)
	opts.SetClientID("simulator-fleet-publisher")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	log.Print("Connected to MQTT broker")

	vehicleID := "B1234XYZ"
	// latitude  := -6.20000
	// longitude := 106.81667
	const offset = 0.0005
	var GEOFENCE_LAT = os.Getenv("GEOFENCE_LAT")
	var GEOFENCE_LNG = os.Getenv("GEOFENCE_LNG")

	geofenceLat, _ := strconv.ParseFloat(GEOFENCE_LAT, 64)
	geofenceLng, _ := strconv.ParseFloat(GEOFENCE_LNG, 64)

	latitude  := geofenceLat + (offset - 2*offset*rand.Float64())
	longitude := geofenceLng + (offset - 2*offset*rand.Float64())

	for {
		lat := latitude + (0.005 - 0.01*rand.Float64())
		lon := longitude + (0.005 - 0.01*rand.Float64())
		timestamp := time.Now().Unix()

		message := VehicleLocation{
			VehicleID: &vehicleID,
			Latitude:  &lat,
			Longitude: &lon,
			Timestamp: &timestamp,
		}

		payload, err := json.Marshal(message)
		topic := "/fleet/vehicle/" + vehicleID + "/location"
		if err != nil {
			log.Printf("Error marshalling message: %v\n", err)
			continue
		}

		token := client.Publish(topic, 0, false, payload)
		token.Wait()
		log.Printf("✅ Published message to topic %s: %s\n", topic, string(payload))

		time.Sleep(2 * time.Second)
	}
}