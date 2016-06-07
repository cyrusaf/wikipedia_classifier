require 'stuff-classifier'
require 'open-uri'
require 'nokogiri'

$cls = nil

File.open("classifier", "r") do |f|
  puts "reading file..."
  content = f.read

  puts "parsing marshal..."
  $cls = Marshal.load(content)
end

def classifyLink(link)
  content = ''

  doc = Nokogiri::HTML(open(link))

  doc.css("#bodyContent").each do |c|
    content += c.content
  end

  return $cls.classify(content)
end


puts "Classifying #{ARGV[0]}..."
puts classifyLink(ARGV[0])
