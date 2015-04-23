# encoding: utf-8


require 'pp'
require 'forwardable'

### csv
require 'csv'
require 'json'

### activerecord w/ sqlite3
##  require 'active_support/all'    ## needed for String#binary? method
require 'active_record'



# our own code

require 'datapak/version'      ## let version always go first
require 'datapak/datapak'




# say hello
puts Datapak.banner    if defined?($RUBYLIBS_DEBUG) && $RUBYLIBS_DEBUG

