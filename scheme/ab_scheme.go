package scheme

type ABScheme struct {
	APP    string            `json:"app"`    // value: app_flag
	Layers map[string]*Layer `json:"layers"` // key: layer_flag
}

type Layer struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Buckets     []int        `json:"buckets"`
	Experiments []Experiment `json:"experiments"`
}

type Experiment struct {
	Name      string   `json:"name"`
	Tag       string   `json:"tag"`
	Param     string   `json:"param"`
	ParamVal  string   `json:"param_val"`
	Bucket    []int    `json:"bucket"`
	Whitelist []string `json:"whitelist"`
}
