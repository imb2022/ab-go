package ab

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/xwi88/log4go"
	"github.com/xwi88/split"

	"github.com/imb2022/ab-go/scheme"
)

var appABScheme *scheme.ABScheme
var lock sync.RWMutex

func updateScheme(abScheme *scheme.ABScheme) {
	if abScheme == nil {
		return
	}

	lock.Lock()
	defer lock.Unlock()
	appABScheme = abScheme
}

func Split(layerFlag, requestId string) (bucketNo int, experiment scheme.Experiment, err error) {
	if appABScheme == nil || len(appABScheme.Layers) == 0 {
		err = errors.New("nil scheme or layers")
		return
	}
	layers := appABScheme.Layers
	if layer, ok := layers[layerFlag]; !ok {
		err = errors.New("not exist layer for the scheme")
		return
	} else {
		if len(layer.Buckets) < 2 || len(layer.Buckets)%2 != 0 {
			err = errors.New(fmt.Sprintf("invalid buckets number in the layer %v", layerFlag))
			return
		}

		bucketNo = splitMurmurHash(requestId, layerFlag, layer.ID)
		index := searchExperimentIndexByBucket(layer.Buckets, bucketNo)
		if index == -1 {
			err = errors.New(fmt.Sprintf("missing bucket in the layer %v", layerFlag))
			return
		}
		experiment = layer.Experiments[index]
	}
	return
}

func splitMurmurHash(requestId, layerFlag string, layerId int64) int {
	hashStr := strings.Builder{}
	hashStr.WriteString(requestId)
	hashStr.WriteString(layerFlag)
	hashStr.WriteString(strconv.FormatInt(layerId, 10))

	ds := hashStr.String() + hex.EncodeToString([]byte(md5Str(hashStr.String())))
	return split.Split(ds, 100, 0)
}

func md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// limited bucket, direct for range deal, -1 missing buckets
func searchExperimentIndexByBucket(buckets []int, bucketNo int) (index int) {
	size := len(buckets)
	if bucketNo < buckets[0] || bucketNo > buckets[size-1] {
		// not over buckets lower
		return -1
	}

	for i := 0; i < size; i = i + 2 {
		if bucketNo >= buckets[i] && bucketNo <= buckets[i+1] {
			return
		}
		index++
	}

	if index == size>>1 {
		// over buckets upper
		return -1
	}
	return
}

func Run(cfg Config) {
	app := cfg.App
	mysqlCfg := cfg.MySql
	if mysqlCfg.Enable {
		err := initMySql(mysqlCfg)
		if err != nil {
			log4go.Error("[ab-go] start mysqlClient error: ", err)
			return
		}
		loadABScheme(app)
	}

	kfcCfg := cfg.Kafka
	if kfcCfg.Enable {
		err := startKafkaListener(kfcCfg, app)
		if err != nil {
			log4go.Error("[ab-go] start kafkaListener error: ", err)
			return
		}
	}

}

func Close() error {
	if kafkaListener != nil {
		log4go.Info("[ab-go] close kafkaListener")
		closeKafkaListener()
	}

	if mysqlClient != nil {
		err := closeMySql()
		if err != nil {
			log4go.Error("[ab-go] close mysqlClient error: ", err)
		}
	}
	return nil
}
