require 'yaml'
require 'parallel'
require 'etc'
require 'fileutils'
require 'optparse'
require 'uri'
require 'ffi-icu'

class WordsExtractor
  def initialize(outdir:, source:, cores: Etc.nprocessors, sorting: false)
    @cores = cores
    @sorting = sorting
    @outdir = outdir
    @source = source
  end

  def clear_output
    FileUtils.rm_rf(@outdir)
    Dir.mkdir(@outdir)
  end

  def get_collation(lang)
    mapper = {
      'ar' => 'ar_SA', # Arabic, Saudi Arabia
      'cs' => 'cs_CZ', # Czech, Czech Republic
      'da' => 'da_DK', # Danish, Denmark
      'de' => 'de_DE', # German, Germany
      'el' => 'el_GR', # Greek, Greece
      'en' => 'en_EN', # English
      'eo' => 'eo',    # Esperanto, not country-specific
      'es' => 'es_ES', # Spanish, Spain
      'fi' => 'fi_FI', # Finnish, Finland
      'fr' => 'fr_FR', # French, France
      'he' => 'he_IL', # Hebrew, Israel
      'hr' => 'hr_HR', # Croatian, Croatia
      'hu' => 'hu_HU', # Hungarian, Hungary
      'it' => 'it_IT', # Italian, Italy
      'lt' => 'lt_LT', # Lithuanian, Lithuania
      'la' => 'en_EN', # Latin locale is the same as English
      'nl' => 'nl_NL', # Dutch, Netherlands
      'pl' => 'pl_PL', # Polish, Poland
      'pt' => 'pt_PT', # Portuguese, Portugal
      'ru' => 'ru_RU', # Russian, Russia
      'sk' => 'sk_SK', # Slovak, Slovakia
      'sv' => 'sv_SE', # Swedish, Sweden
      'uk' => 'uk_UA'  # Ukrainian, Ukraine
    }
    mapper[lang]
  end

  def get_words(filepath)
    IO.readlines(filepath).map do |line|
      line.strip.downcase.split(' ')[2...-1].join(' ').split(/[^\p{L}]/).uniq.select { |s| s.size > 1 }
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
      if @sorting
        collator = ICU::Collation::Collator.new(get_collation(meta['lang']))
        words = words.sort { |a, b| collator.compare(a, b) }
      end
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
  WordsExtractor.new(
    cores: options[:n],
    sorting: options[:s],
    outdir: 'words',
    source: '../data/??/**/*.yml'
  ).run
end
