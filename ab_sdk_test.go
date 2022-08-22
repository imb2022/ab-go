package ab

import (
	"testing"
)

func Test_SDK_Kafka(t *testing.T) {
	appNameOrFlag := "ab-go-test"
	cfg := KafkaConsumer{}
	err := StartKafkaListener(cfg, appNameOrFlag)
	if err != nil {
		t.Error(err)
		return
	}
	defer CloseKafkaListener()

	times := 5
	for i := 0; i < times; i++ {
		layerName := "Layer1"
		requestId := "1234567890"

		bucketNo, experiment, err := Split(layerName, requestId)
		if err != nil {
			t.Error(err)
			return
		}

		t.Logf("layerName:%v, requestId:%v, bucket: %v, experiment: %v",
			layerName, requestId, bucketNo, experiment)
	}
}
