# Pickleback

An implementation of the Apriori Algorithm, published by Rakesh Agrawal and Ramakrishnan Srikant in 1994, in Go.

I've made some modifications that increase performance by over 50x.

## Usage

```shell
pickleback <minimum support> /path/to/transactions.json /path/to/outfile.csv

# e.g.
pickleback 5 sample_transactions.json tmp/results.csv
```
