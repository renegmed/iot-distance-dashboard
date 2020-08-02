init-project:
	go mod init github.com/renegmed/iot-distance



# mosquitto server /usr/sbin/mosquitto should be running

run:
	go run main.go 

test:
	go test sensor-simulators/simulator_test.go 

