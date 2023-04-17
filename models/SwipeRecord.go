package models

import (
	"fmt"
	"log"
)

type SwipeRecordMessage struct {
	Username  string
	Timestamp string
	Time      float32
}

type SwipeRecordMessageInterface interface {
	convertToString(message SwipeRecordMessage)
}

func (srm *SwipeRecordMessage) ConvertToString() string {
	retVal := fmt.Sprintf("%s|%s|%f", srm.Username, srm.Timestamp, srm.Time)
	log.Print("New message converted:", retVal)
	return retVal
}
