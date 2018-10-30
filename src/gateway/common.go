package mqttsn

import (
	"errors"
	"io"
	"log"
	"os"
)

func InitLogger(filename string) error {
	logfile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return errors.New("failed to OpenFile.")
	}

	log.SetOutput(io.MultiWriter(logfile, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return nil
}
