Dir.entries(".").each do |file|
  if file.include? 'PC Computer - Diablo Diablo Hellfire - '
    f = file.sub 'PC Computer - Diablo Diablo Hellfire - ', ''
    f.downcase!
    f.gsub! ' ', '_'
    f.gsub! '_in', ''
    f.gsub! '_with', ''
    f.gsub! '_armor', ''
    f.gsub! '_&', ''
    puts "mv '#{file}' #{f}"
  end
end
