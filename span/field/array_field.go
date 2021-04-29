package field

type ArrayField []Field

func MallocArrayField(n int) *ArrayField {

	var res = ArrayField(make([]Field, 0, n))
	return &res
}

func (f *ArrayField) Append(v Field) {
	*f = append(*f, v)
}

func (f *ArrayField) Length() int {
	return len(*f)
}

func (f *ArrayField) At(i int) (Field, error) {
	if i >= f.Length() {
		return nil, OverIndexError
	}
	return (*f)[i], nil
}
