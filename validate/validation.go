// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package validate

import (
	stdError "errors"
	"fmt"
	"regexp"

	"github.com/asaskevich/govalidator"
	"golang.org/x/net/context"
)

var whiteSpace = regexp.MustCompile(`^\s*$`)

type Validation interface {
	Validate(context.Context) error
}

func IsValidId(id string) error {
	if !govalidator.StringLength(id, 10, 10) {
		return stdError.New("Must be 10 characters long")
	}
	return nil
}

// Validates the passed text is considered valid for a message
func IsMessageText(text string) error {
	if !govalidator.StringLength(text, 1, 300) {
		return stdError.New(fmt.Sprintf("Must be between '%d' and '%d' characters long", 1, 300))
	}
	// Returns true if the text passed contains only whitespace characters \n, \t, ' '
	if whiteSpace.Match(text) {
		return stdError.New(fmt.Sprintf("A message with only white space is not allowed"))
	}
	return nil
}
