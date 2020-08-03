init-project:
	go mod init github.com/renegmed/iot-distance



# mosquitto server /usr/sbin/mosquitto should be running

clean-cache:
	go clean -cache -modcache -i -r

run:
	go run main.go 

test:
	go clean -cache
	go test sensor-simulators/simulator_test.go 

