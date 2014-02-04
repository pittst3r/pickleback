package sets

import (
    "rakesh/nodes"
)

type ItemSet struct {
    Items []*nodes.Item
    Support int64
}

func (s *ItemSet) size() int {
    return len(s.Items)
}

// Returns all items from transactions in transaction store
func AllSingleItemSets(transactionStore *TransactionStore) []*ItemSet {
    itemSets := new([]*ItemSet)
    for _, t := range transactionStore.Transactions {
        newSet := new(ItemSet)
        newSet.Items = t.Items
        *itemSets = append(*itemSets, newSet)
    }
    return *itemSets
}
