# Example of Text File Parsing Across Several Programming Languages

The objective is to extract unique words from UTF-8 files and save the results into separate files.

A notable challenge in sorting words arises from the necessity to accommodate sorting rules specific to different languages' grammars. This is a particularly complex issue that is not present in the English language, where the character set does not exceed the basic ASCII standard.

## Results

The results presented here are based on 123 unique UTF-8 Bible text files in 23 languages, as utilized on the [mybible.pl](ybible.pl) site. These files occupy 504MB. (The repository contains only a few sample files in the 'data' folder. To test additional data, you can create multiple files by cloning the .txt (and the associated .yml) files under different names.)


* Platform: macOS 14.3.1
* Machine: MacBook Pro 16" 64GB 2TB M1Max 10 cores.

<pre>
1. Rust 1.63           = 1.15s, sorting: 1.64s
2. Golang 1.22.0       = 1.40s, sorting with collations: 1.71s
3. Crystal 1.5.1       = 5.61s
4. Python 3.12.2       = 5.69, sorting with collations: 6.04s
5. Elixir 1.14.0       = 7.34s
6. Crystal 1.11.2      = 13.27s
7. Julia 1.8.1         = 12.13s, sorting: 12.22s
8. Ruby 3.3.0          = 12.63s, sorting with collations: 22.00s
</pre>

### Conclusion

The newly optimized version of Golang code demonstrates impressive speed. It is slightly slower than Rust but outperforms other programming languages in terms of execution speed. ~~Golang is the only language at the moment with full mature i18n support for arm64/M1 platform.~~

* Rust = the current example uses [lexical-sort](https://lib.rs/crates/lexical-sort) which is not perfect. [There is no standard mature implementation of i18n in Rust](https://www.arewewebyet.org/topics/i18n/) at the moment.

* Python = has a great implementation of [ICU](https://icu.unicode.org/related) library ~~however it does not support arm64/M1 platform, hence I couldn't use it in this comparison.~~

* Ruby = can sort unicode text ~~but without collations becase it can't use ICU on arm64/M1~~ with collations but it slowes down the code almost 2x.

* Elixir = no ICU for M1 to sort with collations.

* Julia = I couldn't find a good i18 library supporting many languages.

* Crystal = currently supports only Turkish collations. Probably because the language is young and does not have a large enough community or company behind it (as March 2024 still no more collations)

* Golang = has rules for many languages. You can see the influence of a large company and community which makes Golang a mature solution similar to Python. ~~Sorting slowed the whole task down significantly, but the result is correct (in this case I only checked the results for the Polish language)~~

### Kudos

[@romanatnews](https://github.com/romanatnews) (Golang example refactoring)

[@pan93412](https://github.com/pan93412) (Rust example refactoring using Tokyo runtime)

## Changes

2024-03-02

* Updated Python version to 3.12.2, added poetry, solved missing icu4 collations for M1 processors, added a fancy progress bar
* Updated Golang version to 1.22.0
* Updated Ruby version to 3.3.0, added sorting with collations for many languages (this slowed the code almost 2x)
* Updated Crystal version to 1.11.2 (using 10 cores is slower than 8 and uzsing 8 cofres is slower than 4, all is slower than v1.5.1, strange)

2022-09-17

* Updated all version. Slighty faster results. The order not changed.

2022-02-08

* Added improved Ruby code version with correct reading the pure text to tokenize (it ignores sigla in each verse), and with the correct regular expression for extracting words. The code is a little slower but it works almost as expected. (almost because for arm64/M1 it can't use ICU)

2022-02-20

* Added newer Golang 1.1.17 and improved code
* Added fixed Python version ignoring sigla (like in Ruby version)
* Added fixed Julia version ignoring sigla
