package iotControl

import (
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/itsmeadi/tempcontrol/models"
	"github.com/itsmeadi/tempcontrol/tests/mock_tests"
	"go.uber.org/zap"
)

func TestTempControl_TempHandler(t *testing.T) {

	ctrl := gomock.NewController(t)

	logger, _ := zap.NewDevelopment()
	tc := &RoomControl{
		logger:      logger,
		desiredTemp: 22,
	}

	cli := mock_tests.NewMockClient(ctrl)

	tc.client = cli
	t.Run("WithoutPMC", func(t *testing.T) {

		payL := models.Reading{
			SensorID: "sensor-1",
			Type:     "temperature",
			Value:    25.3,
		}
		res := models.Actuator{Level: 0}
		resBytes, _ := json.Marshal(res)
		payLBytes, _ := json.Marshal(payL)
		mqttMsg := mock_tests.NewMockMessage(ctrl)
		mqttMsg.EXPECT().Payload().Return(payLBytes)

		cli.EXPECT().Publish(getRoomTopic("room"), byte(0), false, resBytes)

		tc.TempHandler(nil, mqttMsg)
	})
	t.Run("WithPMC", func(t *testing.T) {

		pmcPayload := models.Reading{
			SensorID: "sensor-1",
			Type:     "pmc",
			Value:    1,
		}
		pmcPayLBytes, _ := json.Marshal(pmcPayload)
		mqttMsg := mock_tests.NewMockMessage(ctrl)
		mqttMsg.EXPECT().Payload().Return(pmcPayLBytes)

		res := models.Actuator{Level: 10}
		resBytes, _ := json.Marshal(res)
		cli.EXPECT().Publish(getRoomTopic("room"), byte(0), false, resBytes)

		tc.PmcHandler(nil, mqttMsg)

		tempPayload := models.Reading{
			SensorID: "sensor-1",
			Type:     "temperature",
			Value:    25.3,
		}
		tempPayLBytes, _ := json.Marshal(tempPayload)
		mqttMsg = mock_tests.NewMockMessage(ctrl)
		mqttMsg.EXPECT().Payload().Return(tempPayLBytes)
		tc.TempHandler(nil, mqttMsg)
	})
}
