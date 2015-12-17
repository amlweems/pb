# shorty

A command line pastebin inspired by ix.io 😊

# Usage

Assume you have created a shortener on the domain `lf.lc`. Users can interact with your shortener as follows.

```
~$ echo Hello world. | curl -F 'f:1=<-' lf.lc
http://lf.lc/fpW

~$ curl lf.lc/fpW
Hello world.
```
