## IOT project on distance dashboard ##

Status: work-in-progress 

This simple application shows the communication between sensorsand LED hardwares
and the web-based dashboard through MQTT broker.

This app:

	1. Receives distance data from distance sensors
	2. Pushes the distance data to a web page using websocket
	3. Based on the distances of two sensor, this app determines where and how to position the object 
	4. Using MQTT broker, the app sends and publishes messages, i.e., sensor distances
    5. Finally, the app displays graphically car in the process of parking in, in real-time (almost)

![Car Parking Video](https://youtu.be/Hd5PZRDyrcg)
