package field

type OverIndexError struct{}

func (e OverIndexError) Error() string {
	return "over index"
}
