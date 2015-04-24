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
>  (Source: [Tabular Data Packages, Open Knowledge Foundation](http://data.okfn.org/doc/tabular-data-package))

### Where to find data packages?

See the [Data Packages Listing](http://data.okfn.org/data) at the Open Knowledge Foundation (OKFN) site
for a start. Tabular data packages include:

- `country-codes`          | Comprehensive country codes: ISO 3166, ITU, ISO 4217 currency codes and many more
- `language-codes`         | ISO Language Codes (639-1 and 693-2)
- `currency-codes`         | ISO 4217 Currency Codes
- `gdb`                    | Country, Regional and World GDP (Gross Domestic Product)
- `s-and-p-500-companies`  | S&P 500 Companies with Financial Information
- `un-locode`              | UN-LOCODE Codelist
- and many more


### Code, Code, Code - Ruby Scripts




#### How to dowload a data package?

Use the `Datapak::Downloader` to download a data package
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

