use Timex

defmodule WordsExtractor do
  @moduledoc nil

  def run do
    outdir = "words"
    File.rm_rf!(outdir)
    File.mkdir!(outdir)
    t = Duration.now()
    paths = Path.wildcard("../data/??/**/*.yml")
    count = length(paths)

    total_size =
      paths
      |> Enum.with_index(1)
      |> Task.async_stream(
        WordsExtractor,
        :worker,
        [outdir, count],
        ordered: false,
        timeout: :infinity
      )
      |> Enum.to_list()
      |> Enum.reduce(0, fn {:ok, num}, acc -> acc + num end)

    elapsed = Duration.diff(Duration.now(), t, :milliseconds)
    IO.puts("Total time: #{elapsed / 1000}s")
    IO.puts("Total size: #{(total_size / 1024 / 1024) |> round} MB")
  end

  def worker({path, i}, outdir, count) do
    %{"code" => code, "lang" => lang} = YamlElixir.read_from_file!(path)
    pat = Regex.compile!("[\W\d]+/u")

    filepath = String.replace(path, ".yml", ".txt")
    %File.Stat{:size => filesize} = File.stat!(filepath)

    content =
      File.read!(filepath)
      |> String.downcase()
      |> String.trim()
      |> then(&Regex.split(pat, &1))
      |> MapSet.new()
      |> Enum.join("\n")

    # sorting does not respect collation so it is ignored
    # |> Enum.sort()

    it = i |> Integer.to_string() |> String.pad_leading(3)
    IO.puts("[#{it}/#{count}] #{path}/#{lang}-#{code}.txt")

    File.write!("./#{outdir}/#{lang}-#{code}.txt", content)
    filesize
  end
end
