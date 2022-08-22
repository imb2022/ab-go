package ab

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/xwi88/log4go"

	"github.com/imb2022/ab-go/scheme"
)

// realtime update
var (
	kafkaListener  sarama.Consumer
	kafkaPartition sarama.PartitionConsumer
	kafkaQuit      chan struct{}
)

func StartKafkaListener(cfg KafkaConsumer, appNameOrFlag string) (err error) {
	config := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(cfg.Version)
	if err != nil {
		return err
	}

	config.Version = version
	brokers := cfg.Brokers
	topic := cfg.Topic
	config.Consumer.Return.Errors = cfg.Errors
	if kafkaListener, err = sarama.NewConsumer(brokers, config); err != nil {
		return err
	}

	kafkaPartition, err = kafkaListener.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	// consume errors
	go func() {
		for err := range kafkaPartition.Errors() {
			log4go.Error("[kafkaListener] %s[%d] err: %s\n", err.Topic, err.Partition, err.Error())
		}
	}()
	kafkaQuit = make(chan struct{}, 1)
	go KafkaPartitionHandler(appNameOrFlag)
	return nil
}

func CloseKafkaListener() {
	kafkaQuit <- struct{}{}
	if kafkaPartition != nil {
		if err := kafkaPartition.Close(); err != nil {
			log4go.Error("[kafkaListener] close err: %s", err.Error())
		}
	}
	if kafkaListener != nil {
		if err := kafkaListener.Close(); err != nil {
			log4go.Error("[kafkaListener] close err: %s", err.Error())
		}
	}
}

func KafkaPartitionHandler(appNameOrFlag string) {
Loop:
	for {
		select {
		case msg, ok := <-kafkaPartition.Messages():
			if ok && msg != nil && string(msg.Key) == appNameOrFlag {
				newScheme := &scheme.ABScheme{}
				err := json.Unmarshal(msg.Value, newScheme)
				if err != nil {
					log4go.Error("[kafkaListener] %s[%d][%d], err: %v", msg.Topic, msg.Partition, msg.Offset, err)
					return
				}
				log4go.Info("[kafkaListener] %s[%d][%d]: %v", msg.Topic, msg.Partition, msg.Offset, newScheme)
				if newScheme != nil && len(newScheme.APP) > 0 && len(newScheme.Layers) > 0 {
					UpdateScheme(newScheme)
				}
			} else {
				log4go.Error("[kafkaListener] consume error: %v", msg)
			}
		case <-kafkaQuit:
			break Loop
		}
	}
	log4go.Info("[kafkaListener] stop listen")
}
