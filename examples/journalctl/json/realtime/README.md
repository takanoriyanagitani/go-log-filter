# Benchmark

- golang wasmtime: go1.20.4 darwin/arm64 & wasmtime 8.0.0
- golang native:   go1.20.4 darwin/arm64
- python native:   Python 3.10.10
- rust   native:   Rust 1.69.0 & serde json 1.0.96

| type            | user  | sys | cpu  | total | jsons / second | ratio | MiB / s |
|:---------------:|:-----:|:---:|:----:|:-----:|:--------------:|:-----:|:-------:|
| golang wasmtime | 22.03 | .43 | 106% | 21.99 |  47,684        | 100%  |  47.9   |
| golang wasmtime | 21.94 | .44 | 106% | 21.93 |  47,815        | 100%  |  48.0   |
| golang wasmtime | 21.98 | .45 | 106% | 21.00 |  49,932        | 105%  |  50.1   |
| golang native   |  7.95 | .26 | 100% |  8.20 | 127,875        | 268%  | 128.4   |
| golang native   |  7.89 | .26 | 100% |  8.14 | 128,818        | 270%  | 129.4   |
| golang native   |  7.92 | .25 | 100% |  8.13 | 128,976        | 270%  | 129.5   |
| python native   |  5.96 | .13 |  98% |  6.17 | 169,947        | 356%  | 170.7   |
| python native   |  5.96 | .13 |  99% |  6.14 | 170,778        | 358%  | 171.5   |
| python native   |  5.70 | .13 |  99% |  5.86 | 178,938        | 375%  | 179.7   |
| rust   native   |  5.68 | .12 |  98% |  5.89 | 178,026        | 373%  | 178.8   |
| rust   native   |  5.67 | .13 |  99% |  5.85 | 179,244        | 376%  | 180.0   |
| rust   native   |  5.44 | .13 |  99% |  5.68 | 184,608        | 387%  | 185.4   |
