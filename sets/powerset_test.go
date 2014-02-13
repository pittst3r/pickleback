package sets

import (
    "math"
    "pickleback/elements"
    "testing"
)

func TestFullPowerset(t *testing.T) {
    set := tSet()
    pset := set.Powerset(1, len(set.Elements))
    expectedSize := psetSize(set)
    if len(pset) > expectedSize {
        t.Error("Set too large; expected", expectedSize, "got", len(pset))
    } else if len(pset) < expectedSize {
        t.Error("Set too small; expected", expectedSize, "got", len(pset))
    }
}

func TestConstrainedPowerset(t *testing.T) {
    set := tSet()
    pset := set.Powerset(2, (len(set.Elements) - 1))
    expectedSize := psetSize(set) - len(set.Elements) - 1
    if len(pset) > expectedSize {
        t.Error("Set too large; expected", expectedSize, "got", len(pset))
    } else if len(pset) < expectedSize {
        t.Error("Set too small; expected", expectedSize, "got", len(pset))
    }
}

func psetSize(set Set) int {
    return int(math.Pow(2, float64(len(set.Elements)))) - 1
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
