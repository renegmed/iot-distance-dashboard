## IOT project on distance dashboard ##

Status: work-in-progress 

This simple application shows the communication between sensorsand LED hardwares
and the web-based dashboard through MQTT broker.

This app:

	1. Receives distance data from distance sensor
	2. Pushes the distance data to a web page using websocket
	3. Based on the distance, this app determines what color of light to turn on
	4. Using MQTT broker, the app sends and publishes message to LED lights which color to turn on
    5. Finally, the app displays, on the web page ,the turned-on light.