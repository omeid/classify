package classify

import (
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

type Term struct {
	Word     string
	Category string
	Weight   float64
}

type Terms []Term

var InvalidCSV = errors.New("Invalid CSV. Expects: Value,Categroy,Weight")

func FromCSV(source io.Reader) (Terms, error) {

	r := csv.NewReader(source)
	headers, err := r.Read()
	if err != nil {
		return nil, err
	}

	if len(headers) < 3 ||
		strings.ToLower(headers[0]) != "word" ||
		strings.ToLower(headers[1]) != "category" ||
		strings.ToLower(headers[2]) != "weight" {
		return nil, InvalidCSV
	}
	var terms Terms
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) < 3 {
			return nil, InvalidCSV
		}

		weight, err := strconv.ParseFloat(record[2], 63)
		if err != nil {
			return nil, err
		}

		terms = append(terms, Term{Word: record[0], Category: record[1], Weight: weight})
	}
	return terms, nil
}

type Results map[string]float64

type Result struct {
	Category string
	Weight   float64
}

func (t Terms) Analyse(text io.Reader) (Results, error) {

	content, err := ioutil.ReadAll(text)
	if err != nil {
		return nil, err
	}

	rs := make(chan Result)
	var wg sync.WaitGroup

	wg.Add(len(t))

	for _, t := range t {
		go func(content string, t Term) {
			defer wg.Done()
			count := float64(strings.Count(content, t.Word))
			if count != 0 {
				rs <- Result{t.Category, t.Weight * count}
			}
		}(string(content), t)
	}

	go func() {
		wg.Wait()
		close(rs)
	}()

	results := make(map[string]float64)
	for t := range rs {
		if v, ok := results[t.Category]; ok {
			results[t.Category] = t.Weight + v
		} else {
			results[t.Category] = t.Weight
		}
	}

	return results, nil
}
