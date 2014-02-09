package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "rakesh/sets"
    "rakesh/elements"
    "strconv"
    "time"
    "sort"
)

// Usage:
// rakesh json_input_file csv_output_file
func main() {
    clock := time.Now()

    // Parse the cli args
    if len(os.Args) < 3 {
        fmt.Println("Usage: rakesh json_input_file csv_output_file")
        return
    }
    infilePath := os.Args[1]
    sup, _ := strconv.ParseInt(os.Args[2], 0, 0)
    minSupport := int(sup)

    // Parse json input
    transactionStore, jsonErr := parseJson(infilePath)
    if jsonErr != nil { return }

    // size represents the current size of the item sets in the computation
    size := 1

    // largeSets is where we will store our results; large because their support is large as defined by minSupport
    largeSets := make(map[int][]*sets.Set)

    // Now the fun begins...

    // Find our first set of 1-item Sets
    allSingleElemSets := sets.AllSingleSets(transactionStore)
    // Count preliminary support
    for _, s := range allSingleElemSets {
        if foundSet, ok := s.FindInSets(largeSets[size]); ok {
            foundSet.Support += 1
        } else {
            largeSets[size] = append(largeSets[size], s)
        }
    }

    // Filter out sets without minimum support
    tmpSet := make([]*sets.Set, len(largeSets[size]))
    counter := 0
    for _, c := range largeSets[size] {
        if c.Support >= minSupport {
            tmpSet[counter] = c
            counter++
        }
    }
    for i, s := range tmpSet {
        if s == nil {
            chompedSet := make([]*sets.Set, i)
            chompedSet = tmpSet[:i]
            largeSets[size] = chompedSet
            break
        }
    }
    
    // for _, s := range largeSets[size] {
    //     fmt.Printf("%v is supported by %d\n", s, s.Support)
    // }
    fmt.Printf("No. of single item sets: %d\n", len(allSingleElemSets))
    fmt.Printf("No. of large sets: %d\n", len(largeSets[size]))

    size++

    largeSets[size] = generateCandidates(size, largeSets[size-1], transactionStore.AllUniqueElements())
    
    
    // fmt.Printf("No. of large sets: %d\n", len(transactionStore.AllUniqueElements()))
    // fmt.Printf("%v\n", transactionStore.AllUniqueElements()[0])
    // fmt.Printf("No. of large sets: %d\n", len(largeSets[size]))
    // for _, s := range largeSets[size] {
    //     fmt.Printf("%v is supported by %d\n", s, s.Support)
    // }

    // Write output csv

    // Print processing time
    fmt.Printf("-> Time to run: %v seconds\n", time.Since(clock).Seconds())
}

func parseJson(filePath string) (store *sets.TransactionStore, err error) {
    rawJson, err := ioutil.ReadFile(filePath)
    store = new(sets.TransactionStore)
    json.Unmarshal(rawJson, store)
    return
}

func generateCandidates(size int, ss []*sets.Set, uniqElems []*elements.Element) []*sets.Set {
    joinedSets := new([]*sets.Set)

    // Join step
    for _, p := range ss {
        for _, q := range uniqElems {
            elems := p.Elements
            sort.Sort(elements.ByElementId(elems))
            candidate := sets.Set{Elements: elems}
            // Dupe prevention
            if elems[len(elems) - 1].Id < q.Id {
                candidate.Elements = append(candidate.Elements, q)
                *joinedSets = append(*joinedSets, &candidate)
            }
        }
    }

    // Prune step
    prunedSets := new([]*sets.Set)
    for _, s := range *joinedSets {
        subsets := s.Powerset(size-1)
        good := true
        for _, sub := range subsets {
            if _, ok := sub.FindInSets(ss); !ok {
                good = false
                break
            }
        }
        if good {
            *prunedSets = append(*prunedSets, s)
        }
    }

    return *prunedSets
}
