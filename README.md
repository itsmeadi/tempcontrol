**Overview**

TempControl is a mqtt based app to control room temperature by changing the temperature of radiator,

**Input topics:**

For temperature `/readings/temperature`

For motion sensor `/readings/pmc`

**Output topic:**

`/actuators/<roomID>` 

`/actuators/room-1` 

**Run command**

`docker build -t iotapp . && docker run -ti iotapp -url tcp://172.17.0.1:1883`

__**Additional params__**

`  -clientID string
          (default "cid-1")`
          
   `-desiredTemp float
          (default 22)`
   
   `-url string
          (default "tcp://172.17.0.1:1883")`
   
   `-motion
            Enable Motion Sensor`       
Sample unit tests `iotControl/control_test.go`

**About**

The app reads the data for temperature and increments the openness of radiator by 10 units if the temp is lesser then desired and vice versa


The app can read data from motion sensor too, if motion sensor is enabled, the app reads data from motion sensor and saves the state, 
when the app reads a temperature reading it checks the state of motion sensor before tuning the radiator

The code can be easily tuned to add more sensors corresponding to multiple rooms
https://github.com/itsmeadi/tempcontrol/blob/master/iotControl/control.go#L137
Currently this fn returns the topic for default room, a map can be maintained that maintains an index for roomID->room actuator topic

The github repo is equipped to run automatic test cases for every push to master branch