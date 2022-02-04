defmodule WordsExtractor do
  @moduledoc nil

  def run do
    outdir = "words"
    clean_dir(outdir)

    walk("../data/pl/", ".yml")
    |> Task.async_stream(
      WordsExtractor,
      :worker,
      [outdir],
      ordered: false,
      timeout: :infinity
    )
    |> Enum.to_list()
  end

  def clean_dir(path) do
    File.rm_rf!(path)
    File.mkdir!(path)
  end

  def worker(path, outdir) do
    %{"code" => code} = YamlElixir.read_from_file!(path)

    words =
      File.read!(String.replace(path, ".yml", ".txt"))
      |> String.downcase()
      |> String.trim()
      |> then(&Regex.split(~r/[\W\d]+/u, &1))
      |> MapSet.new()
      # sorting does not respect collation
      |> Enum.sort()

    File.write!("#{outdir}/extracted-#{code}.txt", Enum.join(words, "\n"))
    IO.puts(path)
  end

  def walk(path, pattern) do
    dir = String.to_charlist(path)
    regexp = String.to_charlist(pattern)

    :filelib.fold_files(dir, regexp, true, fn file, acc -> [file | acc] end, [])
    |> Enum.map(fn filepath -> to_string(filepath) end)
  end
end
