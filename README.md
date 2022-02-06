# words_extractor

### Info

Example of a text file parsing in several programming languages. The goal is to extract unique words from utf-8 files and save results them into separate files.

### Results

The following results are for 123 unique utf-8 Bible text files in 23 languages (used at mybible.pl site) They take 504MB. (The repo contains only a few sample files in the 'data' folder. For testing more data you could multiple files by cloning *.txt (and the associated*.yml) file under different names)

* Platform: MacOS 12.2
* Machine: MacBook Pro 16" 64GB 2TB M1Max 10 cores.

<pre>
1. Rust 1.58      = 0.38s
2. Python 3.10.2  = 2.80s
3. Julia 1.7.1    = 4.522
4. Crystal 1.3.2  = 5.72s
5. Elixir 1.13.2  = 7.82s
6. Ruby 3.1.0     = 8.31s

Golang 1.17.6    = UNDER REFACTORING, stay tuned
</pre>

### Conclusion

The difficulty in sorting words is due to the need to handle sorting rules according to the language. This is quite a complex problem that does not exist for the English language where the character set does not exceed the basic ASCII standard.

* Rust = I couldn't find collations for sort rules in other languages.

* Julia = same as Rust

* Elixir = same as Rust

* Crystal = currently has Turkish-only collations. Probably because the language is young and does not have a large enough community or company behind it. The manual sorting was not perfect here, the algorithm needs to be improved.

* Python = has a great implementation of ICU library but unfortunately, it is not still available for the arm64 / M1 platform hence I couldn't use it in this comparison.

* Ruby = same as Python

* Golang = has rules for many languages. You can see the influence of a large company and community which makes Golang a mature solution. Sorting slowed the whole task down significantly, but the result is correct (in this case I only checked the results for the Polish language)