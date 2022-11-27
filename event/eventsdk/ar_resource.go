package eventsdk

import (
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/event/version"
	"encoding/json"
)

// resource 自定义 event resource 和 Trace 输出格式一致。
type resource struct {
	SchemaURL     string                 `json:"SchemaURL"`
	AttributesMap map[string]interface{} `json:"Attributes"`
}

// newResource 创建新的 *resource 。
func newResource() *resource {
	r := &resource{
		SchemaURL:     version.SchemaURL,
		AttributesMap: defaultAttributes(),
	}
	return r
}

func (r *resource) GetSchemaURL() string {
	return r.SchemaURL
}

func (r *resource) GetAttributes() map[string]interface{} {
	return r.AttributesMap
}

func (r *resource) private() {}

// MarshalJSON 只输出 AttributesMap 。
func (r *resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.AttributesMap)
}

// UnmarshalJSON AttributesMap 变回 resource 。
func (r *resource) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &r.AttributesMap)
	return err
}
