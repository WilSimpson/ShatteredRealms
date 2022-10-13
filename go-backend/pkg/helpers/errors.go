package helpers

import log "github.com/sirupsen/logrus"

func Check(err error, errContext string) {
	if err != nil {
		log.Fatalf("%s: %v", errContext, err)
	}
}
