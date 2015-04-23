# encoding: utf-8

###
#  to run use
#     ruby -I ./lib -I ./test test/test_pak.rb


require 'helper'

class TestPak < MiniTest::Test

  def test_s_and_p_500_companies

    pak = Datapak::Pak.new( './pak/s-and-p-500-companies/datapackage.json' )

    puts "name: #{pak.name}"
    puts "title: #{pak.title}"
    puts "license: #{pak.license}"

    pp pak.tables
    pp pak.table[0]['Symbol']
    pp pak.table[495]['Symbol']

    ## pak.table.each do |row|
    ##  pp row
    ## end

    puts pak.tables[0].dump_schema
    puts pak.tables[1].dump_schema

    # database setup 'n' config
    ActiveRecord::Base.establish_connection( adapter:  'sqlite3', database: ':memory:' )
    ActiveRecord::Base.logger = Logger.new( STDOUT )

    pak.table.up!
    pak.table.import!

    pak.tables[1].up!
    pak.tables[1].import!


    pp pak.table.ar_clazz


    company = pak.table.ar_clazz

    puts "Company.count: #{company.count}"
    pp company.first
    pp company.find_by!( symbol: 'ABT' )
    pp company.find_by!( name: '3M Co' )
    pp company.where( sector: 'Industrials' ).count
    pp company.where( sector: 'Industrials' ).all


    ### todo: try a join w/ belongs_to ??

    assert true  # if we get here - test success
  end

end # class TestPak

