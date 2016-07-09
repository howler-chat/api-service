// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import "github.com/Sirupsen/logrus"

func ToFields(tags map[string]string) logrus.Fields {
	result := logrus.Fields{}
	for key, value := range tags {
		result[key] = value
	}
	return result
}
