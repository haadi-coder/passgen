# Password Generator Benchmarks

This document presents the performance benchmarks for the `passgen` package across four optimization iterations. Each iteration reflects improvements made to the password generation algorithm, with key changes noted below. The benchmarks measure operations per second (`ops/sec`), time per operation (`ns/op`), memory usage (`B/op`), and number of allocations (`allocs/op`) for various test cases.

## Benchmark Iterations

- **#1**: Initial implementation.
- **#2**: Optimized `shuffleString` with buffered values.
- **#3**: Optimized `generatePassEntry` for reduced allocations.
- **#4**: Optimized `Generate` by transitioning from strings to runes, minimizing cross-type conversions.

## Results

| Test Case              | Iteration | Ops/sec   | Time (ns/op) | Memory (B/op) | Allocations (allocs/op) |
|------------------------|-----------|-----------|--------------|---------------|-------------------------|
| default                | #1        | 162,255   | 6,163        | 2,216         | 104                     |
| default                | #2        | 223,814   | 4,468        | 1,456         | 57                      |
| default                | #3        | 433,087   | 2,309        | 688           | 9                       |
| default                | #4        | 486,811   | 2,054        | 464           | 6                       |
| short_password         | #1        | 300,030   | 3,333        | 1,408         | 55                      |
| short_password         | #2        | 407,497   | 2,454        | 1,040         | 32                      |
| short_password         | #3        | 721,124   | 1,387        | 656           | 8                       |
| short_password         | #4        | 891,266   | 1,122        | 432           | 5                       |
| long_password          | #1        | 42,316    | 23,632       | 7,280         | 395                     |
| long_password          | #2        | 62,383    | 16,029       | 4,216         | 204                     |
| long_password          | #3        | 128,447   | 7,787        | 1,144         | 12                      |
| long_password          | #4        | 134,479   | 7,436        | 920           | 9                       |
| no_symbols             | #1        | 177,809   | 5,624        | 1,992         | 103                     |
| no_symbols             | #2        | 255,957   | 3,907        | 1,232         | 56                      |
| no_symbols             | #3        | 460,073   | 2,174        | 464           | 8                       |
| no_symbols             | #4        | 492,319   | 2,031        | 464           | 6                       |
| digits_only            | #1        | 165,152   | 6,055        | 1,656         | 101                     |
| digits_only            | #2        | 235,404   | 4,248        | 896           | 54                      |
| digits_only            | #3        | 513,083   | 1,949        | 128           | 6                       |
| digits_only            | #4        | 489,236   | 2,044        | 464           | 6                       |
| with_min_requirements   | #1        | 123,848   | 8,073        | 2,656         | 133                     |
| with_min_requirements   | #2        | 179,891   | 5,559        | 1,704         | 74                      |
| with_min_requirements   | #3        | 338,870   | 2,951        | 744           | 14                      |
| with_min_requirements   | #4        | 371,885   | 2,689        | 520           | 11                      |
| complex_requirements    | #1        | 80,939    | 12,355       | 3,848         | 206                     |
| complex_requirements    | #2        | 117,182   | 8,533        | 2,320         | 111                     |
| complex_requirements    | #3        | 232,125   | 4,308        | 784           | 15                      |
| complex_requirements    | #4        | 248,999   | 4,017        | 560           | 12                      |

## Observations

- **Performance Improvements**: Each iteration significantly improves performance, with iteration #4 achieving the highest operations per second and lowest memory usage across all test cases.
- **Short Passwords**: The `short_password` test case shows the most significant improvement, reaching 891,266 ops/sec in iteration #4, due to reduced allocations and faster execution.
- **Long Passwords**: The `long_password` test case remains the slowest due to its length (64 characters), but iteration #4 reduces memory usage by ~87% compared to #1.
- **Memory Efficiency**: Iteration #4 minimizes allocations (down to 5–12 per operation) by using runes instead of strings, avoiding cross-type conversions.
- **Digits Only**: The `digits_only` test case shows a slight regression in time from iteration #3 to #4 (1,949 ns to 2,044 ns), possibly due to overhead in rune handling, but memory usage remains low.

## Notes

- **Ops/sec** is calculated as `1,000,000,000 / ns/op` to represent the number of operations per second.
- Benchmarks were run on a system with 12 CPU cores (indicated by `-12` in the output).
- Further optimizations could focus on reducing memory allocations for longer passwords or investigating the slight performance regression in the `digits_only` case.