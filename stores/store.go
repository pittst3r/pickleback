package stores

import (
    "encoding/json"
    "io/ioutil"
    "github.com/ledbury/pickleback/elements"
    "github.com/ledbury/pickleback/sets"
)

type Store struct {
    Transactions []*sets.Transaction
}

func (s *Store) SortTransactionElements() {
    for _, t := range s.Transactions {
        t.SortElements()
    }
}

func AllSingleSets(s *Store) []*sets.Set {
    ss := []*sets.Set{}
    for _, t := range s.Transactions {
        for _, e := range t.Elements {
            newSet := sets.Set{Elements: []*elements.Element{e}}
            ss = append(ss, &newSet)
        }
    }
    return ss
}

func (s *Store) Read(infilePath string) (error) {
    rawJson, err := ioutil.ReadFile(infilePath)
    json.Unmarshal(rawJson, s)
    return err
}
