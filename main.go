package main

import (
    "encoding/csv"
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
    sup, _ := strconv.ParseInt(os.Args[1], 0, 0)
    minSupport := int(sup)
    infilePath := os.Args[2]
    outfilePath := os.Args[3]

    // Parse json input
    transactionStore, jsonErr := parseJson(infilePath)
    if jsonErr != nil { return }

    // Sort elements in transactions
    // This is necessary for when we search the transactions for
    // candidate matches. It's most efficient to do them all at once
    // right here rather than sorting each transaction over and
    // over again later.
    for _, t := range transactionStore.Transactions {
        sort.Sort(elements.ByElementId(t.Elements))
    }

    // Now the fun begins...

    // largeSets is where we will store our results; large because their support is large as defined by minSupport
    largeSets := make(map[int][]*sets.Set)

    // Find our first set of 1-item Sets
    allSingleElemSets := sets.AllSingleSets(transactionStore)
    sort.Sort(sets.ByFirstElementId(allSingleElemSets))
    // Count preliminary support
    for _, s := range allSingleElemSets {
        if foundSet, ok := s.FindInSets(largeSets[1]); ok {
            foundSet.Support += 1
        } else {
            s.Support = 1
            largeSets[1] = append(largeSets[1], s)
        }
    }
    // Filter out sets without minimum support
    tmpSet := make([]*sets.Set, len(largeSets[1]))
    counter := 0
    for _, c := range largeSets[1] {
        if c.Support >= (minSupport * 2) {
            tmpSet[counter] = c
            counter++
        }
    }
    for i, s := range tmpSet {
        if s == nil {
            chompedSet := make([]*sets.Set, i)
            chompedSet = tmpSet[:i]
            largeSets[1] = chompedSet
            break
        }
    }

    for size := 2; len(largeSets[size-1]) > 0; size++ {

        candidates := generateCandidates(size, largeSets[size-1], largeSets[1])

        // Tally up support for candidates
        for _, t := range transactionStore.Transactions {
            for _, c := range candidates {
                for _, s := range t.Powerset(size) {
                    if c.Eql(s) {
                        c.Support += 1
                        break
                    }
                }
            }
        }

        // Filter out unsupported candidates
        supportedCandidates := make([]*sets.Set, len(candidates))
        counter := 0
        for _, c := range candidates {
            if c.Support >= minSupport {
                supportedCandidates[counter] = c
                counter++
            }
        }
        for i, c := range supportedCandidates {
            if c == nil {
                supportedCandidates = supportedCandidates[:i]
                break
            }
        }

        sort.Sort(sets.BySupport(supportedCandidates))
        sort.Sort(sort.Reverse(sets.BySupport(supportedCandidates)))

        largeSets[size] = supportedCandidates

    }

    for i := 2; len(largeSets[i]) > 0; i++ {
        fmt.Printf("Size %d count: %d\n", i, len(largeSets[i]))
        for _, s := range largeSets[i] {
            fmt.Printf("%v x %d\n", s, s.Support)
        }
    }

    // Write output csv
    outfile, _ := os.Create(outfilePath)
    csvWr := csv.NewWriter(outfile)
    resLineCount := 1
    for i := 2; len(largeSets[i]) > 0; i++ {
        resLineCount += len(largeSets[i])
    }
    resultSlice := make([][]string, resLineCount)
    counter = 0
    resultSlice[counter] = []string{"Support", "Elements"}
    counter++
    for i := 2; len(largeSets[i]) > 0; i++ {
        for _, s := range largeSets[i] {
            resultSlice[counter] = make([]string, (len(s.Elements) + 1))
            for c := range resultSlice[counter] {
                if c == 0 {
                    resultSlice[counter][c] = fmt.Sprintf("%d", s.Support)
                } else {
                    resultSlice[counter][c] = s.Elements[c - 1].Name
                }
            }
            counter++
        }
    }
    csvWr.WriteAll(resultSlice)

    // Print processing time
    fmt.Printf("-> Time to run: %v seconds\n", time.Since(clock).Seconds())
}

func parseJson(filePath string) (store *sets.TransactionStore, err error) {
    rawJson, err := ioutil.ReadFile(filePath)
    store = new(sets.TransactionStore)
    json.Unmarshal(rawJson, store)
    return
}

func generateCandidates(size int, largeSets []*sets.Set, singleSets []*sets.Set) []*sets.Set {
    joinedSets := new([]*sets.Set)

    // Join step
    for _, p := range largeSets {
        for _, q := range singleSets {
            elems := p.Elements
            sort.Sort(elements.ByElementId(elems))
            // Dupe prevention
            if elems[len(elems) - 1].Id < q.Elements[0].Id {
                *joinedSets = append(*joinedSets, sets.Spawn(elems, q.Elements[0]))
            }
        }
    }

    // Prune step
    prunedSets := new([]*sets.Set)
    for _, s := range *joinedSets {
        good := true
        for _, sub := range s.Powerset(size-1) {
            if _, ok := sub.FindInSets(largeSets); !ok {
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
