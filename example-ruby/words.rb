require 'yaml'
require 'yaml'
require 'parallel'
require 'etc'
require 'fileutils'

outdir = 'words'

FileUtils.rm_rf(outdir)
Dir.mkdir(outdir)

t = Time.now

sorted = false

paths = Dir['../data/??/**/*.yml']
count = paths.count

sizes = Parallel.map_with_index(paths, in_processes: Etc.nprocessors) do |yaml_path, i|
  meta = YAML.load_file(yaml_path)
  filepath = yaml_path.gsub('.yml', '.txt')
  words = IO.read(filepath).downcase.strip.split(/[^\p{word}]+/).uniq
  words = words.sort if sorted
  outpath = "#{outdir}/#{meta['lang']}-#{meta['code']}.txt"
  puts(format('[%3d/%d] %s/%s', i, count, yaml_path, outpath))
  File.write(outpath, words.join("\n"))
  File.size(filepath)
end

puts "Total size: #{(sizes.sum / 1024.0 / 1024).round} MB"
puts "Total time: #{Time.now - t} s"
