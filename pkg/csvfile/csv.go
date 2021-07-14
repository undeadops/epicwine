package csvfile

import (
	"encoding/csv"
	"epicwine/pkg/wine"
	"fmt"
	"os"
	"strconv"
	"sync"
)

type CsvFile struct {
	file string
	mu   *sync.RWMutex
}

func NewCsvFile(file string) *CsvFile {
	return &CsvFile{
		file: file,
		mu:   &sync.RWMutex{},
	}
}

func (c *CsvFile) VerifyCsv() error {
	// check if file exists
	if _, err := os.Stat(c.file); os.IsNotExist(err) {
		return fmt.Errorf("csv file %q not found", c.file)
	}
	// Read Lock
	c.mu.RLock()
	defer c.mu.RUnlock()
	// Can we open the file in read/write mode
	f, err := os.OpenFile(c.file, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("unable to open and validate csv file, check permissions")
	}
	// Close our file
	defer f.Close()

	// Check if file is a csv file
	return nil
}

func (c *CsvFile) WriteRecord(wine wine.WineRequestInput) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Convert our stuct to a []string to be writtne to csv
	var row []string
	row = append(row, strconv.Itoa(wine.ID))
	row = append(row, wine.Country)
	row = append(row, wine.Description)
	row = append(row, wine.Designation)
	row = append(row, wine.Points)
	row = append(row, wine.Price)
	row = append(row, wine.Province)
	row = append(row, wine.Region1)
	row = append(row, wine.Region2)
	row = append(row, wine.TasterName)
	row = append(row, wine.TasterTwitterHandle)
	row = append(row, wine.Title)
	row = append(row, wine.Variety)
	row = append(row, wine.Winery)

	// Open CSV for writing
	f, err := os.OpenFile(c.file, os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("unable to open csv")
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	err = csvWriter.Write(row)
	if err != nil {
		return fmt.Errorf("unable to write to csv")
	}
	// Flush to disk
	csvWriter.Flush()
	return nil
}

func (c *CsvFile) GetCount() (int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	f, err := os.Open(c.file)
	if err != nil {
		return 0, fmt.Errorf("unable to open csv")
	}
	defer f.Close()

	// Using csvReader still, so as to count CSV only lines
	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("unable to read csv")
	}
	// Count number of rows minus header
	count := len(rows[1:])
	return count, nil
}

func (c *CsvFile) GetRecords(limit, offset int) ([]map[string]interface{}, error) {
	// Setup Return data variable
	var data []map[string]interface{}
	// Read Locks
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Open our CSV file
	f, err := os.Open(c.file)
	if err != nil {
		return nil, fmt.Errorf("unable to open csv")
	}
	defer f.Close()

	// loop
	n := 0

	// Employ our CSV Reader
	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll() // `rows` is of type [][]string
	if err != nil {
		return nil, fmt.Errorf("unable to read csv")
	}

	// index position for title
	title := find(rows[0], "title")

	for i, row := range rows[offset:] {
		// skip header
		if i == 0 {
			continue
		}

		n++
		w := map[string]interface{}{
			"id":    row[0],
			"title": row[title],
		}
		// add to my returned data
		data = append(data, w)

		if n == limit {
			break
		}
	}

	return data, nil
}

// GetRecord - Pull one line from CSV file
func (c *CsvFile) GetRecord(id int) (map[string]interface{}, error) {
	// Read Locks
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Open our CSV file
	f, err := os.Open(c.file)
	if err != nil {
		return nil, fmt.Errorf("unable to open csv")
	}
	defer f.Close()

	// increment id by 1 because of header
	id = id + 1

	// Employ our CSV Reader
	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to read csv")
	}

	// Check Length of Row
	if len(rows) <= id {
		return nil, fmt.Errorf("index out of range")
	}

	// index position for title
	title := find(rows[0], "title")

	// +1 because header is at 0, but ids start at 0
	row := rows[id]

	data := map[string]interface{}{
		"id":    row[0],
		"title": row[title],
	}

	return data, nil
}

// find - used to find `title` index
func find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}

	}
	return len(a)
}
