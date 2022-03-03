package runner

import (
	"errors"
	"strconv"

	"github.com/hoangvvo/hcmut-co1027/conf"
)

var ErrNoTestResult = errors.New("no test result")
var ErrResultMismatch = errors.New("result not matched")
var ErrDeadlineExceeded = errors.New("command timed out (" + strconv.Itoa(conf.TimeoutSec) + "s)")
