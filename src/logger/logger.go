package logger

import (
	log "github.com/sirupsen/logrus"
)

var Log = log.New()

func InitLog(){
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
	Log.Formatter = new(prefixed.TextFormatter)
	Log.Level = log.DebugLevel
}

//func qeq(){
//	Log.Formatter
//}
//var ioWriter = os.Stderr
//var RootLogger = log.New(ioWriter, "[", 1)


//
//func Println(args...interface{}){
//	RootLogger.Println(args)
//}