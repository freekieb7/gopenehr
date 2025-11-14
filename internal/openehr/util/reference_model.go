package util

type ReferenceModel interface {
	HasModelName() bool
	Validate(path string) []ValidationError
}
