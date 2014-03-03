package algorithms

import (
    "fmt"
    "sort"
    // "github.com/ledbury/pickleback/elements"
    "github.com/ledbury/pickleback/results"
    "github.com/ledbury/pickleback/sets"
    "github.com/ledbury/pickleback/stores"
)

type Apriori struct {
    Data *stores.Store
    MinSupport int
}

func (algo Apriori) Run() *results.Result {

    // Sort elements in transactions
    // This is necessary for when we search the transactions for
    // candidate matches.
    algo.Data.SortTransactionElements()

    // Now the fun begins...

    resultSet := &results.Result{}

    // Find our first set of 1-item Sets
    allSingleElemSets := stores.AllSingleSets(algo.Data)
    sort.Sort(sets.ByFirstElementId(allSingleElemSets))
    // Count preliminary support
    for _, s := range allSingleElemSets {
        if foundSet, ok := s.FindInSets(resultSet.OfSize(1)); ok {
            foundSet.Support += 1
        } else {
            s.Support = 1
            resultSet.AddSet(s)
        }
    }
    // Filter out sets without minimum support
    tmpSet := []*sets.Set{}
    for _, c := range resultSet.OfSize(1) {
        if c.Support >= algo.MinSupport {
            tmpSet = append(tmpSet, c)
        }
    }
    resultSet.ReplaceSets(1, &tmpSet)

    for size := 2; len(resultSet.OfSize(size-1)) > 0; size++ {

        candidates := generateCandidates(size, resultSet.OfSize(size-1), resultSet.OfSize(1))

        // Tally up support for candidates
        for _, t := range algo.Data.Transactions {
            for _, c := range candidates {
                if _, ok := c.FindInSets(t.Powerset(1, size)); ok {
                    c.Support += 1
                }
            }
        }

        // Filter out unsupported candidates
        supportedCandidates := []*sets.Set{}
        for _, c := range candidates {
            if c.Support >= algo.MinSupport {
                supportedCandidates = append(supportedCandidates, c)
            }
        }

        sort.Sort(sets.BySupport(supportedCandidates))
        sort.Sort(sort.Reverse(sets.BySupport(supportedCandidates)))

        resultSet.AddSets(size, supportedCandidates)

    }

    for i := 2; len(resultSet.OfSize(i)) > 0; i++ {
        fmt.Printf("Size %d count: %d\n", i, len(resultSet.OfSize(i)))
        for _, s := range resultSet.OfSize(i) {
            fmt.Printf("%v x %d\n", s, s.Support)
        }
    }

    return resultSet

}

func generateCandidates(size int, resultSet []*sets.Set, singleSets []*sets.Set) []*sets.Set {
    joinedSets := []*sets.Set{}

    // Join step
    for _, p := range resultSet {
        for _, q := range singleSets {
            elems := p.Elements
            // sort.Sort(elements.ByElementId(elems))
            // Dupe prevention
            if elems[len(elems) - 1].Id < q.Elements[0].Id {
                newSet := sets.Spawn(elems, q.Elements[0])
                joinedSets = append(joinedSets, newSet)
            }
        }
    }

    // Prune step
    prunedSets := []*sets.Set{}
    for _, s := range joinedSets {
        good := true
        sz := size - 1
        for _, sub := range s.Powerset(sz, sz) {
            if _, ok := sub.FindInSets(resultSet); !ok {
                good = false
                break
            }
        }
        if good {
            prunedSets = append(prunedSets, s)
        }
    }

    return prunedSets
}
