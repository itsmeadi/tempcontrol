**Overview**

TempControl is a mqtt based app to control room temperature by changing the temperature of radiator

**Input topics:**

For temperature `/readings/temperature`
```json
{
     "sensorID": "sensor-1",
     "type": "temperature",
     "value": 25.3
}
```

For motion sensor `/readings/pmc`
```json
{
     "sensorID": "sensor-1",
     "type": "pmc",
     "value": 1
}
```
The value can be set 0 or 1 depending on motion

**Output topic:**

`/actuators/<roomID>` 

`/actuators/room-1` 


**Run command**

`docker build -t iotapp . && docker run -ti iotapp -url tcp://172.17.0.1:1883`

**Additional params**

`  -clientID string
          (default "cid-1")`
          
   `-desiredTemp float
          (default 22)`
   
   `-url string
          (default "tcp://172.17.0.1:1883")`
   
   `-motion
            Enable Motion Sensor`    
               

**About**

The app reads the data for temperature and increments the openness of radiator by 10 units if the temp is lesser then desired, and vice versa


The app can read data from motion sensor too, if motion sensor is enabled, the app saves the state, 
whenever the app reads a temperature reading it checks the state of motion sensor before tuning the radiator

The code can be easily tuned to add more sensors corresponding to multiple rooms
https://github.com/itsmeadi/tempcontrol/blob/master/iotControl/control.go#L142
Currently this fn returns the topic for default room, a map can be maintained that maintains an index for roomID->room actuator topic

The github repo is equipped to run automatic test cases for every push to master branch


Sample unit tests `iotControl/control_test.go`
