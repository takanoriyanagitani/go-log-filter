# Benchmark

- golang wasmtime: go1.20.4 darwin/arm64 & wasmtime 8.0.0
- golang native:   go1.20.4 darwin/arm64
- python native:   Python 3.10.10

| type            | user  | sys | cpu  | total | jsons / second | ratio |
|:---------------:|:-----:|:---:|:----:|:-----:|:--------------:|:-----:|
| golang wasmtime | 22.03 | .43 | 106% | 21.99 |  47,684        | 100%  |
| golang wasmtime | 21.94 | .44 | 106% | 21.93 |  47,815        | 100%  |
| golang wasmtime | 21.98 | .45 | 106% | 21.00 |  49,932        | 105%  |
| golang native   |  7.95 | .26 | 100% |  8.20 | 127,875        | 268%  |
| golang native   |  7.89 | .26 | 100% |  8.14 | 128,818        | 270%  |
| golang native   |  7.92 | .25 | 100% |  8.13 | 128,976        | 270%  |
| python native   |  5.96 | .13 |  98% |  6.17 | 169,947        | 356%  |
| python native   |  5.96 | .13 |  99% |  6.14 | 170,778        | 358%  |
| python native   |  5.70 | .13 |  99% |  5.86 | 178,938        | 375%  |
