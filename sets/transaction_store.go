package sets

import (
    "rakesh/elements"
    "sort"
)

type TransactionStore struct {
    Transactions []*Transaction
}

func (transactionStore *TransactionStore) AllUniqueElements() []*elements.Element {
    elems := new([]*elements.Element)
    for _, t := range transactionStore.Transactions {
        for _, e := range t.Elements {
            *elems = append(*elems, e)
        }
    }
    sort.Sort(elements.ByElementId(*elems))
    uniqElems := new([]*elements.Element)
    for i, e := range *elems {
        if i != 0 {
            if e.Id > (*elems)[i-1].Id {
                *uniqElems = append(*uniqElems, e)
            }
        }
    }
    // for i := range uniqElems {
    //     if uniqElems[i] == nil {
    //         *elems = uniqElems[:i]
    //         break
    //     }
    // }
    return *uniqElems
}
