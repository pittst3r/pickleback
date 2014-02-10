package sets

import (
    "rakesh/elements"
    "fmt"
    "testing"
)

func TestFullPowerset(t *testing.T) {
    set := tSet()
    pset := set.Powerset(len(set.Elements))
    expectedSize := set.PsetSize() - 1
    fmt.Printf("Expected %d sets got %d.\n", expectedSize, len(pset))
    if len(pset) > expectedSize {
        t.Error("Set too large; expected", expectedSize, "got", len(pset))
    } else if len(pset) < expectedSize {
        t.Error("Set too small; expected", expectedSize, "got", len(pset))
    }
}

func TestConstrainedPowerset(t *testing.T) {
    var set Set
    set = tSet()
    pset := set.Powerset(len(set.Elements) - 1)
    expectedSize := set.PsetSize() - 2
    fmt.Printf("Expected %d sets got %d.\n", expectedSize, len(pset))
    if len(pset) > expectedSize {
        t.Error("Set too large; expected", expectedSize, "got", len(pset))
    } else if len(pset) < expectedSize {
        t.Error("Set too small; expected", expectedSize, "got", len(pset))
    }
}

func tSet() Set {
    elems := []*elements.Element{
        &elements.Element{int64(1), "1"},
        &elements.Element{int64(2), "2"},
        &elements.Element{int64(3), "3"},
        &elements.Element{int64(4), "4"},
        &elements.Element{int64(5), "5"},
        &elements.Element{int64(6), "6"},
        &elements.Element{int64(7), "7"},
        &elements.Element{int64(8), "8"},
        &elements.Element{int64(9), "9"},
    }
    set := Set{Elements: elems}
    return set
}
