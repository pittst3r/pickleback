package sets

import (
    "rakesh/elements"
    "math"
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

// Given the above algorithm:
// Let s be any set in S
// Let p(s) be the parent of s
// Let z(s) be the size of s
// Let t(s, n) be the subsets of size n of s
// Let the children of s be defined as
//   any sets beginning with s
// Let c(s) be the children of s
// t(s, z(s)) = c(p(s))


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

func (set *Set) Powerset(limit int) []*Set {
    sets := new([]*Set)
    *sets = append(*sets, new(Set))
    for _, e := range set.Elements {
        for _, s := range *sets {
            if s.Size() < limit {
                *sets = append(*sets, Spawn(s.Elements, e))
            }
        }
    }
    // We don't care about the empty set.
    return (*sets)[1:]
}
