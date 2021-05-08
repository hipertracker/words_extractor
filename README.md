# words_extractor

Example of words extracting in Go, Crystal, Rust and Python

Text source: 79.4MB in 20 files

- Python 3.9.5 with sorting: 11s, without sorting 10s
- Go 1.16.3 with sorting: 30s, without sorting: 20s
- Rust 1.51.0 with sorting: 1m31s, without sorting: 1m10s
- Crystal with sorting: 2m55s, without sorting: 27s

macOS 11.3.1, MacBook Pro (Retina, 15-inch, Late 2013)

```
cd words_extractor_py
python words.py

cd words_extractor_rs
cargo run

cd words_extractor_go
go run .

cd words_extractor_cr
crystal run src/fast_words_cr.cr
```
