package algorithms

import (
    "github.com/ledbury/pickleback/results"
)

type Runner interface {
    Run() *results.Result
}

func RunAlgorithm(r Runner) *results.Result {
    return r.Run()
}
