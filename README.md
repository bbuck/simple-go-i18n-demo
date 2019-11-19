A simple demo of a Go application that runs a server that renders some templates.

- Template engine is [extemplate](https://github.com/dannyvankooten/extemplate)
  which is built on top of Go's html/template library to allow extending base
  templates.
- [Chi](https://github.com/go-chi/chi) was used as the webserver because it's
  idiomatic and has great utilities (IMO). I used this over the plain http
  library since it already has the concept of middleware, route params and
  contains a router.
- The i18n was rolled as an example implementation. The implementation is a bit
  cumbersome for a large amount of usage but it gets the point across.
