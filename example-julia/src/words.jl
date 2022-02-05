module words_extractor_jl

using Distributed
using YAML
using Glob

const outdir = "words"

function worker(yaml_path, i, count)
    path = get_filepath(yaml_path)
    words = get_words(yaml_path)
    write(path, join(words, "\n"))
    println("[$(lpad(i, 3, ' '))/$count] $path")
end

function get_words(yaml_path)
    text_path = replace(yaml_path, ".yml" => ".txt")
    text = read(text_path, String) |> lowercase
    split(text, r"[\W\d]+") |> Set |> collect
end

function get_filepath(path)
    meta = YAML.load_file(path)
    """./$outdir/$(meta["lang"])-$(meta["code"]).txt"""
end

function rdir(dir::AbstractString, pat::Glob.FilenameMatch)
    result = String[]
    for (root, dirs, files) in walkdir(dir)
        filepaths = joinpath.(root, files)
        append!(result, filter!(f -> occursin(pat, f), filepaths))
    end
    return result
end

rdir(dir::AbstractString, pat::AbstractString) = rdir(dir, Glob.FilenameMatch(pat))

function main()
    if ispath(outdir)
        rm(outdir, recursive = true)
    end
    mkdir(outdir)
    paths = rdir("../data", fn"../data/??/*.yml")
    count = length(paths)
    i = 1
    Threads.@threads for path in paths
        # println("Spawn $path")
        worker(path, i, count)
        i += 1
    end
end

addprocs()
println(string("Workers ", nworkers()))
println(string("Processing... using ", Threads.nthreads(), " threads"))
@time main()
end # module
