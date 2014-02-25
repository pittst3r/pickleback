package sets

import (
    "fmt"
    "github.com/ledbury/pickleback/elements"
)

type Set struct {
    Elements []*elements.Element
    Support int
}

type BySupport []*Set
func (x BySupport) Len() int           { return len(x) }
func (x BySupport) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x BySupport) Less(i, j int) bool { return x[i].Support < x[j].Support }

type ByFirstElementId []*Set
func (x ByFirstElementId) Len() int           { return len(x) }
func (x ByFirstElementId) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x ByFirstElementId) Less(i, j int) bool { return x[i].Elements[0].Id < x[j].Elements[0].Id }

func (set Set) Size() int {
    return len(set.Elements)
}

func (set Set) String() string {
    return fmt.Sprint(set.Elements)
}

// Returns true if receiver is found in slice, returns the matched set
func (set *Set) FindInSets(ss []*Set) (foundSet *Set, ok bool) {
    ok = false
    for _, s := range ss {
        if set.Eql(s) {
            foundSet = s
            ok = true
            break
        }
    }
    return
}

func (set *Set) Eql(s *Set) bool {
    if len(set.Elements) != len(s.Elements) {
        return false
    }
    for i := range set.Elements {
        if s.Elements[i].Id != set.Elements[i].Id {
            return false
        }
    }
    return true
}
