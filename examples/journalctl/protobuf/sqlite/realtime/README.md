# Benchmark

- python 1: python 3.10.10 std lib(including sqlite3) & protobuf
- golang 1: go 1.20.4 & sqlite3 & protobuf
- golang 2: go 1.20.4 & sqlite3 & protobuf, with non-0 cost abstraction

| type     | user | sys  | cpu  | total | logs / s | ratio | MiB / s | RSS      |
|:--------:|:----:|:----:|:----:|:-----:|:--------:|:-----:|:-------:|:--------:|
| python 1 | 2.42 | 1.05 |  81% | 4.27  | 246K     |   -   | 321     | 1,724 MB |
| python 1 | 2.42 | 0.98 |  81% | 4.19  | 250K     |   -   | 327     | 1,724 MB |
| python 1 | 2.41 | 0.98 |  83% | 4.07  | 258K     |   -   | 337     | 1,888 MB |
| golang 1 | 7.55 | 0.55 |  99% | 8.11  | 129K     |   -   | 169     | 18.09 MB |
| golang 1 | 7.49 | 0.44 | 104% | 7.56  | 139K     |   -   | 181     | 17.69 MB |
| golang 1 | 7.50 | 0.44 | 104% | 7.56  | 139K     |   -   | 181     | 18.27 MB |
| golang 2 | 14.2 | 0.58 | 104% | 14.1  |          |   -   |         | 18.10 MB |
| golang 2 | 14.1 | 0.58 | 104% | 14.1  |          |   -   |         | 17.55 MB |
| golang 2 | 14.1 | 0.57 | 104% | 14.1  |          |   -   |         | 17.45 MB |
