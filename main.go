package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "github.com/ledbury/pickleback/elements"
    "github.com/ledbury/pickleback/sets"
    "github.com/ledbury/pickleback/stores"
    "sort"
    "strconv"
    "time"
    "path/filepath"
)

func DBFilename() string {
    return filepath.Join(os.TempDir(), "pickleback.db")
}

func main() {
    clock := time.Now()

    // Parse the cli args
    if len(os.Args) < 3 {
        fmt.Println("Usage: pickleback <min support> infile/path.json outfile/path.csv")
        return
    }
    sup, _ := strconv.ParseInt(os.Args[1], 0, 0)
    minSupport := int(sup)
    infilePath := os.Args[2]
    outfilePath := os.Args[3]

    stores.Initialize(DBFilename())
    defer os.Remove(DBFilename())

    // Parse json input
    transactionStore, jsonErr := sets.ParseJson(infilePath)
    if jsonErr != nil { return }

    // Get all of our transactions ready
    for _, t := range transactionStore.Transactions {

        // Sort elements in transactions
        // This is necessary for when we search the transactions for
        // candidate matches.
        sort.Sort(elements.ByElementId(t.Elements))

        stores.StoreTransaction(DBFilename(), t)
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
            foundSet.TransactionIds = append(foundSet.TransactionIds, s.TransactionIds...)
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
        if c.Support >= minSupport {
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
        for _, c := range candidates {
            tids := c.TransactionIds
            c.TransactionIds = []string{}
            for _, t := range tids {
                trans, _ := stores.FindTransaction(DBFilename(), t)
                for _, s := range trans.Powerset(1, size) {
                    if c.Eql(s) {
                        c.TransactionIds = append(c.TransactionIds, t)
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
    writeResults(&largeSets, outfilePath)

    // Print processing time
    fmt.Printf("-> Time to run: %v seconds\n", time.Since(clock).Seconds())
}

func generateCandidates(size int, largeSets []*sets.Set, singleSets []*sets.Set) []*sets.Set {
    joinedSets := []*sets.Set{}

    // Join step
    for _, p := range largeSets {
        for _, q := range singleSets {
            elems := p.Elements
            sort.Sort(elements.ByElementId(elems))
            // Dupe prevention
            if elems[len(elems) - 1].Id < q.Elements[0].Id {
                newSet := sets.Spawn(elems, q.Elements[0])
                // We add these transactions to the set so we can search them later for the
                // larger set we just made a couple lines ago. The set's transactions will
                // be accurate after that point.
                newSet.TransactionIds = []string{}
                newSet.TransactionIds = append(newSet.TransactionIds, p.TransactionIds...)
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
            if _, ok := sub.FindInSets(largeSets); !ok {
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

func writeResults(largeSets *map[int][]*sets.Set, outfilePath string) (err error) {
    outfile, _ := os.Create(outfilePath)
    csvWr := csv.NewWriter(outfile)
    resLineCount := 1
    for i := 2; len((*largeSets)[i]) > 0; i++ {
        resLineCount += len((*largeSets)[i])
    }
    resultSlice := make([][]string, resLineCount)
    counter := 0
    resultSlice[counter] = []string{"Support", "Elements"}
    counter++
    for i := 2; len((*largeSets)[i]) > 0; i++ {
        for _, s := range (*largeSets)[i] {
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
    return csvWr.WriteAll(resultSlice)
}
