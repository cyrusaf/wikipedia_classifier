require 'stuff-classifier'

cls = nil
File.open("classifier", "r") do |f|
  puts "reading file..."
  content = f.read

  puts "parsing marshal..."
  cls = Marshal.load(content)
end

puts "classifying..."
puts cls.classify("The Warriors have to be frustrated right now. They breezed through the 2015-16 NBA regular season with such ease that the 3-1 Western Conference Finals deficit they are facing against the Thunder at the moment has to be a shock to their systems. And in particular, Draymond Green must be upset with the way he has played in Games 3 and 4 of the WCF. You could argue they have been two of the worst games of his entire NBA career. Following Golden State’s loss to Oklahoma City on Tuesday night, Green spoke with UNINTERRUPTED about how he was feeling. He was surprisingly cool, calm, and collected for a guy who is part of a team that is now facing a monumental challenge. But he talked at length about what he’s dealing with personally and even mentioned that Kobe Bryant texted him some words of encouragement after the loss. Kobe may be retired, but it sounds like he's still making an impact on the NBA Playoffs, just like he did for so many years.")
puts cls.classify("Visual Basic is a fantastic for programming.")
puts cls.classify("The Warriors scored 110 points going 2 for 10 from the field and 1 from 5 from the three point line.")
