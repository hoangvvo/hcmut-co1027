package conf

import (
	"os"
	"path"
)

var cwd, _ = os.Getwd()
var CasesDir = path.Join(cwd, "cases")
var ArchiveDir = path.Join(cwd, "archive")
var AppURI = os.Getenv("APP_URI")
