# words_extractor

|   |   |   |   |   |
|--- |--- |--- |--- |--- |
|   |   |   |   |   |
|   |   |   |   |   |
|   |   |   |   |   |

Example of a text file parsing in several programming languages

MacOS 12.2
Rust 1.58.1
MBP 16" 64GB 2TB M1Max 10 cores
Tested on 123 files (504MB)

Results:

1. Rust with waitgroup 1.58.1 -> 0.3521 s
2. Ruby 3.1 with Parallel -> 2.0542 s
3. Python 3.10.2 with multiprocessing -> 2.9403 s
4. Crystal 1.3.2 with channels ->  6.0035 s
5. Go 1.18beta1 with waitgroup -> 7.2166 s

For more files (just got into data/pl/* and clone files several times)

Total files: 936
Total size: 2745 MB

1. Rust:     7.5436 s
2. Python:  15.3445 s
3. Crystal: 32.3369 s
4. Ruby:    40.9438 s
5. Golang:  73.0321 s
