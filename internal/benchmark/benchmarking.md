# Benchmarking

To execute the benchmark, run from the root of the project:

```
make benchmark
```

If you wish to trace the execution of the benchmark, run:

```
make benchmark-trace
```

This will create a number of `*.out` files you can use to trace the execution of the benchmark with a command like the following:

```
go tool trace <name of file>
```

## Benchmarking Results

As of 2023-08-25, the benchmark results are as follows for operations.

This was executed on a Mac with an Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz

### Tree Loading

This benchmark test loads up to N amount of items and records the time it takes to load the tree.

At its worst, it takes approximately 4.7 seconds to create a single trie tree with 3 million items.

| Test Name                   | Nanoseconds per Operation | Milliseconds per Load Operation |
|-----------------------------|---------------------------|---------------------------------|
| BenchmarkLoadTree100-12     | 0.0003089                 | 3.089e-10                       |
| BenchmarkLoadTree1000-12    | 0.001288                  | 1.288e-9                        |
| BenchmarkLoadTree10000-12   | 0.01869                   | 1.869e-8                        |
| BenchmarkLoadTree100000-12  | 0.1653                    | 1.653e-7                        |
| BenchmarkLoadTree1000000-12 | 1753077726                | 1753.077726                     |
| BenchmarkLoadTree3000000-12 | 4706770540                | 4706.770540                     |

### Tree Searching

The benchmarks evaluate two different types of searches:

* Searching a single-depth tree
* Searching a two-depth tree

These benchmarks measure the time taken to find results similar to the character 'e', which should have a high frequency of occurrence across multiple different branches of the tree.

This search term is deliberately chosen to exercise an absolute-worst-case scenario.

The results are:

#### Single-Depth Tree

At its worst, it takes approximately 2.23 seconds to search a single-dimensional tree with 3 million items.

| Test Name                           | Nanoseconds per Operation | Milliseconds per Search Operation | Count of Matching Search Results |
|-------------------------------------|---------------------------|-----------------------------------|----------------------------------|
| BenchmarkSearchSingleTree100-12     | 0.0000846                 | 8.46e-11                          | 67                               |
| BenchmarkSingleSearchTree1000-12    | 0.0009899                 | 9.899e-10                         | 515                              |
| BenchmarkSingleSearchTree10000-12   | 0.009069                  | 9.069e-9                          | 4,716                            |
| BenchmarkSingleSearchTree100000-12  | 0.1308                    | 1.308e-7                          | 51,725                           |
| BenchmarkSingleSearchTree1000000-12 | 1164539029                | 1164.539029                       | 536,064                          |
| BenchmarkSingleSearchTree3000000-12 | 2952145574                | 2952.145573                       | 1,596,337                        |


#### Two-Depth Tree

The following are test cases exercising multi-tree searches.

##### Single-Character Search Phrases

At its worst, it takes approximately 8.2 seconds to search a two-dimensional tree with 3 million items using a single-character search that results in approximately 2.5 million results.

| Test Name                    | Nanoseconds per Operation | Milliseconds per Search Operation | Count of Matching Search Results |
|------------------------------|---------------------------|-----------------------------------|----------------------------------|
| BenchmarkMultiTree100-12     | 0.0001868                 | 1.868e-10                         | 75                               |
| BenchmarkMultiTree1000-12    | 0.001407                  | 1.407e-9                          | 663                              |
| BenchmarkMultiTree10000-12   | 0.02112                   | 2.112e-8                          | 6,909                            |
| BenchmarkMultiTree100000-12  | 0.2253                    | 2.253e-7                          | 77,596                           |
| BenchmarkMultiTree1000000-12 | 2761922627                | 2761.922627                       | 839,744                          |
| BenchmarkMultiTree3000000-12 | 8204723017                | 8204.723017                       | 2,498,677                        |

##### Multi-Character Search Phrase

At its worst, it takes approximately 2.0 seconds to search a two-dimensional tree with 3 million items using a multi-character search that results in 68,055 results.

| Test Name                                   | Nanoseconds per Operation | Milliseconds per Search Operation | Count of Matching Search Results |
|---------------------------------------------|---------------------------|-----------------------------------|----------------------------------|
| BenchmarkMultiTreeMultiCharPhrase100-12     | 0.0000854                 | 8.54e-11                          | 3                                |
| BenchmarkMultiTreeMultiCharPhrase1000-12    | 0.0009689                 | 9.689e-10                         | 11                               |
| BenchmarkMultiTreeMultiCharPhrase10000-12   | 0.009723                  | 9.723e-9                          | 221                              |
| BenchmarkMultiTreeMultiCharPhrase100000-12  | 0.10583                   | 1.058e-7                          | 2,059                            |
| BenchmarkMultiTreeMultiCharPhrase1000000-12 | 0.8378                    | 8.378e-7                          | 19,686                           |
| BenchmarkMultiTreeMultiCharPhrase3000000-12 | 2093077778                | 2093.077778                       | 68,055                           |

##### Multi-Character, Low Yield Phrase

At its worst, it takes approximately 2.6 seconds to search a two-dimensional tree with 3 million items using a multi-character search that results in a low result count of 2,604 results.

| Test Name                                        | Nanoseconds per Operation | Milliseconds per Search Operation | Count of Matching Search Results |
|--------------------------------------------------|---------------------------|-----------------------------------|----------------------------------|
| BenchmarkMultiTreeLowResultCountPhrase100-12     | 0.0001059                 | 1.059e-10                         | 0                                |
| BenchmarkMultiTreeLowResultCountPhrase1000-12    | 0.0006468                 | 6.468e-10                         | 0                                |
| BenchmarkMultiTreeLowResultCountPhrase10000-12   | 0.008481                  | 8.481e-9                          | 20                               |
| BenchmarkMultiTreeLowResultCountPhrase100000-12  | 0.1121                    | 1.121e-7                          | 115                              |
| BenchmarkMultiTreeLowResultCountPhrase1000000-12 | 78020647                  | 78.020647                         | 1,044                            |
| BenchmarkMultiTreeLowResultCountPhrase3000000-12 | 2605484087                | 2605.484087                       | 2,604                            |