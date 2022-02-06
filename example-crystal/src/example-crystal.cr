require "json"
require "yaml"

module Example::Crystal
  VERSION = "0.3.0"
  CHARSET = "aąbcćdeęfghijklłmnńoópqrsśtuvwxyzźż"

  def self.main(outdir = "words")
    with_sorting = false
    concurrent = true

    prepare_folder(outdir, "*.txt")

    file_count = 0
    total_size = 0

    channel = Channel(Tuple(String, Int64)).new
    srcPath = "../data/??/**/*.yml"
    paths = Dir.glob(srcPath, follow_symlinks: true)
    count = paths.size
    paths.each do |path|
      if concurrent
        spawn do
          channel.send worker(path, outdir, with_sorting)
        end
        file_count += 1
      else
        worker(path, outdir, with_sorting)
      end
    end
    if concurrent
      file_count.times do |i|
        path, size = channel.receive
        total_size += size
        puts(::sprintf("[%3d/%d] %s", i + 1, file_count, path))
      end
    end
    puts("Total size: #{(total_size / 1024 / 1024).round} MB")
  end

  def self.worker(path, outdir, with_sorting)
    filepath = path.gsub(".yml", ".txt")
    filesize = File.size(filepath)
    text = File.read(filepath).gsub("\n", " ").downcase

    words = text.split(/[^\p{L}]+/).to_set

    if with_sorting
      words = words.to_a.sort { |x, y| self.word_cmp(x, y) }
    end

    meta = File.open(path) { |file| YAML.parse(file) }
    outfilepath = %Q(#{outdir}/#{meta["lang"]}-#{meta["code"]}.txt)
    File.write(outfilepath, words.join("\n"))
    filesize = File.size(filepath)
    {filepath, filesize}
  end

  def self.prepare_folder(folder : String, pattern : String)
    Dir.mkdir_p(folder) unless File.exists?(folder)
    Dir.glob("#{folder}/#{pattern}").each do |path|
      File.delete path
    end
  end

  def self.word_cmp(str1 : String, str2 : String)
    tokens2 = str2.chars
    str1.chars.each_with_index do |s1, i|
      return 1 unless tokens2[i]?
      idx1 = CHARSET.index(s1) || -1
      idx2 = CHARSET.index(tokens2[i]) || -1
      return 1 if idx1 > idx2
      return -1 if idx1 < idx2
    end
    0
  end
end

elapsed = Time.measure do
  Example::Crystal.main
end

puts("Total time: #{elapsed}")
