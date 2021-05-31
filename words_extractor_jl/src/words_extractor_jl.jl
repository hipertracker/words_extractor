module words_extractor_jl

using Distributed
using Pipe
using YAML

folder = "words"

function worker(yaml_path) 
    path = get_filepath(yaml_path)      
    words = get_words(yaml_path)
    write(path, join(words, "\n")) 
    println(string("Saved...", path))
end

function get_words(yaml_path)
    text_path = replace(yaml_path, ".yml" => ".txt")
    text = read(text_path, String) |> lowercase
    split(text, r"[\W\d]+") |> Set |> collect
end

function get_filepath(path)
    meta = YAML.load_file(path)
    string(folder, "/extracted-words-for-", meta["label"], ".txt")
end

function walk(path, file_ext)
    res = []
    for (root, _, files) in walkdir(path, topdown = true)
        for file in files
            if endswith(file, file_ext)
                filepath = joinpath(root, file)                    
                push!(res, filepath)
            end
        end
    end
    res
end

function main()
    if ispath(folder)
        rm(folder, recursive = true)
    end
    mkdir(folder)
    Threads.@threads for path in walk("../data/pl/", ".yml")
        println("Spawn $path")
        worker(path)
    end
end

# addprocs()
# println(string("Workers ", nworkers()))
println(string("Processing... using ", Threads.nthreads(), " threads"))
@time main()
end # module
