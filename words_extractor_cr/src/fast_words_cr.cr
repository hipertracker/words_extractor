require "json"
require "yaml"

# TODO: Write documentation for `FastWordsCr`

module FastWordsCr
  VERSION = "0.2.0"
  CHARSET = "aąbcćdeęfghijklłmnńoópqrsśtuvwxyzźż"

  def self.main(outpath = "words")
    with_sorting = true
    concurrent = true

    prepare_folder(outpath, "*.txt")

    file_count = 0

    processed_files = Channel(Bool).new
    Dir.glob("../data/pl/**/*.yml").each do |path|
      if concurrent
        spawn do
          worker(path, outpath, with_sorting)
          processed_files.send true
        end
        file_count += 1
      else
        worker(path, outpath, with_sorting)
      end
    end
    if concurrent
      file_count.times do
        processed_files.receive
      end
    end
  end

  def self.worker(path, outpath, with_sorting)
    text = File.read(path.gsub(".yml", ".txt")).gsub("\n", " ").downcase

    words = text.split(/[^\p{L}]+/).to_set

    if with_sorting
      words = words.to_a.sort { |x, y| self.word_cmp(x, y) }
    end

    meta = File.open(path) { |file| YAML.parse(file) }
    filepath = %Q(#{outpath}/extracted-words-for-#{meta["label"]}.txt)
    File.write(filepath, words.join("\n"))
    puts "Saved #{filepath}"
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

elapsed_time = Time.measure do
  FastWordsCr.main
end
puts elapsed_time
