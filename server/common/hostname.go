package common

import (
	"fmt"
	"go-link-shortener/server/globals"
)

func GetHostname() string {
	return GetHostnameWithPort(globals.Config.Dashboard.Showport)
}

func GetHostnameWithPort(showport bool) string {
	hostname := globals.Config.Server.Hostname
	if showport {
		hostname = fmt.Sprintf("%s:%d", globals.Config.Server.Hostname, globals.Config.Server.Port)
	}
	return hostname
}
