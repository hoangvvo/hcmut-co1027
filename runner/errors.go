package runner

import "errors"

var ErrNoTestResult = errors.New("no test result")
var ErrResultMismatch = errors.New("result not matched")
