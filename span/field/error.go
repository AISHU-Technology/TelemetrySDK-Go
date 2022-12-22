package field

import "errors"

const defualtError = "TelemetrySDK-GO/Span(Log).Error: "

type StringError string

func (e StringError) Error() string {
	return string(e)
}

const (
	NilPointerError = StringError("nil pointer")
	OverIndexError  = StringError("over index")
)

func GenerateSpecificError(e error) error {
	return errors.New(defualtError + e.Error())
}
