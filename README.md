# words_extractor

### Info

Example of a text file parsing in several programming languages. The goal is to extract unique words from utf-8 files and save results them into separate files.

The difficulty in sorting words is due to the need to handle sorting rules according to the the different languages grammary. This is quite a complex problem that does not exist for the English language where the character set does not exceed the basic ASCII standard.

### Results

The following results are for 123 unique utf-8 Bible text files in 23 languages (used at mybible.pl site) They take 504MB. (The repo contains only a few sample files in the 'data' folder. For testing more data you could multiple files by cloning *.txt (and the associated*.yml) file under different names)

* Platform: macOS 12.6
* Machine: MacBook Pro 16" 64GB 2TB M1Max 10 cores.

<pre>
1. Rust 1.63           = 1.15s, sorting: 1.64s
2. Golang 1.22.0       = 1.40s, sorting with collations: 1.71s
3. Python 3.12.2       = 5.69, sorting with collations: 6.04s
4. Crystal 1.5.1       = 5.61s
5. Elixir 1.14.0       = 7.34s
6. Julia 1.8.1         = 12.13s, sorting: 12.22s
7. Ruby 3.3.0          = 12.63s, sorting with collations: 22.00s
</pre>

### Conclusion

The new optimized Golang code version is very fast, slower than Rust but faster than other languages. Golang is the only language at the moment with full mature i18n support for arm64/M1 platform.

* Rust = the current example uses [lexical-sort](https://lib.rs/crates/lexical-sort) which is not perfect. [There is no standard mature implementation of i18n in Rust](https://www.arewewebyet.org/topics/i18n/) at the moment.

* Python = has a great implementation of [ICU](https://icu.unicode.org/related) library however it does not support arm64/M1 platform, hence I couldn't use it in this comparison.

* Ruby = can sort unicode text but without collations becase it can't use ICU on arm64/M1

* Elixir = same as Python, no ICU for M1.

* Julia = I couldn't find a good i18 library supporting many languages.

* Crystal = currently supports only Turkish collations. Probably because the language is young and does not have a large enough community or company behind it.

* Golang = has rules for many languages. You can see the influence of a large company and community which makes Golang a mature solution. Sorting slowed the whole task down significantly, but the result is correct (in this case I only checked the results for the Polish language)

### Kudos

[@romanatnews](https://github.com/romanatnews) (Golang example refactoring)

[@pan93412](https://github.com/pan93412) (Rust example refactoring using Tokyo runtime)

## Changes

2024-03-02

* Updated Python version to 3.12.2, added poetry, solved missing icu4 collations for M1 processors, added a fancy progress bar
* Updated Golang version to 1.22.0
* Updated Ruby version to 3.3.0, added sorting with collations for many languages (this slowed the code almost 2x)

2022-09-17

* Updated all version. Slighty faster results. The order not changed.

2022-02-08

* Added improved Ruby code version with correct reading the pure text to tokenize (it ignores sigla in each verse), and with the correct regular expression for extracting words. The code is a little slower but it works almost as expected. (almost because for arm64/M1 it can't use ICU)

2022-02-20

* Added newer Golang 1.1.17 and improved code
* Added fixed Python version ignoring sigla (like in Ruby version)
* Added fixed Julia version ignoring sigla
