# Notes


## Old Heroku Doku / Setup



**Try Online**

Try the sample artbase web server / service installation running online
at [**pixelartexchange.herokuapp.com »**](https://pixelartexchange.herokuapp.com/).



**Install & Run Online Using Heroku**

Yes, you can. Run your own artbase (web) server / service
copy online using heroku.

Step 1 - Login to Heroku

```
$ heroku login
```

Step 2 - Create a Heroku app(lication)

```
$ heroku create  [app_name]
$ git push heroku master
```

That's it.
Test your online artbase (web) server / service
running at `https://[app_name].herokuapp.com`


---


Add `Procfile`:

```
web: artbase.server
```

