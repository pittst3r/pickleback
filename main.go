package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "rakesh/sets"
    "rakesh/nodes"
    "os"
)

// Usage:
// rakesh json_input_file csv_output_file
func main() {

    // Parse the cli args
    if len(os.Args) < 3 {
        fmt.Println("Usage: rakesh json_input_file csv_output_file")
        return
    }
    infilePath := os.Args[1]
    // outfilePath := os.Args[2]

    // Parse json input
    transactionStore, jsonErr := parseJson(infilePath)
    if jsonErr != nil { return }
    // Used this to check successful parsing of json
    // for _, t := range transactionStore.Transactions[0:3] {
    //     fmt.Println(t.Id)
    //     for _, i := range t.Items {
    //         fmt.Printf("    %d: %s\n", i.Id, i.Name)
    //     }
    // }


    // Run algorithm

    // The size of our initial set of ItemSets (basically all of our items)
    // size := 1

    // Large sets have satisfied preliminary support checks
    // Map of arrays of ItemSet pointers, with the keys being the size of the underlying sets
    // largeSets := new(map[int][]*sets.ItemSet)

    // Find our first set of 1-item ItemSets
    allSingleItemSets := sets.AllSingleItemSets(transactionStore)
    candidates := new([]*sets.ItemSet)
    findItemInCandidates := func(item *nodes.Item) (candidate *sets.ItemSet, found bool) {
        found = false
        for _, c := range *candidates {
            for _, i := range c.Items {
                if i.Id == item.Id {
                    candidate = c
                    found = true
                }
            }
        }
        return
    }
    for i := range allSingleItemSets {
        if c, found := findItemInCandidates(allSingleItemSets[i].Items[0]); found {
            c.Support += int64(1)
            fmt.Printf("Incrementing support for %v to %d\n", c.Items[0].Name, c.Support, &allSingleItemSets[i], &c)
        } else {
            fmt.Printf("Adding candidate: %v\n", allSingleItemSets[i].Items[0].Name)
            *candidates = append(*candidates, allSingleItemSets[i])
        }
    }
    largeSets := new([]*sets.ItemSet)
    for _, c := range *candidates {
        if c.Support > int64(3) {
            *largeSets = append(*largeSets, c)
        }
    }
    fmt.Printf("No. of single item sets: %d\n", len(allSingleItemSets))
    fmt.Printf("No. of candidates: %d\n", len(*candidates))
    fmt.Printf("No. of large sets: %d\n", len(*largeSets))
    // fmt.Printf("Candidate 1: (%d) %v\n", (*largeSets)[1].Items[0].Id, (*largeSets)[1].Items[0].Name)

    // def find_large_sets(size, large_sets)
    //   candidates = []
    //   large_sets.each do |candidate|
    //     if found_set = candidates.find { |c| c == candidate }
    //       found_set.support += 1
    //     else
    //       candidates << ItemSet.new([*candidate.items])
    //     end
    //   end
    //   candidates.reject! { |set| set.support < min_support }
    //   candidates
    // end

    // Write output csv
}

func parseJson(filePath string) (store *sets.TransactionStore, err error) {
    rawJson, err := ioutil.ReadFile(filePath)
    store = new(sets.TransactionStore)
    json.Unmarshal(rawJson, store)
    return
}
