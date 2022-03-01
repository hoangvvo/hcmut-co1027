package conf

import (
	"os"
	"path"
)

var cwd, _ = os.Getwd()
var CASEDIR = path.Join(cwd, "cases")
