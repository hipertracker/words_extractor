# words_extractor

### Info

Example of a text file parsing in several programming languages. The goal is to extract unique words from utf-8 files and save results them into separate files.

The difficulty in sorting words is due to the need to handle sorting rules according to the the different languages grammary. This is quite a complex problem that does not exist for the English language where the character set does not exceed the basic ASCII standard.

### Results

The following results are for 123 unique utf-8 Bible text files in 23 languages (used at mybible.pl site) They take 504MB. (The repo contains only a few sample files in the 'data' folder. For testing more data you could multiple files by cloning *.txt (and the associated*.yml) file under different names)

* Platform: MacOS 12.2
* Machine: MacBook Pro 16" 64GB 2TB M1Max 10 cores.

<pre>
1. Rust 1.58      = 1.14s (with sorting: 1.59s) with tokyo (previous: 1.34s, with sorting: 1.79)
2. Golang 1.17.6  = 1.34s (with sorting: 6.55s)
3. Python 3.10.2  = 2.80s (with multiprocessing)
4. Julia 1.7.1    = 4.522
5. Crystal 1.3.2  = 5.72s
6. Elixir 1.13.2  = 7.82s
7. Ruby 3.1.0     = 8.31s (with Parallel)
</pre>

### Conclusion

The new optimized Golang code version is very fast, slower than Rust but faster than other languages. Golang is the only language at the moment with full mature i18n support for arm64/M1 platform.

* Rust = the current example uses [lexical-sort](https://lib.rs/crates/lexical-sort) which is not perfect. [There is no standard mature implementation of i18n in Rust](https://www.arewewebyet.org/topics/i18n/) at the moment.

* Python = has a great implementation of [ICU](https://icu.unicode.org/related) library however it does not support arm64/M1 platform, hence I couldn't use it in this comparison.

* Ruby = same as Python, no ICU for M1.

* Elixir = same as Python, no ICU for M1.

* Julia = I couldn't find a good i18 library supporting many languages.

* Crystal = currently supports only Turkish collations. Probably because the language is young and does not have a large enough community or company behind it.

* Golang = has rules for many languages. You can see the influence of a large company and community which makes Golang a mature solution. Sorting slowed the whole task down significantly, but the result is correct (in this case I only checked the results for the Polish language)

### Kudos

[@romanatnews](https://github.com/romanatnews) (Golang example refactoring)

[@pan93412](https://github.com/pan93412) (Rust example refactoring using Tokyo runtime)
