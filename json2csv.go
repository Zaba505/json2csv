package json2csv

import (
  "io"
  "encoding/csv"
  "encoding/json"
  "errors"
  "fmt"
)

// Option
type Option func(*config)

// MapFieldToColumn
func MapFieldToColumn(fieldName string, columnIndex int) Option {
  return func(cfg *config) {
    cfg.mapFieldsToColumns[fieldName] = columnIndex
  }
}

// SkipColumnTitles
func SkipColumnTitles() Option {
  return func(cfg *config) {
    cfg.skipColumnTitles = true
  }
}

type config struct {
  skipColumnTitles bool
  mapFieldsToColumns map[string]int
}

// Convert maps JSON objects to CSV.
func Convert(w io.Writer, r io.Reader, opts ...Option) error {
  cfg := buildConfig(opts...)

  objReader := readJSON(r)
  csvWriter := csv.NewWriter(w)
  defer csvWriter.Flush()

  row := make([]string, len(cfg.mapFieldsToColumns))
  if !cfg.skipColumnTitles {
		for title, idx := range cfg.mapFieldsToColumns {
			row[idx-1] = title
		}

    if err := csvWriter.Write(row); err != nil {
      return err
    }
	}

	for {
    obj := objReader.nextObject()
    if obj == nil {
      return csvWriter.Error()
    }

		err := validateJSON(cfg.mapFieldsToColumns, obj)
		if err != nil {
			return err
		}

		for field, val := range obj {
      idx := cfg.mapFieldsToColumns[field]
      row[idx-1] = fmt.Sprintf("%v", val)
		}

    if err := csvWriter.Write(row); err != nil {
      return err
    }
	}
}

func buildConfig(opts ...Option) (cfg config) {
  cfg.skipColumnTitles = false
  cfg.mapFieldsToColumns = make(map[string]int)

  for _, opt := range opts {
    opt(&cfg)
  }

  return
}

type reader struct {
  objChan chan map[string]interface{}
}

func newReader() *reader {
  return &reader{
    objChan: make(chan map[string]interface{}, 1),
  }
}

func (r *reader) read(src io.Reader) {
  var data []map[string]interface{}

	dec := json.NewDecoder(src)

  err := dec.Decode(&data)
  if err != nil {
    panic(err)
  }
  defer close(r.objChan)

  for _, obj := range data {
    r.objChan <- obj
  }
}

func (r *reader) nextObject() map[string]interface{} {
  return <-r.objChan
}

func readJSON(src io.Reader) *reader {
	r := newReader()
  go r.read(src)

  return r
}

func validateJSON(fieldColMap map[string]int, data map[string]interface{}) error {
	for k, v := range data {
		if _, exists := fieldColMap[k]; !exists {
			return errors.New(fmt.Sprintf("json2xlsx: missing column mapping for json field: %s", k))
		}

		err := validateJSONValue(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateJSONValue(v interface{}) error {
	switch v.(type) {
	case map[string]interface{}:
		return errors.New("json2xlsx: nested objects are not supported values")
	default:
		return nil
	}
}
