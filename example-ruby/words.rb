require 'yaml'
require 'parallel'
require 'etc'
require 'fileutils'
require 'optparse'

class WordsExtractor
  def initialize(cores: Etc.nprocessors, sorting: false, outdir: 'words', source: '../data/??/**/*.yml')
    @cores = cores
    @sorting = sorting
    @outdir = outdir
    @source = source
  end

  def clear_output
    FileUtils.rm_rf(@outdir)
    Dir.mkdir(@outdir)
  end

  def get_words(filepath)
    IO.readlines(filepath).map do |line|
      line.strip.downcase.split(' ')[2...-1].join(' ').split(/[^\p{L}]+/).uniq
    end.flatten.uniq
  end

  def save_words(words:, meta:, yaml_path:, count:, i:)
    outpath = "#{@outdir}/#{meta['lang']}-#{meta['code']}.txt"
    puts(format('[%3d/%d] %s/%s', i, count, yaml_path, outpath))
    File.write(outpath, words.join("\n"))
  end

  def run
    print "Running using #{@cores} processes"
    print ' with sorting' if @sorting
    puts '...'
    clear_output
    start = Time.now
    paths = Dir[@source]
    count = paths.count
    sizes = Parallel.map_with_index(paths, in_processes: @cores) do |yaml_path, i|
      meta = YAML.load_file(yaml_path)
      filepath = yaml_path.gsub('.yml', '.txt')
      words = get_words(filepath)
      words.sort! if @sorting
      save_words(words:, meta:, yaml_path:, count:, i:)
      File.size(filepath)
    end
    puts "Total size: #{(sizes.sum / 1024.0 / 1024).round} MB"
    puts "Total time: #{Time.now - start} s"
  end
end

if __FILE__ == $PROGRAM_NAME
  cores = Etc.nprocessors
  options = { s: false, n: cores }
  OptionParser.new do |opts|
    opts.banner = "Usage: ruby #{__FILE__} [options]"
    opts.on('-n [NUM]', OptionParser::DecimalInteger, "Number of cores to run (default #{cores})") do |val|
      options[:n] = if val.negative? || val > cores
                      cores
                    else
                      val
                    end
    end
    opts.on('-s', 'Sort results') { |v| options[:s] = v }
  end.parse!
  WordsExtractor.new(cores: options[:n], sorting: options[:s]).run
end
