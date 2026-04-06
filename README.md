# manajemen_armada

Backend service untuk fleet management system yang mendukung:
- Ingest data lokasi kendaraan via MQTT
- Penyimpanan data ke PostgreSQL
- Event geofence menggunakan RabbitMQ
- REST API untuk akses data kendaraan

## **Clone Repository**
Git clone or git pull from github
```bash
# Git clone
git clone https://github.com/SatriaNata/fleet_management.git

#Or

#Git pull
git pull origin main
```

## **Move into folder**
```bash
cd fleet_management
```

## **Initialization**
```bash
# create .env file base on example file
cp .env.example .env
```

## **Run with Docker**
```bash
docker compose up --build
```

sistem akan berjalan di
API -> http://localhost:8080
Check sistem berjalan http://localhost:8080/test
RabbitMQ dashboard -> http:/localhost:15672 
- username: guest
- password: guest

## **Run MQTT vehicle publisher**
 Run for send vehicle simulation cordinate every 2 seconds, to run this simulation you need to open new terminal for run this command promt
```bash
    go run cmd/publisher/main.go
```

## **Run Locally (Optional)**
1. setup dependencies
    - Postgresql, 
    - rabbitMQ, 
    - Mosquitto
2. setup environtment
```bash 
    cp .env.example .env.js 
```
3. run aplication
```bash
    go run cmd/api/main.go
```

list End point 
Test connection
- GET "/test
Get vehicle location using vehicle_id
- GET "/vehicles/:vehicle_id/location
Get vehicle history location
- GET "/vehicles/:vehicle_id/history
