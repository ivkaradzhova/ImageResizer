
def raise_and_rescue
  a = true
  begin

    puts "begin #{a} Arise!"

    # using raise to create an exception
    raise 'Exception Created!'


    # using Rescue method
  rescue


    a = false
    puts "Finally #{a} Saved!"


  ensure

    puts "ensure #{a} Saved!"
  end

  puts "Outside from Begin Block! #{a}"

end

# calling method
raise_and_rescue