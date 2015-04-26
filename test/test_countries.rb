# encoding: utf-8

###
#  to run use
#     ruby -I ./lib -I ./test test/test_countries.rb


require 'helper'

class TestCountries < MiniTest::Test

  def test_country_list
    pak = Datapak::Pak.new( './pak/country-list/datapackage.json' )

    puts "name: #{pak.name}"
    puts "title: #{pak.title}"
    puts "license: #{pak.license}"

    pp pak.tables

    ## pak.table.each do |row|
    ##  pp row
    ## end

    puts pak.table.dump_schema

    # database setup 'n' config
    ActiveRecord::Base.establish_connection( adapter:  'sqlite3', database: ':memory:' )
    ActiveRecord::Base.logger = Logger.new( STDOUT )

    pak.table.up!
    pak.table.import!

    pp pak.table.ar_clazz

    assert true  # if we get here - test success
  end

end # class TestCountries

