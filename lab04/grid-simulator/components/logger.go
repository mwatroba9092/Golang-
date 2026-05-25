package components

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type CSVLogger struct {
	logChan chan []string
	file    *os.File
	writer  *bufio.Writer
	csvW    *csv.Writer
	wg      *sync.WaitGroup
}

func NewCSVLogger(wg *sync.WaitGroup) (*CSVLogger, error) {
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		return nil, err
	}

	f, err := os.Create(filepath.Join("logs", "grid_history.csv"))
	if err != nil {
		return nil, err
	}

	w := bufio.NewWriter(f)
	cw := csv.NewWriter(w)

	cw.Write([]string{"Timestamp", "Event", "Details"})

	return &CSVLogger{
		logChan: make(chan []string, 200),
		file:    f,
		writer:  w,
		csvW:    cw,
		wg:      wg,
	}, nil
}

func (l *CSVLogger) Start(ctx context.Context) {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		for {
			select {
			case <-ctx.Done():
				l.drainChannel()
				l.Flush()
				return
			case record := <-l.logChan:
				l.csvW.Write(record)
			}
		}
	}()
}

func (l *CSVLogger) drainChannel() {
	close(l.logChan)
	for record := range l.logChan {
		l.csvW.Write(record)
	}
}

func (l *CSVLogger) LogEvent(event interface{}) {
	record := []string{
		time.Now().Format(time.RFC3339Nano),
		"GRID_EVENT",
		fmt.Sprintf("%v", event),
	}

	select {
	case l.logChan <- record:
	default:
	}
}

func (l *CSVLogger) Flush() error {
	l.csvW.Flush()
	err := l.writer.Flush()
	if err != nil {
		return err
	}
	return l.file.Close()
}