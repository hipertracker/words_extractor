require 'yaml'
require 'parallel'
require 'etc'
require 'fileutils'

outdir = 'words'
start = Time.now

FileUtils.rm_rf(outdir)
Dir.mkdir(outdir)

sorted = false

paths = Dir['../data/pl/**/*.yml']

Parallel.each(paths, in_processes: Etc.nprocessors) do |yaml_path| 
  meta = YAML.load_file(yaml_path)
  words = IO.read(yaml_path.gsub('.yml', '.txt')).downcase.strip.split(/[^\p{word}]+/).uniq
  if sorted
    words = words.sort
  end
  outpath = "#{outdir}/#{meta['lang']}-#{meta['code']}.txt"
  puts outpath
  File.write(outpath, words.join("\n"))
end

secs = Time.now - start
puts "Total time: #{secs} s"

