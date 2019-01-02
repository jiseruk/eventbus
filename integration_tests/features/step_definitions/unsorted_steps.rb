Cuando("pasan {int} segundos") do |int|
  puts "Sleeping #{int} seconds" if $debug
  int.times do
    sleep 1
    print "." if $debug
  end
  puts
end
