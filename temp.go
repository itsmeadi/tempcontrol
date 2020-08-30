package main

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/itsmeadi/tempcontrol/models"
	"go.uber.org/zap"
)

func (tc *TempControl) TempHandler(client mqtt.Client, msg mqtt.Message) {

	tc.logger.Info("Received new models.Reading")
	var rd models.Reading
	err := json.Unmarshal(msg.Payload(), &rd)
	if err != nil {
		tc.logger.Error("Error unmarshall Payload", zap.Error(err))
		return
	}
	tc.logger.Debug("Received New models.Reading", zap.String("models.Reading", string(msg.Payload())))
	roomID := getRoomID(rd)
	if motionSensorState, ok := tc.pmc.Load(roomID); ok {
		if motionSensorState.(bool) {
			actVal := tc.getNewActState(roomID, rd.Value)
			tc.setActuator(roomID, actVal)
			return
		}
	}

	//set Thermistor state to 0 in case motion detector detects nothing
	tc.setActuator(roomID, 0)
}

func (tc *TempControl) getNewActState(roomID string, temp float64) int {

	var newValue, prevValue int
	prevValueIt, ok := tc.actLevel.Load(roomID)
	if !ok {
		prevValue, ok = prevValueIt.(int)
		if !ok {
			tc.logger.Error("Invalid value set to prevValueIt")
		}
	}
	if temp > tc.desiredTemp {
		newValue = prevValue + 10
	} else {
		newValue = prevValue - 10
	}
	return newValue
}

func (tc *TempControl) PmcHandler(client mqtt.Client, msg mqtt.Message) {

	var rd models.Reading
	err := json.Unmarshal(msg.Payload(), &rd)
	if err != nil {
		tc.logger.Error("Error unmarshall Payload", zap.Error(err))
		return
	}

	var val bool
	if rd.Value == 1 {
		val = true
	} else {
		val = false
	}
	tc.pmc.Store(getRoomID(rd), val)
}

func (tc *TempControl) setActuator(roomID string, val int) {

	act := models.Actuator{Level: val}

	actJSON, err := json.Marshal(act)
	if err != nil {
		tc.logger.Error("Unable to Marshal models.Actuator json", zap.Error(err))
		return
	}
	roomActTopic := getRoomTopic(roomID)

	tc.logger.Info("Writing value to models.Actuator", zap.String("topic", roomActTopic), zap.Int("value", val))
	tc.actLevel.Store(roomID, val)
	tc.client.Publish(roomActTopic, 0, false, actJSON)
}

//Find the topic for Room models.Actuator for room with `id`
func getRoomTopic(id string) string {
	return "/models.Actuators/room-1"
}

//Get room id from the models.Reading msg
func getRoomID(rd models.Reading) string {
	return "room-1"
}
