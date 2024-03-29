# artbase - zero-config web server / services; (auto-)downloads & serves pre-configured pixel art collections "out-of-the-box"; incl. 2x/4x/8x zoom for bigger image sizes and more; binaries for easy "xcopy" installation for windows, linux & friends




## Download Binaries For Easy "Xcopy" Installation

Find the archives to download  - about 3 Megabytes (MB) - for Windows, Linux and Friends at the [**Releases Page »**](https://github.com/pixelartexchange/artbase.server/releases).

Unpack the archive (e.g. `artbase-*.tar.gz` or `artbase-*.zip`) and than start / run the binary:

```
$ artbased
```

This will start-up a (web) server (listening on port 8080). To test open up `http://localhost:8080` in your browser (to get the index web page listing all collections).


## Build & Run From Source


Use / issue / type  (in the `/artbase.server` directory):

```
$ go build artbased.go
```

to get a zero-config x-copy binary for your operation system / architecture.
To run use:

```
$ artbased
```

This will start-up a (web) server (listening on port 8080). To test open up `http://localhost:8080` in your browser (to get the index web page listing all collections).



## Artbase - The Server Edition

The artbase (web) server will (auto-)download on demand the first-time only pre-configured
pixel art collections (using all-in-one image composites)
to your working directory and use the "cached" version from the next time on (incl. server restarts).



The pixel art collections pre-configured¹ include:

- [punks](https://github.com/cryptopunksnotdead/awesome-24px/blob/master/collection/punks.png) (24x24)
- [morepunks](https://github.com/cryptopunksnotdead/awesome-24px/blob/master/collection/morepunks.png) (24x24)
- [readymadepunks](https://github.com/cryptopunksnotdead/punks.readymade/blob/master/readymades.png) (24x24)
- [coolcats](https://github.com/cryptopunksnotdead/awesome-24px/blob/master/collection/coolcats.png)  (24x24)
- [mooncatresuce](https://github.com/cryptocopycats/awesome-mooncatrescue-bubble/blob/master/i/mooncatrescue.png)  (24x24)
- [boredapes](https://github.com/cryptopunksnotdead/awesome-24px/blob/master/collection/boredapes.png)  (28x28)
- And many more


¹: see [collections.csv](collections.csv)



Note:  You can use your own collections  - use the `-c/--config` flag
and pass along a file path or a web url to the collection dataset in
the comma-separated values (.csv) format e.g.  `$ artbased --config ./collections.csv`.





### (Web) Services


To get pixel art images, use `/:name/:id[.png|.svg]`.

Let's try the (default)
binary raster graphics format
using the portable network graphics (.png) encoding.
Example:

`/punks/0`, `/punks/1`, `/punks/2`,
(same as `/punks/0.png`, `/punks/1.png`, `/punks/2.png`) ...

![](i/punks-000000.png)
![](i/punks-000001.png)
![](i/punks-000002.png)


or `/coolcats/0`, `/coolcats/1`, `/coolcats/2`,
(same as `/coolcats/0.png`, `/coolcats/1.png`, `/coolcats/2.png`)  ...

![](i/coolcats-000000.png)
![](i/coolcats-000001.png)
![](i/coolcats-000002.png)



Let's try the scalable vector graphics (.svg) format in text.
Example:


`/punks/0.svg`, `/punks/1.svg`, `/punks/2.svg`, ...

![](i/punks-000000.svg)
![](i/punks-000001.svg)
![](i/punks-000002.svg)


or `/coolcats/0.svg`, `/coolcats/1.svg`, `/coolcats/2.svg`,  ...

![](i/coolcats-000000.svg)
![](i/coolcats-000001.svg)
![](i/coolcats-000002.svg)


Note: Pixels get "encoded" as rectangle "shapes" with a width and height
of one (1×1).  Transparent pixels
with the red/green/blue/alpha (rgba) value of (0 or 0/0/0/0)
get dropped.



#### z (zoom) Parameter - 2x, 4x, 8x, 10x, 20x ...   (.png only)


Note: The default image size is the default
(minimum) pixel size of the collection e.g. 24x24 for punks, morepunks,
coolcats and so on.
Use the z (zoom) parameter to upsize.

Let's try 2x:


`/punks/0?z=2`, `/punks/1?z=2`, `/punks/2?z=2`, ...

![](i/punks-000000@2x.png)
![](i/punks-000001@2x.png)
![](i/punks-000002@2x.png)


or `/coolcats/0?z=2`, `/coolcats/1?z=2`, `/coolcats/2?z=2`, ...

![](i/coolcats-000000@2x.png)
![](i/coolcats-000001@2x.png)
![](i/coolcats-000002@2x.png)




Let's try 8x:


`/punks/0?z=8`, `/punks/1?z=8`, `/punks/2?z=8`, ...

![](i/punks-000000@8x.png)
![](i/punks-000001@8x.png)
![](i/punks-000002@8x.png)


or `/coolcats/0?z=8`, `/coolcats/1?z=8`, `/coolcats/2?z=8`, ...  And so on.

![](i/coolcats-000000@8x.png)
![](i/coolcats-000001@8x.png)
![](i/coolcats-000002@8x.png)



#### bg (background) Parameter    (.png only)

Let's try adding the classic gray-ish/blue-ish v2 background
in red/green/blue (rgb) hexcode `#638596`.
Use the bg (background) parameter:

`/punks/0?bg=638596`, `/punks/1?bg=638596`, `/punks/2?bg=638596`, ...

![](i/punks-000000_(v2).png)
![](i/punks-000001_(v2).png)
![](i/punks-000002_(v2).png)


Let's try adding the baby blue-ish v3 background
in red/green/blue (rgb) hexcode `#60a4f7`:

`/punks/0?bg=60a4f7`, `/punks/1?bg=60a4f7`, `/punks/2?bg=60a4f7`, ...

![](i/punks-000000_(v3).png)
![](i/punks-000001_(v3).png)
![](i/punks-000002_(v3).png)




#### silhouette Parameter    (.png only)

Let's try a black silhouette.
Use the silhouette parameter:

`/punks/0?silhouette=black`, `/punks/1?silhouette=black`, `/punks/2?silhouette=black`, ...

![](i/punks-000000_silhouette(black).png)
![](i/punks-000001_silhouette(black).png)
![](i/punks-000002_silhouette(black).png)



#### Bonus -  Glory to Ukraine! Fuck (Vladimir) Putin! Stop the War! - Send A Stop The War Message To The World With Your Profile Picture


Let's try two-colored with the background in blue
and the silhouette (foreground) in yellow:

`/punks/0?bg=ukraineblue&silhouette=ukraineyellow`, `/punks/1?bg=ukraineblue&silhouette=ukraineyellow`, `/punks/2?bg=ukraineblue&silhouette=ukraineyellow`, ...

![](i/punks-000000_(ukraineblue)_silhouette(ukraineyellow).png)
![](i/punks-000001_(ukraineblue)_silhouette(ukraineyellow).png)
![](i/punks-000002_(ukraineblue)_silhouette(ukraineyellow).png)



Let's try the ukraine flag in the background:

`/punks/0?flag=ukraine`, `/punks/1?flag=ukraine`, `/punks/2?flag=ukraine`, ...

![](i/punks-000000_flag(ukraine).png)
![](i/punks-000001_flag(ukraine).png)
![](i/punks-000002_flag(ukraine).png)



####  Bonus  -  Philip! Phree the Phunks!

Let's try to flip vertically, that is, mirror, the images -
that turns right-looking images into left-looking and vice versa.
Use the m (mirror) parameter:


`/punks/0?m=t`, `/punks/1?m=t`, `/punks/2?m=t`, ...

![](i/punks-000000_mirror.png)
![](i/punks-000001_mirror.png)
![](i/punks-000002_mirror.png)


That's it for now.





## License

The `artbase` sources & binaries are dedicated to the public domain.
Use it as you please with no restrictions whatsoever.




## Questions? Comments?

Post them over at the [Help & Support](https://github.com/geraldb/help) page. Thanks.