package mqtt

import (
	"log"
	"os"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"fleet-management/internal/geofence"
	"fleet-management/internal/rabbitmq"
	"fleet-management/internal/service"
	"time"
)

type VehicleLocation struct {
	VehicleID *string  `json:"vehicle_id"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Timestamp *int64   `json:"timestamp"`
}

func StartSubscriber() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(os.Getenv("MQTT_BROKER_URL"))
	opts.SetClientID("fleet-subscriber")
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)

	// client := mqtt.NewClient(opts)
	// if token := client.Connect(); token.Wait() && token.Error() != nil {
	// 	log.Fatal(token.Error())
	// }

	var client mqtt.Client
	var err error

	for i := 0; i < 10; i++ {
		client = mqtt.NewClient(opts)

		token := client.Connect()
		if token.Wait() && token.Error() == nil {
			log.Println("✅ Connected to MQTT broker")
			break
		}

		err = token.Error()
		log.Printf("MQTT not ready (%v), retrying... (%d/10)", err, i+1)

		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("❌ Failed to connect MQTT: %v", err)
	}



	log.Println("Connected to MQTT broker")

	topic := "/fleet/vehicle/+/location"
	client.Subscribe(topic, 0, messageHandler)
}

func messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message on topic: %s\n", string(msg.Topic()))
	log.Printf("Received message payload: %s\n", string(msg.Payload()))

	var loc VehicleLocation
	err := json.Unmarshal(msg.Payload(), &loc)
	if err != nil {
		log.Printf("Invalid JSON Format: %v", err)
		return
	}

	isValidFormat := isValidationFormatLoc(loc)
	if isValidFormat != true {
		log.Printf("❌ invalid message format: %v", isValidFormat)
		return 
	}

	err = service.InsertLocation(*loc.VehicleID, *loc.Latitude, *loc.Longitude, *loc.Timestamp)
	if err != nil {
		log.Printf("❌ Error inserting location into database: %v\n", err)
	} else {
		log.Printf("✅ Location Vehicle successfully insert to DB: %s", *loc.VehicleID)
	}

	if geofence.IsInsideGeofence(*loc.Latitude, *loc.Longitude) {
		log.Printf("Send data to rabbit MQ - Vehicle entered geofence: %s", string(*loc.VehicleID))
		rabbitmq.PublishGeofenceAlert(
			*loc.VehicleID,
			*loc.Latitude,
			*loc.Longitude,
			*loc.Timestamp,
		)
	} else {
		log.Printf("❌ Vehicle is outside geofence: %s", string(*loc.VehicleID))
	}
}

func isValidationFormatLoc(loc VehicleLocation) interface{} {
	if loc.VehicleID == nil || *loc.VehicleID == "" {
		return "missing vehicle_id"
	}

	if (loc.Latitude == nil || *loc.Latitude == 0) || (loc.Longitude == nil || *loc.Longitude == 0) {
		return "missing latitude or longitude"
	}

	if (loc.Timestamp == nil || *loc.Timestamp == 0) {
		return "missing timestamp"
	}

	return true
}