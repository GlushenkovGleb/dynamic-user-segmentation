package repository

import (
	"dynamic-user-segmentation/pkg/model"
	"encoding/csv"
	"github.com/google/uuid"
	"io"
	"os"
)

const (
	dateTimeFormat = "2006-01-02 15:04:05"
)

type HistoryCSV struct {
	storagePath string
}

func NewHistoryCSV(storagePath string) *HistoryCSV {
	return &HistoryCSV{storagePath: storagePath}
}

// SaveEvents Takes a list of events to save in csv file
// returns name of saved file
func (h *HistoryCSV) SaveEvents(events []model.MemberEvent) (string, error) {
	rows := prepareRows(events)
	fileName := getRandomCSVName()
	filePath := h.storagePath + "/" + fileName
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	// Init Writer
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	headers := []string{"user_id, slug, status, date_time"}
	writer.Write(headers)
	err = writer.WriteAll(rows)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// UploadEvents takes file name to read from
// and writes content to dest, returns fileSize
func (h *HistoryCSV) UploadEvents(fileName string, dest io.Writer) (int, error) {
	filePath := h.storagePath + "/" + fileName
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	fileStat, _ := file.Stat()
	fileSize := int(fileStat.Size())

	_, err = file.Seek(0, 0)
	if err != nil {
		return 0, err
	}
	_, err = io.Copy(dest, file)
	if err != nil {
		return 0, err
	}

	return fileSize, nil
}

func prepareRows(events []model.MemberEvent) [][]string {
	rows := make([][]string, 0, len(events))
	for _, event := range events {
		prettyTime := event.HappenedAt.Format(dateTimeFormat)
		row := []string{event.UserId, event.Slug, event.Status, prettyTime}
		rows = append(rows, row)
	}
	return rows
}

func getRandomCSVName() string {
	return uuid.NewString() + ".csv"
}
