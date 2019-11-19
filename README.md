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

To see the demo in practice, clone this repo somewhere (it's go modules so you
probably still need the GO111MODULE var set and it does not necessarly need to
go in \$GOPATH/src).

Once downloaded just run `go build` then `./http-server` and you can visit
`/` which should render Index Page (in all "locales"), `/home` which should
render "Home Page" or "Página de inicio" in Spanish, and finally `/hello/{name}`
(where `{name}` is a name, like `/hello/Brandon`) which should render "Hello, {name}!"
or "¡Hola, {name}!" depending on locale.

To change locales just add a query string `?locale=es` and you should see the
changes on home/hello pages. `es` is the only locale hardcoded so you can simply
add a new one to the `localeMap` and then use it on subsequent builds for testing
purposes.
