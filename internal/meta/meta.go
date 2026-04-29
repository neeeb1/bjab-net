package meta

import (
	"sort"
	"time"
)

type Metadata struct {
	Title       string   `yaml:"title"`
	Date        string   `yaml:"date"`
	Slug        string   `yaml:"slug"`
	Tags        []string `yaml:"tags"`
	Description string   `yaml:"description"`
	Draft       bool     `yaml:"draft"`
}

type Dated interface {
	GetDate() string
}

func SortByDate[T Dated](items []T) ([]T, error) {
	result := make([]T, len(items))
	copy(result, items)
	var err error
	sort.SliceStable(result, func(i, j int) bool {
		t1, e1 := time.Parse("2006-01-02", result[i].GetDate())
		t2, e2 := time.Parse("2006-01-02", result[j].GetDate())
		if e1 != nil {
			err = e1
		}
		if e2 != nil {
			err = e2
		}
		return t1.After(t2)
	})
	return result, err
}
