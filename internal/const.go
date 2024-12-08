package internal

import "runtime"

var goVersion = runtime.Version()
var osName = runtime.GOOS

var userAgent = "skyway-cli/0.0.1 (" + osName + "; " + goVersion + "; Go-http-client/1.1; +https://github.com/kadoshita/skyway-cli)"
