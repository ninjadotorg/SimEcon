package tracker

import (
	"encoding/csv"
	"os"
)

type Tracker struct {
}

var tracker *Tracker

func GetTrackerInstance() *Tracker {
	if tracker != nil {
		return tracker
	}
	tracker = &Tracker{}
	return tracker
}

func (t *Tracker) OpenFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (t *Tracker) WriteToCSV(fileName string, record []string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	defer file.Close()
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	err = writer.Write(record)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}
