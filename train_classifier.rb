require 'stuff-classifier'

cls = StuffClassifier::Bayes.new("Wikipedia")

total_cats = Dir.entries('documents').length - 2
i = 1

Dir.foreach('documents') do |category|
  next if category == '.' or category == '..'

  puts "#{i}/#{total_cats}: #{category}"
  i += 1

  # do work on real items
  Dir.foreach("documents/#{category}") do |article|
      next if article == '.' or article == '..'

      File.open("documents/#{category}/#{article}", "r") do |f|
        content = f.read

        cls.train(category, content)
      end
  end
end

# Save classifier
File.open("classifier", "w") do |f|
    f.write Marshal.dump(cls)
end
