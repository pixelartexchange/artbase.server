# encoding: utf-8

module CsvPack

class Downloader

  def initialize( cache_dir='./pack' )
    @cache_dir = cache_dir   # todo: check if folder exists now (or on demand)?
    @worker = Fetcher::Worker.new
  end

  SHORTCUTS = {
    ## to be done
  }

  def fetch( name_or_shortcut_or_url )   ## todo/check: use (re)name to get/update/etc. why? why not??

    name = name_or_shortcut_or_url

    ##
    ## e.g. try
    ##   country-list
    ##

    ## url_base = "http://data.okfn.org/data/core/#{name}"
    url_base = "https://datahub.io/core/#{name}"
    url = "#{url_base}/datapackage.json"

    dest_dir = "#{@cache_dir}/#{name}"
    FileUtils.mkdir_p( dest_dir )

    pack_path = "#{dest_dir}/datapackage.json"
    @worker.copy( url, pack_path )

    h = JSON.parse( File.read( pack_path ) )
    pp h

    ## copy resources (tables)
    h['resources'].each do |r|
      puts "== resource:"
      pp r

      res_url       = r['url']

      res_name          = r['name']
      res_relative_path = r['path']
      if res_relative_path.nil?
        res_relative_path = "#{res_name}.csv"
      end

      res_path = "#{dest_dir}/#{res_relative_path}"
      puts "[debug] res_path: >#{res_path}<"
      res_dir   = File.dirname( res_path )
      FileUtils.mkdir_p( res_dir )

      @worker.copy( res_url, res_path )
    end
  end

end # class Downloader

end # module CsvPack
