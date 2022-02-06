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

Golang 1.17    = UNDER REFACTORING, stay tuned
</pre>

### Conclusion

Rust is the fastest language beyond doubt.

The high Python performance is interesting. Although it is using a multiprocessing standard library for full CPU cores utilization this is still dynamic interpreted language after all, which is rather expected to be slower than statically typed languages.
