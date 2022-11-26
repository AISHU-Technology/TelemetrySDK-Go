package common

import "devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporters/arevent/model"

// Attribute 自定义 Event Attribute 和 Trace 输出格式一致。
type Attribute struct {
	Key   string `json:"Key"`
	Value Value  `json:"Data"`
}

// NewAttribute 创建新的 Attribute 。
func NewAttribute(key string, value model.ARValue) *Attribute {
	return &Attribute{
		Key: key,
		Value: Value{
			Type: value.GetType(),
			Data: value.GetData(),
		},
	}
}

func (a *Attribute) Valid() bool {
	return a.keyNotNil() && a.keyNotCollide() && a.valueTyped()
}

// keyNotNil 校验 Attribute.Key 不为空，即有含义。
func (a *Attribute) keyNotNil() bool {
	return len(a.Key) > 0
}

// keyNotCollide 校验 Attribute.Key 不与默认值冲突。
func (a *Attribute) keyNotCollide() bool {
	switch a.Key {
	case "host":
		return false
	case "os":
		return false
	case "telemetry":
		return false
	case "service":
		return false
	default:
		return true
	}
}

// valueTyped 校验 Value.Type 是枚举类型。
func (a *Attribute) valueTyped() bool {
	switch a.Value.GetType() {
	case "BOOL":
		return true
	case "BOOLARRAY":
		return true
	case "INT":
		return true
	case "INTARRAY":
		return true
	case "FLOAT":
		return true
	case "FLOATARRAY":
		return true
	case "STRING":
		return true
	case "STRINGARRAY":
		return true
	default:
		return false
	}
}

func (a *Attribute) GetKey() string {
	return a.Key
}

func (a *Attribute) GetValue() model.ARValue {
	return a.Value
}
