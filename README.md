# words_extractor

Example of words extracting in Go, Crystal, Rust and Python

Text source: 79.4MB in 30 files

- Python 3.9.5 (parallel) with sorting: 4.01s, without sorting: 2.47
- Go 1.16.4 (parallel) with sorting: 7.32s, without sorting: 4.06s
- Python 3.9.5 with sorting: 10s, without sorting 8.32s
- Go 1.16.4 with sorting: 21s, without sorting: 11s
- Rust 1.51.0 with sorting: 1m31s, without sorting: 1m10s
- Crystal 1.0.0 with sorting: 2m55s, without sorting: 27s

macOS 11.3.1, MacBook Pro (Retina, 15-inch, Late 2013)

```
cd words_extractor_py
python words.py

cd words_extractor_rs
cargo run

cd words_extractor_go
make run

cd words_extractor_cr
crystal run src/fast_words_cr.cr
```

## Running Python

1. Install the latest Python 3.9.5
2. Create venv and dependencies

```
cd words_extractor_py
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

3. Run the code

```
python words_parallel.py
```
