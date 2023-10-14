package sanitizer

type Validator interface {
	Struct(s interface{}) error
}
