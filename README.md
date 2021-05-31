# words_extractor

Example of text file parsing in Python, Golang, Elixir, Rust, Crystal and Julia

Text source: 79.4MB in 30 files

- Rust 1.51.0 (parallel) with sorting: 2s, without sorting: 1s
- Python 3.9.5 (parallel) with sorting: 3.42s, without sorting: 2.47
- Crystal 1.0.0 (parallel) with sorting: 5.78s, withouts sorting: 2.48
- Go 1.16.4 (parallel) with sorting: 6.41s, without sorting: 3.75s
- Rust 1.51.0 with sorting: 7s, without sorting: 5s (no parallelism)
- Python 3.9.5 with sorting: 10s, without sorting 8.32s (no parallelism)
- Julia 1.6.1 (8 threads) 9s, (1 thread) 10.3s without sorting
- Crystal 1.0.0 with sorting: 13s, without sorting: 7s (no parallelism)
- Go 1.16.4 with sorting: 21s, without sorting: 11s (no parallelism)
- Elixir 1.12 (parallel) with sorting: 33s (without release build)

macOS 11.3.1, MacBook Pro (Retina, 15-inch, Late 2013)

Python
```bash
cd words_extractor_py
python words.py
```

Rust
```
cd words_extractor_rs
cargo build --release
target/release/words_extractor_rs

Golang
```
cd words_extractor_go
make build
GOGC=2000 ./main

Crystal
```
cd words_extractor_cr
crystal build --release -Dpreview_mt src/fast_words_cr.cr -o main
CRYSTAL_WORKES=8 ./main
```

Julia
```
JULIA_NUM_THREADS=8 julia src/words_extractor_jl.jl
```

Elixir
```
cd words_extractor_ex
mix run -e "WordsExtractor.run"
```

## Running Python

1. Install the latest Python 3.9.5
2. Create venv and dependencies

```bash
cd words_extractor_py
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

3. Run the code

```bash
python words_parallel.py
```
