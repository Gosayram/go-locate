# Benchmark Results

Generated on `Sat Jun 14 13:16:22 +05 2025`

## Go Benchmarks
| Test | Iterations | Time/op | Memory/op | Allocs/op |
|------|------------|---------|-----------|-----------|
| BenchmarkAdvancedFuzzyMatch/Short-12 | 191004991 | 6.320 | 0 | 0 |
| BenchmarkAdvancedFuzzyMatch/Medium-12 | 77873618 | 15.39 | 0 | 0 |
| BenchmarkAdvancedFuzzyMatch/Long-12 | 70379856 | 16.75 | 0 | 0 |
| BenchmarkAdvancedFuzzyMatch/Complex-12 | 100000000 | 11.88 | 0 | 0 |
| BenchmarkAdvancedFuzzyMatch/NoMatch-12 | 77193580 | 15.76 | 0 | 0 |
| BenchmarkAdvancedShouldExclude-12 | 1000000 | 1158 | 0 | 0 |
| BenchmarkAdvancedMatches-12 | 521324 | 2239 | 24 | 2 |
| BenchmarkSearchRoots-12 | 1000000000 | 1.177 | 0 | 0 |
| BenchmarkConcurrentSearch-12 | 1297 | 940171 | 156606 | 1522 |
| BenchmarkPatternMatching/Wildcard-12 | 4038162 | 297.4 | 0 | 0 |
| BenchmarkPatternMatching/Exact-12 | 10104397 | 116.3 | 0 | 0 |
| BenchmarkPatternMatching/Complex-12 | 4846804 | 245.8 | 0 | 0 |
| BenchmarkPatternMatching/MultipleWildcards-12 | 1552206 | 772.1 | 0 | 0 |
| BenchmarkMemoryAllocation-12 | 153493 | 7764 | 3569 | 104 |
| BenchmarkLargeDirectorySearch-12 | 153 | 7801491 | 1042155 | 8484 |
| BenchmarkExtensionFiltering-12 | 1927033 | 623.5 | 0 | 0 |
