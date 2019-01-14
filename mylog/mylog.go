package mylog

import (
	"log"
	"io"
	"sync"
	"os"
)

type logger struct {
	filename string
	*log.Logger
}

var (
	once sync.Once
)

func GetInstance() *logger {
	logger := createLogger("gyw.log")
	return logger
}

func createLogger(fname string) *logger {
	f, err := os.OpenFile("gywlog", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		
	wrt := io.MultiWriter(os.Stdout, f)	
	return &logger{
        filename: fname,
		Logger:   log.New(wrt, "gyw logs:", log.Ldate | log.Ltime | log.Lshortfile),
	}
}
