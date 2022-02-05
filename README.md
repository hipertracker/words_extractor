# words_extractor

### Info
Example of a text file parsing in several programming languages. The goal is to extract unique words from utf-8 files and save results them into separate files.

### Results

The following results are for 936 files (2745 MB) on MacOS 12.2 and MacBook Pro 16" 64GB 2TB M1Max 10 cores. (For more text files go into data/pl/* and duplicate files several times.) All examples are using a similar logic and approach.

<pre>
1. Rust v1.58.1      =  7.54s
2. Python v3.10.2    = 15.34s (with multiprocessing)
3. Julia v1.7.1      = 17.00s
4. Crystal v1.3.2    = 26.32s 
5. Ruby v3.1.0       = 40.94s (with Parallel)
6. Golang v1.18beta1 = 73.00s
7. Elixir v1.13.2    = 2m43s 
</pre>

### Conclusion

Rust is the fastest language beyond doubt.

What is surprised is pretty poor Golang's performance on this task. Crystal is faster than Golang but in this task, it is still slower than Python which is also surprising. (Neither Golang nor Crystal is my main field of expertise so maybe there is some room for improvement. Although I showed this code to people and nobody so far could improve it in any significant way. But if I find a better implementation I will update this comparison.)

The high Python performance is interesting. Although it is using a multiprocessing standard library for full CPU cores utilization this is still dynamic interpreted language after all, which is rather expected to be slower than statically typed languages. 