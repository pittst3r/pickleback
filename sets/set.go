package sets

import (
    "rakesh/elements"
    "fmt"
)

type Set struct {
    Elements []*elements.Element
    Support int
}

func (set Set) Size() int {
    return len(set.Elements)
}

func (set Set) String() string {
    return fmt.Sprint(set.Elements)
}

// Returns all items from transactions in transaction store
func AllSingleSets(transactionStore *TransactionStore) []*Set {
    ss := new([]*Set)
    for _, t := range transactionStore.Transactions {
        for _, e := range t.Elements {
            newSet := Set{Elements: []*elements.Element{e}}
            *ss = append(*ss, &newSet)
        }
    }
    return *ss
}

// Returns true if receiver is found in slice, returns the matched set
func (set *Set) FindInSets(ss []*Set) (foundSet *Set, ok bool) {
    ok = false
    for s := range ss {
        for i := range ss[s].Elements {
            for p := range set.Elements {
                if set.Elements[p].Id == ss[s].Elements[i].Id {
                    ok = true
                    foundSet = ss[s]
                }
            }
        }
    }
    return
}
