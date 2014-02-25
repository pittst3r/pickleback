package sets

import (
    "sort"
    "github.com/ledbury/pickleback/elements"
)

type Transaction struct {
    Id string
    Set
}

func (t *Transaction) SortElements() {
    sort.Sort(elements.ByElementId(t.Elements))
}
