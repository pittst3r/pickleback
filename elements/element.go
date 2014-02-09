package elements

import (
    "fmt"
    "sort"
)

type Element struct {
    Id int64
    Name string
}

// Implementation of sort.Interface for elements
type ByElementId []*Element

func (x ByElementId) Len() int           { return len(x) }
func (x ByElementId) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x ByElementId) Less(i, j int) bool { return x[i].Id < x[j].Id }

// Elements must be sorted in ascending order by id
func SearchElements(elements *[]*Element, element *Element) int {
    return sort.Search(len(*elements), func(i int) bool { return (*elements)[i].Id >= element.Id })
}

func (e Element) String() string {
    return fmt.Sprintf("%s", e.Name)
}
