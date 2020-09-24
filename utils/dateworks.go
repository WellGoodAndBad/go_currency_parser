package utils

import (
	"time"
)

func DatesParse() []string{
	layout := "2006-01-02"
	var datesParse []string
	for i:=-30; 0 > i; i++ {
		date:=time.Now().AddDate(0, 0, i)
		datesParse = append(datesParse, date.Format(layout))
	}
	return datesParse
}
