package field

type StringError string

func (e StringError) Error() string {
	return string(e)
}

const (
	NilPointerError = StringError("nil pointer")
	OverIndexError  = StringError("over index")
)
