package models

import ()

////////////

type validationResult struct {
	Valid  bool
	Errors []string
}

func newValidationResult() validationResult {
	return validationResult{Valid: true, Errors: []string{}}
}

func (x *validationResult) AddError(err string) {
	x.Valid = false
	x.Errors = append(x.Errors, err)
}
