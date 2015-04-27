# datapak

yet another library to work with tabular data packages (*.csv files w/ datapackage.json)

* home  :: [github.com/textkit/datapak](https://github.com/textkit/datapak)
* bugs  :: [github.com/textkit/datapak/issues](https://github.com/textkit/datapak/issues)
* gem   :: [rubygems.org/gems/datapak](https://rubygems.org/gems/datapak)
* rdoc  :: [rubydoc.info/gems/datapak](http://rubydoc.info/gems/datapak)
* forum :: [ruby-talk@ruby-lang.org](http://www.ruby-lang.org/en/community/mailing-lists/)


## Usage


### What's a tabular data package?

> Tabular Data Package is a simple structure for publishing and sharing
> tabular data with the following key features:
>
> - Data is stored in CSV (comma separated values) files
> - Metadata about the dataset both general (e.g. title, author)
>   and the specific data files (e.g. schema) is stored in a single JSON file
>   named `datapackage.json` which follows the Data Package format
>
> <!-- break -->
>
> (Source: [Tabular Data Packages, Open Knowledge Foundation](http://data.okfn.org/doc/tabular-data-package))

Here's a minimal example of a tabular data package holding two files, that is, `data.csv` and `datapackage.json`:
 
`data.csv`:

~~~~
Brewery,City,Name,Abv
Andechser Klosterbrauerei,Andechs,Doppelbock Dunkel,7%
Augustiner Bräu München,München,Edelstoff,5.6%
Bayerische Staatsbrauerei Weihenstephan,Freising,Hefe Weissbier,5.4%
Brauerei Spezial,Bamberg,Rauchbier Märzen,5.1%
Hacker-Pschorr Bräu,München,Münchner Dunkel,5.0%
Staatliches Hofbräuhaus München,München,Hofbräu Oktoberfestbier,6.3%
...
~~~~

`datapackage.json`:

~~~~
{
  "name": "beer",
  "resources": [
    {
      "path": "data.csv",
      "schema": {
        "fields": [ { "name": "Brewery",   "type": "string" },
                    { "name": "City",      "type": "string" },
                    { "name": "Name",      "type": "string" },
                    { "name": "Abv",       "type": "number" } ]
      }
    }
  ]
}
~~~~

### Where to find data packages?

For some more real world examples see the [Data Packages Listing](http://data.okfn.org/data) at the Open Knowledge Foundation (OKFN) site for a start. Tabular data packages include:

Name                     | Comments
------------------------ | -------------
`country-codes`          | Comprehensive country codes: ISO 3166, ITU, ISO 4217 currency codes and many more
`language-codes`         | ISO Language Codes (639-1 and 693-2)
`currency-codes`         | ISO 4217 Currency Codes
`gdb`                    | Country, Regional and World GDP (Gross Domestic Product)
`s-and-p-500-companies`  | S&P 500 Companies with Financial Information
`un-locode`              | UN-LOCODE Codelist

and many more


### Code, Code, Code - Ruby Scripts

~~~
require 'datapak`

Datapak.import(
  's-and-p-500-companies',
  'gdb'
)
~~~

Using `Datapak.import` will:

1) download all data packages to the `./pak` folder

2) (auto-)add all tables to an in-memory SQLite database using SQL `create_table`
   commands via `ActiveRecord` migrations e.g.

~~~
create_table :constituents_financials do |t|
  t.string :symbol          # Symbol         (string)
  t.string :name            # Name           (string)
  t.string :sector          # Sector         (string)
  t.float  :price           # Price          (number)
  t.float  :dividend_yield  # Dividend Yield (number)
  t.float  :price_earnings  # Price/Earnings (number)
  t.float  :earnings_share  # Earnings/Share (number)
  t.float  :book_value      # Book Value     (number)
  t.float  :_52_week_low    # 52 week low    (number)
  t.float  :_52_week_high   # 52 week high   (number)
  t.float  :market_cap      # Market Cap     (number)
  t.float  :ebitda          # EBITDA         (number)
  t.float  :price_sales     # Price/Sales    (number)
  t.float  :price_book      # Price/Book     (number)
  t.string :sec_filings     # SEC Filings    (string)
end
~~~

3) (auto-)import all datasets using SQL inserts e.g.

~~~
INSERT INTO constituents_financials
  (symbol,
   name,
   sector,
   price,
   dividend_yield,
   price_earnings,
   earnings_share,
   book_value,
   _52_week_low,
   _52_week_high,
   market_cap,
   ebitda,
   price_sales,
   price_book,
   sec_filings)
VALUES
  ('MMM',
   '3M Co',
   'Industrials',
   162.27,
   2.11,
   22.28,
   7.284,
   25.238,
   123.61,
   162.92,
   104.0,
   8.467,
   3.28,
   6.43,
   'http://www.sec.gov/cgi-bin/browse-edgar?action=getcompany&CIK=MMM')
~~~

4) (auto-)add ActiveRecord models for all tables.


So what? Now you can use all the "magic" of ActiveRecord to query
the datasets. Example:

~~~
puts "Constituent.count: #{Constituent.count}"

# SELECT COUNT(*) FROM "constituents"
# => 496


pp Constituent.first

# SELECT  "constituents".* FROM "constituents" ORDER BY "constituents"."id" ASC LIMIT 1
# => #<Constituent:0x9f8cb78
         id:     1,
         symbol: "MMM",
         name:   "3M Co",
         sector: "Industrials">


pp Constituent.find_by!( symbol: 'MMM' )

# SELECT  "constituents".*
         FROM "constituents"
         WHERE "constituents"."symbol" = "MMM"
         LIMIT 1
# => #<Constituent:0x9f8cb78
         id:     1,
         symbol: "MMM",
         name:   "3M Co",
         sector: "Industrials">


pp Constituent.find_by!( name: '3M Co' )

# SELECT  "constituents".*
          FROM "constituents"
          WHERE "constituents"."name" = "3M Co"
          LIMIT 1
# => #<Constituent:0x9f8cb78
         id:     1,
         symbol: "MMM",
         name:   "3M Co",
         sector: "Industrials">


pp Constituent.where( sector: 'Industrials' ).count

# SELECT COUNT(*) FROM "constituents"
         WHERE "constituents"."sector" = "Industrials"
# => 63


pp Constituent.where( sector: 'Industrials' ).all

# SELECT "constituents".*
         FROM "constituents"
         WHERE "constituents"."sector" = "Industrials"
# => [#<Constituent:0x9f8cb78
          id:     1,
          symbol: "MMM",
          name:   "3M Co",
          sector: "Industrials">,
      #<Constituent:0xa2a4180
          id:     8,
          symbol: "ADT",
          name:   "ADT Corp (The)",
          sector: "Industrials">,...]

and so on
~~~


#### How to dowload a data package ("by hand")?

Use the `Datapak::Downloader` class to download a data package
to your disk (by default data packages get stored in `./pak`).

~~~
dl = Datapak::Downloader.new
dl.fetch( 'language-codes' )
dl.fetch( 's-and-p-500-companies' )
dl.fetch( 'un-locode`)
~~~

Will result in:

~~~
-- pak
   |-- language-codes
   |   |-- data
   |   |   |-- language-codes-3b2.csv
   |   |   |-- language-codes.csv
   |   |   `-- language-codes-full.csv
   |   `-- datapackage.json
   |-- s-and-p-500-companies
   |   |-- data
   |   |   |-- constituents.csv
   |   |   `-- constituents-financials.csv
   |   `-- datapackage.json
   `-- un-locode
       |-- data
       |   |-- code-list.csv
       |   |-- country-codes.csv
       |   |-- function-classifiers.csv
       |   |-- status-indicators.csv
       |   `-- subdivision-codes.csv
       `-- datapackage.json
~~~

#### How to add and import a data package ("by hand")?

Use the `Datapak::Pak` class to read-in a data package
and add and import into an SQL database. 

~~~
pak = Datapak::Pak.new( './pak/un-locode/datapackage.json' )
pak.tables.each do |table|
  table.up!      # (auto-) add table  using SQL create_table via ActiveRecord migration
  table.import!  # import all records using SQL inserts
end
~~~


#### How to connect to a different SQL database?

You can connect to any database supported by ActiveRecord. If you do NOT
establish a connection in your script - the standard (default fallback)
is using an in-memory SQLite3 database.

##### SQLite

For example, to create an SQLite3 database on disk - lets say `datapak.db` -
use in your script (before the `Datapak.import` statement):

~~~
ActiveRecord::Base.establish_connection( adapter:  'sqlite3
                                         database: './datapak.db' )
~~~

##### PostgreSQL 

For example, to connect to a PostgreSQL database use in your script
(before the `Datapak.import` statement):

~~~
require 'pg'       ##  pull-in PostgreSQL (pg) machinery

ActiveRecord::Base.establish_connection( adapter:  'postgresql'
                                         username: 'ruby'",
                                         password: 'topsecret',
                                         database: 'database' )
~~~


## Install

Just install the gem:

    $ gem install datapak


## Alternatives

See the "[Tools and Plugins for working with Data Packages](http://data.okfn.org/tools)"
page at the Open Knowledge Foundation (OKFN).


## License

The `datapak` scripts are dedicated to the public domain.
Use it as you please with no restrictions whatsoever.

## Questions? Comments?

Send them along to the ruby-talk mailing list. Thanks!

