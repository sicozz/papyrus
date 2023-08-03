package dtos

import (
	"errors"
	"regexp"
	"strings"
)

type ErrDto struct {
	Message string `json:"error"`
}

func NewErrDto(msg string) (dto ErrDto) {
	dto = ErrDto{msg}
	return
}

type ValidationErrDto map[string]string

func NewValidationErrDto(errMsg string) (vErrDto ValidationErrDto, err error) {
	errors := strings.Split(errMsg, "\n")
	vErrDto = make(ValidationErrDto)
	for _, errDesc := range errors {
		k, v, err := parseValidationString(errDesc)
		if err != nil {
			return nil, err
		}
		vErrDto[k] = v
	}

	return
}

/*
* Validation error string parsing:
*
* Example string:
*		"Key: 'User.Username' Error:Field validation for 'Username' failed on the 'required' tag"
*
* The regexp matches all values inside simple quotes ('') which happen to
* have this meaning according to their possition:
*
* 0. Type.Field
* 1. Field
* 2. Failed condition #0
* 3. Failed condition #1
* ... And so on with every failed condition
 */
func parseValidationString(errDesc string) (k string, v string, err error) {
	re := regexp.MustCompile(`'([^']*)'`) // If this dosen't compile then program must stop
	matches := re.FindAllString(errDesc, -1)

	if len(matches) < 3 {
		return "", "", errors.New("bad format for validation string")
	}

	for i, m := range matches {
		m = strings.ToLower(m)
		matches[i] = strings.Trim(m, "'")
	}

	k = matches[1]
	v = strings.Join(matches[2:], ",")
	return
}
