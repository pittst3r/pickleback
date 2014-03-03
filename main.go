package main

import (
    "errors"
    "fmt"
    "os"
    "strconv"
    "time"
    "github.com/ledbury/pickleback/algorithms"
    "github.com/ledbury/pickleback/results"
    "github.com/ledbury/pickleback/stores"
)

func main() {
    
    options, err := parseOptions(os.Args[1:])
    if err != nil { return }

    // Parse data store
    store := new(stores.Store)
    infilePath := options["infilePath"].(string)
    store.Read(infilePath)

    // Run algorithm
    var algorithm algorithms.Runner
    minSupport := options["minSupport"].(int)
    switch options["algorithm"] {
    case "apriori":
        algorithm = algorithms.Apriori{Data: store, MinSupport: minSupport}
    }

    // Run algorithm
    clock := time.Now()
    resultSet := algorithms.RunAlgorithm(algorithm)
    duration := time.Since(clock).Seconds()

    // Write results
    outfilePath := options["outfilePath"].(string)
    results.WriteResults(outfilePath, resultSet)

    fmt.Println("-> Run in", duration, "seconds.")

}

func parseOptions(args []string) (options map[string]interface{}, err error) {

    options = map[string]interface{}{}

    if len(args) < 4 {
        err = errors.New("Usage: pickleback algorithm minsupport infile/path.json outfile/path.csv")
        return
    }

    options["algorithm"] = args[0]

    sup, _ := strconv.ParseInt(args[1], 0, 0)
    options["minSupport"] = int(sup)

    options["infilePath"] = args[2]

    options["outfilePath"] = args[3]

    return

}
