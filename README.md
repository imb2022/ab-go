# ab-go

ab sdk for go, split with `github.com/spaolacci/murmur3`

## Usage

> import by `go get -u github.com/imb2022/ab-go`

- run `func Run(cfg Config)`
- split `func Split(layerFlag, requestId string) (bucketNo int, experiment scheme.Experiment, err error) `
- close `func Close() error`
