module words_extractor_jl

import Pkg
Pkg.add("ArgParse")
Pkg.add("YAML")
Pkg.add("Glob")

using ArgParse
using Distributed
using YAML
using Glob


const outdir = "words"

function parse_commandline()
    s = ArgParseSettings()
    @add_arg_table s begin
        "-s"
        help = "Sort results"
        action = :store_true
    end
    return parse_args(s)
end

function worker(yaml_path, sorting, i, count)
    path = get_filepath(yaml_path)
    words = get_words(yaml_path, sorting)
    write(path, join(words, "\n"))
    println("[$(lpad(i, 3, ' '))/$count] $path")
end

function get_words(yaml_path, sorting = false)
    pattern = r"[\W\d]+"
    unique_words = Set()
    open(replace(yaml_path, ".yml" => ".txt")) do file
        for line in readlines(file)
            # exclude beginning book refrence from the line
            text = split(line, " ")[begin+2:end] |> t -> join(t, " ")
            tokens =
                text |>
                lowercase |>
                t -> split(t, pattern) |> t -> filter(token -> length(token) > 1, t)
            union!(unique_words, tokens)
        end
    end
    # unique_words = Set(words)
    if sorting
        arr = collect(unique_words)
        sort(arr)
    else
        unique_words
    end
end

function get_filepath(path)
    meta = YAML.load_file(path)
    """./$outdir/$(meta["lang"])-$(meta["code"]).txt"""
end

function rdir(dir::AbstractString, pat::Glob.FilenameMatch)
    result = String[]
    for (root, _dirs, files) in walkdir(dir)
        filepaths = joinpath.(root, files)
        append!(result, filter!(f -> occursin(pat, f), filepaths))
    end
    result
end

rdir(dir::AbstractString, pat::AbstractString) = rdir(dir, Glob.FilenameMatch(pat))

function main()
    parsed_args = parse_commandline()
    sorting = parsed_args["s"]

    addprocs()
    println(string("Workers ", nworkers()))
    println(string("Processing... using ", Threads.nthreads(), " threads"))
    if sorting
        println("with sorting")
    end
    if ispath(outdir)
        rm(outdir, recursive = true)
    end
    mkdir(outdir)
    paths = rdir("../data", fn"../data/??/*.yml")
    count = length(paths)
    i = 1
    Threads.@threads for path in paths
        # println("Spawn $path")
        worker(path, sorting, i, count)
        i += 1
    end
end

@time main()
end # module
