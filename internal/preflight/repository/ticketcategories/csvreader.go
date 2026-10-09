package ticketcategories

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/yoyojr404/internal/common/entity"
)

type CSVReaderRepository struct {
	filepath string
}

func NewCSVReaderRepository(filepath string) *CSVReaderRepository {
	return &CSVReaderRepository{filepath: filepath}
}

func (c *CSVReaderRepository) Read(ctx context.Context, process func(row entity.TicketCategory) error) error {
	file, err := os.Open(c.filepath)
	if err != nil {
		return fmt.Errorf("os.Open(): %w", err)
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	if _, err := csvReader.Read(); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		record, err := csvReader.Read()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return fmt.Errorf("csv parse error: %w", err)
		}

		catId, err := strconv.Atoi(record[0])
		if err != nil {
			return fmt.Errorf("strconv.Atoi(): %w", err)
		}

		item := entity.TicketCategory{
			ID:   int64(catId),
			Name: record[1],
		}

		if err := process(item); err != nil {
			return err
		}
	}

	return nil
}
