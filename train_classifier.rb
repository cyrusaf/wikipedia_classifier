require 'stuff-classifier'

cls = StuffClassifier::Bayes.new("Wikipedia")

total_cats = Dir.entries('data').length
i = 0

Dir.foreach('data') do |category|
  next if category == '.' or category == '..'

  puts "#{i}/#{total_cats}"
  i += 1

  # do work on real items
  Dir.foreach("data/#{category}") do |article|
      next if article == '.' or article == '..'

      File.open("data/#{category}/#{article}", "r") do |f|
        content = f.read

        cls.train(category, content)
      end
  end
end

# Save classifier
File.open("classifier", "w") do |f|
    f.write Marshal.dump(cls)
end
