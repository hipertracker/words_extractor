require "json"
require "yaml"

# TODO: Write documentation for `FastWordsCr`

module FastWordsCr
  VERSION = "0.1.0"

  def self.main(outpath = "words")
    prepare_folder(outpath, "*.txt")
    Dir.glob("../data/pl/**/*.yml").each do |path|
      # spawn do
        worker(path, outpath)
      # end
    end
    # Fiber.yield
  end

  def self.worker(path, outpath)
    text = File.read(path.gsub(".yml", ".txt")).gsub("\n", " ").downcase

    # 35sec
    # words_json = (text.split(/[^\p{L}]+/).to_set - Set{""}).to_json.downcase
    # words = Array(String).from_json(words_json).sort { |x, y| self.word_cmp(x, y) }

    # 35s
    # words = (text.split(/[^\p{L}]+/).to_set - Set{""}).to_a.sort do |x, y|
    #   self.word_cmp(x, y)
    # end

    # 7s (no sort)
    words = (text.split(/[^\p{L}]+/).to_set - Set{""}).to_a

    meta = File.open(path) { |file| YAML.parse(file) }
    filepath = %Q(#{outpath}/słowa - #{meta["label"]}.txt)
    puts filepath
    File.write(filepath, words.join("\n"))
  end

  def self.prepare_folder(folder : String, pattern : String)
    Dir.mkdir_p(folder) unless File.exists?(folder)
    Dir.glob("#{folder}/#{pattern}").each do |path|
      File.delete path
    end
  end

  def self.word_cmp(str1 : String, str2 : String, charset = "aąbcćdeęfghijklłmnńoópqrsśtuvwxyzźż")
    tokens1 = str1.downcase.split("")
    tokens2 = str2.downcase.split("")
    tokens1.each_with_index do |s1, i|
      return 1 unless tokens2[i]?
      idx1 = charset.index(s1) || -1
      idx2 = charset.index(tokens2[i]) || -1
      return 1 if idx1 > idx2
      return -1 if idx1 < idx2
    end
    0
  end
end

FastWordsCr.main
