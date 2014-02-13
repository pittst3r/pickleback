package sets

import (
    "encoding/json"
    "io/ioutil"
    "pickleback/elements"
    "sort"
)

type TransactionStore struct {
    Transactions []*Transaction
}

func (transactionStore *TransactionStore) AllUniqueElements() []*elements.Element {
    elems := []*elements.Element{}
    for _, t := range transactionStore.Transactions {
        for _, e := range t.Elements {
            elems = append(elems, e)
        }
    }
    sort.Sort(elements.ByElementId(elems))
    uniqElems := []*elements.Element{}
    for i, e := range elems {
        if i != 0 {
            if e.Id > (elems)[i-1].Id {
                uniqElems = append(uniqElems, e)
            }
        }
    }
    return uniqElems
}

func ParseJson(filePath string) (store *TransactionStore, err error) {
    rawJson, err := ioutil.ReadFile(filePath)
    store = new(TransactionStore)
    json.Unmarshal(rawJson, store)
    return
}
