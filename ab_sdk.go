package ab

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/xwi88/split"

	"github.com/imb2022/ab-go/scheme"
)

var appABScheme *scheme.ABScheme
var lock sync.RWMutex

func UpdateScheme(abScheme *scheme.ABScheme) {
	if abScheme == nil {
		return
	}

	lock.Lock()
	defer lock.Unlock()
	appABScheme = abScheme
}

func Split(abScheme *scheme.ABScheme, layerName, requestId string) (bucketNo int, experiment scheme.Experiment, err error) {
	if abScheme == nil || len(abScheme.Layers) == 0 {
		err = errors.New("nil scheme or layers")
		return
	}
	layers := abScheme.Layers
	if layer, ok := layers[layerName]; !ok {
		err = errors.New("not exist layer for the scheme")
		return
	} else {
		if len(layer.Buckets)%2 != 0 {
			err = errors.New(fmt.Sprintf("invalid buckets number in the layer %v", layerName))
			return
		}

		bucketNo = splitMurmurHash(requestId, layerName, layer.ID)
		index := searchExperimentIndexByBucket(layer.Buckets, bucketNo)
		experiment = layer.Experiments[index]
	}
	return
}

func splitMurmurHash(requestId, layerName string, layerId int64) int {
	hashStr := strings.Builder{}
	hashStr.WriteString(requestId)
	hashStr.WriteString(layerName)
	hashStr.WriteString(strconv.FormatInt(layerId, 10))

	ds := hashStr.String() + hex.EncodeToString([]byte(md5Str(hashStr.String())))
	return split.Split(ds, 100, 0)
}

func md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// limited bucket, direct for range deal
func searchExperimentIndexByBucket(buckets []int, bucketNo int) (index int) {
	for i := 0; i < len(buckets); i = i + 2 {
		if bucketNo <= buckets[i+1] {
			return
		}
		index++
	}
	return
}
