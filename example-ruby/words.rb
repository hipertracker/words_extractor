require 'yaml'
require 'parallel'
require 'etc'
require 'fileutils'

outdir = 'words'
start = Time.now

FileUtils.rm_rf(outdir)
Dir.mkdir(outdir)

sorted = false

paths = Dir['../data/??/**/*.yml']
total_size = 0
count = paths.count

Parallel.each_with_index(paths, in_processes: Etc.nprocessors) do |yaml_path, i|
  meta = YAML.load_file(yaml_path)
  words = IO.read(yaml_path.gsub('.yml', '.txt')).downcase.strip.split(/[^\p{word}]+/).uniq
  words = words.sort if sorted
  total_size += words.count
  outpath = "#{outdir}/#{meta['lang']}-#{meta['code']}.txt"
  # puts "[#{i}/#{count}] #{yaml_path}/#{outpath}"
  puts(format('[%3d/%d] %s/%s', i, count, yaml_path, outpath))
  File.write(outpath, words.join("\n"))
end

secs = Time.now - start
puts "Total size: #{total_size}"
puts "Total time: #{secs} s"
