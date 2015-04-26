# encoding: utf-8


require 'pp'
require 'forwardable'

### csv
require 'csv'
require 'json'


### downloader
require 'fetcher'

### activerecord w/ sqlite3
##  require 'active_support/all'    ## needed for String#binary? method
require 'active_record'



# our own code

require 'datapak/version'      ## let version always go first
require 'datapak/datapak'
require 'datapak/downloader'

module Datapak
  
  def self.import(*args)
    ## step 1: download
    dl = Downloader.new
    args.each do |arg|
      dl.fetch( arg )
    end

    ## step 2: up 'n' import
    args.each do |arg|
      pak = Pak.new( "./pak/#{arg}/datapackage.json" )
      pak.tables.each do |table|
        table.up!
        table.import!
      end
    end
  end

end # module Datapak



# say hello
puts Datapak.banner    if defined?($RUBYLIBS_DEBUG) && $RUBYLIBS_DEBUG

