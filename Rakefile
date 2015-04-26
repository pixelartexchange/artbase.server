require 'hoe'
require './lib/datapak/version.rb'

Hoe.spec 'datapak' do

  self.version = Datapak::VERSION

  self.summary = 'datapak - yet another library to work with tabular data packages (*.csv files w/ datapackage.json)'
  self.description = summary

  self.urls    = ['https://github.com/textkit/datapak']

  self.author  = 'Gerald Bauer'
  self.email   = 'ruby-talk@ruby-lang.org'

  # switch extension to .markdown for gihub formatting
  self.readme_file  = 'README.md'
  self.history_file = 'HISTORY.md'

  self.extra_deps = [
    ['logutils', '>=0.6.1'],
    ['fetcher', '>=0.4.5'],
    ['activerecord'],
  ]

  self.licenses = ['Public Domain']

  self.spec_extras = {
    required_ruby_version: '>= 1.9.2'
  }

end
