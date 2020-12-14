package types

// ValidateConfig is a jQuery validate
// parameters struct as defined
// https://jqueryvalidation.org/validate/
type ValidateConfig struct {
	ErrorClass string                     `json:"errorClass"`
	Rules      map[string]ValidateRule    `json:"rules"`
	Messages   map[string]ValidateMessage `json:"messages"`
}
type ValidateRule struct {
	Required bool           `json:"required"`
	Remote   ValidateRemote `json:"remote"`
}
type ValidateRemote struct {
	URL        string      `json:"url"`
	Type       string      `json:"type"`
	BeforeSend interface{} `json:"beforeSend"`
}

type ValidateMessage struct {
	Required  string `json:"required"`
	MinLength string `json:"minlength"`
}

func (jq Jquery) Validate(config ValidateConfig) {

	configMap := StructToMap(config)
	jq.Call("validate", configMap)

}

func (jq Jquery) Valid() bool {

	return jq.Call("valid").Bool()

}
