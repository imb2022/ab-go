package ab

import (
	"fmt"
	"testing"
	"time"

	"github.com/Shopify/sarama"

	"github.com/imb2022/ab-go/scheme"
)

func Test_SDK(t *testing.T) {
	var enableMysql bool
	var enableKafka bool
	var testWithNormalRun bool // true run test with Run; else run test with init kfk each time
	var initialOffset int64

	testWithNormalRun = true
	// enableKafka = true
	enableMysql = true
	initialOffset = sarama.OffsetOldest

	kafkaCfg := KafkaConsumer{
		Enable:        enableKafka,
		Version:       sarama.MaxVersion.String(),
		Brokers:       []string{"10.14.41.57:9092", "10.14.41.58:9092"},
		Initial:       initialOffset,
		MaxWaitTime:   time.Second,
		SessionTimout: 20 * time.Second,
		Topic:         "d-ab-data",
		Errors:        true,
	}

	mysqlCfg := MySql{
		Enable:             enableMysql,
		Addr:               "10.14.41.52:3306",
		User:               "root",
		Password:           "Password123",
		DBName:             "ab",
		Collation:          "utf8mb4_general_ci",
		MaxOpenConnections: 4,
		MaxIdleConnections: 4,
		ConnMaxLifetime:    time.Minute * 5,
		Debug:              true,
		TLS:                false,
		ParseTime:          true,
		Timeout:            time.Second * 5,
		ReadTimeout:        time.Second * 5,
		WriteTimeout:       time.Second * 5,
		LOC:                "Asia/Shanghai",
	}

	cfg := Config{
		App:   "",
		Kafka: kafkaCfg,
		MySql: mysqlCfg,
	}

	appNameLayerMapping := map[string][]string{
		"fe_test":  {"layer_fe_test"},
		"engine10": {"update_test"},
		"engine6":  {"layer_sdk_test_1", "layer_sdk_test_2"},
	}

	times := 5
	for i := 0; i < times; i++ {
		requestId := fmt.Sprintf("1234567890%v", i)
		for appNameOrFlag, layerFlags := range appNameLayerMapping {
			if testWithNormalRun {
				splitWithRun(t, i, cfg, layerFlags, appNameOrFlag, requestId)
			} else {
				splitWithKafka(t, i, cfg, layerFlags, appNameOrFlag, requestId)
			}
		}
	}
}

func splitWithKafka(t *testing.T, i int, cfg Config, layerFlags []string, appNameOrFlag, requestId string) {
	cfg.App = appNameOrFlag

	for _, layerFlag := range layerFlags {
		bucketNo, experiment, err := consumerAndSplitWithKafka(t, cfg.Kafka, appNameOrFlag, layerFlag, requestId)
		if err == nil {
			t.Logf("[splitWithKafka] index[%v], app[%v], layer[%v], requestId[%v], bucket[%v], experiment: %+v",
				i, appNameOrFlag, layerFlag, requestId, bucketNo, experiment)
		}
	}
}

func splitWithRun(t *testing.T, i int, cfg Config, layerFlags []string, appNameOrFlag, requestId string) {
	cfg.App = appNameOrFlag
	Run(cfg)
	time.Sleep(time.Millisecond * 150)

	for _, layerFlag := range layerFlags {
		bucketNo, experiment, err := consumerAndSplit(layerFlag, requestId)
		// bucketNo, experiment, err := consumerAndSplitWithKafka(t, cfg.Kafka, appNameOrFlag, layerFlag, requestId)
		if err == nil {
			t.Logf("[split] index[%v], app[%v], layer[%v], requestId[%v], bucket[%v], experiment: %+v",
				i, appNameOrFlag, layerFlag, requestId, bucketNo, experiment)
		}
	}

	err := Close()
	if err != nil {
	}
}

func consumerAndSplitWithKafka(t *testing.T, cfg KafkaConsumer, appNameOrFlag, layerFlag, requestId string) (bucketNo int, experiment scheme.Experiment, err error) {
	err = startKafkaListener(cfg, appNameOrFlag)
	if err != nil {
		t.Error(err)
		return
	}
	defer closeKafkaListener()
	time.Sleep(time.Millisecond * 150)

	bucketNo, experiment, err = Split(layerFlag, requestId)
	if err != nil {
		return
	}
	return
}

func consumerAndSplit(layerFlag, requestId string) (bucketNo int, experiment scheme.Experiment, err error) {
	bucketNo, experiment, err = Split(layerFlag, requestId)
	if err != nil {
		return
	}
	return
}
