Video automation, generating videos out of animated HTML5 pages
===============================================================

This little project does only one simple thing: it generates a series
of videos based on a `data.yaml` (list of data objects) and a
`template.html`, rendering each as MP4 files (using ffmpeg) from the
moment an element `.capture` is displayed, until it is hidden.

My primary goal was to create dynamic videos to introduce speakers in
the video recording of the sessions at Golang Montréal meetups
(https://golangmontreal.org).  I wanted to cut the post-production
cost of doing videos recordings of the session, so I coded the
repetitive stuff.

This program reads a `data.yaml` similar to this:

```
- slug: 01-alexandre-bourget
  first: Alexandre
  last: Bourget
  twitter: bourgetalexndre
  tagline: Data Scientist, Intel Security

  meetup: "#gomtl-01"
  date: 22 février 2016

- slug: 02-robert-de-niro
  first: Robert
  last: De Niro
  tagline: POS Dev, Redshift Inc

  meetup: "#gomtl-01"
  date: February 22nd 2016
```

See the included `template.html` for an example template.

It then uses ChromeDriver (using
[the agouti WebDriver implementation](https://github.com/sclevine/agouti)
to open the browser, resize it properly, load the template from its
own internal web server.  It waits until an HTML node with a `capture`
class appears (on the body for example). It then starts `ffmpeg` to
capture the screen at 25 fps, writing an .mp4 file to disk for each
`slug` in the `data.yaml`.  When the `capture` class disappears, the
video stops.  It then proceeds to the next video.

Pretty simple eh ?

Dependencies
------------

You need to have the binary for `ffmpeg` available in your PATH.

Limitations
-----------

This currently runs only on Linux (uses the x11grab ffmpeg
module). The window-decorations dimensions are modeled around the
Ubuntu Unity environment and are currently hard-coded.

It could be adapted to run on Windows though.
