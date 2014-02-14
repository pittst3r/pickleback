package sets

import (
    "math"
    "github.com/ledbury/pickleback/elements"
)

// for each element in the set:
//     for each subset constructed so far:
//         new subset = (subset + element)

// { 1, 2, 3, 4 }

// { }
// { 1 }
// { 1, 2 } { 2 }
// { 1, 3 } { 1, 2, 3 } { 2, 3 } { 3 }
// { 1, 4 } { 1, 2, 4 } { 2, 4 } { 1, 3, 4 } { 1, 2, 3, 4 } { 2, 3, 4 } { 3, 4 } { 4 }

func Spawn(existing []*elements.Element, new *elements.Element) *Set {
    sliceLen := len(existing) + 1
    elems := make([]*elements.Element, sliceLen, sliceLen)
    for i, e := range existing {
        elems[i] = e
    }
    elems[len(elems) - 1] = new
    return &Set{Elements: elems}
}

func (set *Set) PsetSize() int {
    return int(math.Pow(2, float64(len(set.Elements))))
}

// Powerset generates the power set of the receiver except for
// the empty set, because we don't actually care about that.
//
// Additionally, Powerset stores the entire power set to the
// receiver's Subset field. Subsequent calls to Powerset will
// simply return Subset. Zero the Subset field and call Powerset
// again to re-generate the Powerset (e.g. if the set
// elements change).
//
// Min and max constrain the size of the returned sets but does
// not constrain the size of the generated and stored sets.
func (set *Set) Powerset(min, max int) []*Set {

    if set.Subsets == nil {

        sets := []*Set{new(Set)}
        for _, e := range set.Elements {
            for _, s := range sets {
                sets = append(sets, Spawn(s.Elements, e))
            }
        }

        // We don't care about the empty set
        sets = sets[1:]

        // Squirrel our power set away
        set.Subsets = sets

    }

    toReturn := []*Set{}
    for _, s := range set.Subsets {
        if s.Size() >= min && s.Size() <= max {
            toReturn = append(toReturn, s)
        }
    }
    return toReturn

}
