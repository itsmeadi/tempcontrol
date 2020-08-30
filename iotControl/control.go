package iotControl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/itsmeadi/tempcontrol/models"
	"go.uber.org/zap"
)

type RoomControl struct {
	client      mqtt.Client
	logger      *zap.Logger
	pmc         sync.Map
	pmcEnable   bool
	actLevel    sync.Map
	desiredTemp float64
}

func NewRoomControl(clientID, urlStr string, logger *zap.Logger, desiredTemp float64, pmcEnable bool) *RoomControl {

	logger = logger.Named("RoomControl")
	uri, err := url.Parse(urlStr)
	if err != nil {
		logger.Fatal("Invalid url", zap.Error(err))
	}
	opts := createClientOptions(clientID, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		logger.Fatal("Error while Init", zap.Error(err), zap.String("url", urlStr))
	}

	if pmcEnable {
		logger.Info("Motion Sensor enabled, need both motion sensor and temperature reading to work")
	}
	logger.Info("Connected to broker successfully")
	return &RoomControl{client: client, logger: logger, pmc: sync.Map{}, actLevel: sync.Map{}, desiredTemp: desiredTemp, pmcEnable: pmcEnable}
}
func (tc *RoomControl) TempHandler(client mqtt.Client, msg mqtt.Message) {

	tc.logger.Info("Received new Temperature Reading")
	var rd models.Reading
	payLoad := msg.Payload()
	err := json.Unmarshal(payLoad, &rd)
	if err != nil {
		tc.logger.Error("Error unmarshall Payload", zap.Error(err))
		return
	}
	tc.logger.Debug("Received New Temperature Reading", zap.String("Reading", string(payLoad)))
	roomID := getRoomID(rd)
	if motionSensorState, ok := tc.pmc.Load(roomID); ok && motionSensorState.(bool) || !tc.pmcEnable {
		actVal := tc.getNewActState(roomID, rd.Value)
		tc.setActuator(roomID, actVal)
		return
	}
	if tc.pmcEnable {
		tc.logger.Info("Motion Sensor shows offline, setting radiator to 0")
	}
	//set Thermistor state to 0 in case motion detector detects nothing
	tc.setActuator(roomID, 0)
}

func (tc *RoomControl) getNewActState(roomID string, temp float64) int {

	var newValue, prevValue int
	prevValueIt, ok := tc.actLevel.Load(roomID)
	if ok {
		prevValue, ok = prevValueIt.(int)
		if !ok {
			tc.logger.Error("Invalid value set to prevValueIt")
		}
	}
	if temp < tc.desiredTemp {
		newValue = prevValue + 10
	} else {
		newValue = prevValue - 10
	}
	newValue = max(0, newValue)
	newValue = min(100, newValue)
	return newValue
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func (tc *RoomControl) PmcHandler(client mqtt.Client, msg mqtt.Message) {

	if !tc.pmcEnable {
		return
	}
	payLoad := msg.Payload()
	var rd models.Reading
	err := json.Unmarshal(payLoad, &rd)
	if err != nil {
		tc.logger.Error("Error unmarshall Payload", zap.Error(err))
		return
	}
	tc.logger.Debug("Received New Motion Sensor Reading", zap.String("Reading", string(payLoad)))

	var val bool
	if rd.Value == 1 {
		val = true
	} else {
		val = false
		tc.setActuator(getRoomID(rd), 0)
	}

	tc.pmc.Store(getRoomID(rd), val)
}

func (tc *RoomControl) setActuator(roomID string, val int) {

	act := models.Actuator{Level: val}

	actJSON, err := json.Marshal(act)
	if err != nil {
		tc.logger.Error("Unable to Marshal models.Actuator json", zap.Error(err))
		return
	}
	roomActTopic := getRoomTopic(roomID)

	tc.logger.Info("Writing value to Actuator", zap.String("topic", roomActTopic), zap.Int("value", val))
	tc.actLevel.Store(roomID, val)
	tc.client.Publish(roomActTopic, 0, false, actJSON)
}

//Find the topic for Room models.Actuator for room with `id`
func getRoomTopic(id string) string {
	return "/actuators/" + id
}

//Get room id from the models.Reading msg
func getRoomID(rd models.Reading) string {
	return "room-1"
}

func (tc *RoomControl) InitSubscribers(ctx context.Context) {

	tc.client.Subscribe("/readings/temperature", 0, tc.TempHandler)
	tc.client.Subscribe("/readings/pmc", 0, tc.PmcHandler)
}
func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s", uri.Scheme, uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
}
