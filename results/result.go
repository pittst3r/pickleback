package results

import (
    "encoding/csv"
    "fmt"
    "os"
    "github.com/ledbury/pickleback/sets"
)

type Result map[int][]*sets.Set

func (r *Result) OfSize(s int) []*sets.Set {
    return (*r)[s]
}

func (r *Result) AddSet(s *sets.Set) {
    (*r)[s.Size()] = append((*r)[s.Size()], s)
}

func (r *Result) AddSets(size int, ss []*sets.Set) {
    if (*r)[size] == nil {
        (*r)[size] = []*sets.Set{}
    }
    (*r)[size] = append((*r)[size], ss...)
}

func WriteResults(outfilePath string, results *Result) (error) {
    outfile, _ := os.Create(outfilePath)
    csvWr := csv.NewWriter(outfile)
    resLineCount := 1
    for i := 2; len((*results)[i]) > 0; i++ {
        resLineCount += len((*results)[i])
    }
    resultSlice := make([][]string, resLineCount)
    counter := 0
    resultSlice[counter] = []string{"Support", "Elements"}
    counter++
    for i := 2; len((*results)[i]) > 0; i++ {
        for _, s := range (*results)[i] {
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
