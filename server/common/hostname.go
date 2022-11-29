package common

import (
	"fmt"
	"go-link-shortener/server/globals"
)

func GetHostname() string {
	hostname := globals.Config.Server.Hostname
	if globals.Config.Server.Port != 443 {
		hostname = fmt.Sprintf("%s:%d", globals.Config.Server.Hostname, globals.Config.Server.Port)
	}
	return hostname
}
