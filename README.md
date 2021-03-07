# TCOD

The `tcod` package is a fork of `libtcod-go` that's been updated to support the 
latest libtcod 1.16.6.  It's also been refactored to better support Go standards.

The libtcod library provides many utilities frequently used in roguelike games, like:

* truecolor console (SDL and OpenGL backends)
* keyboard/mouse input
* misc algorithms (line drawing, pathfinding, field of view, dungeon generation)
* terrain and noise generators
* widget toolkit
* config parser
* name generator

Most of libtcod API (version 1.5.1) is wrapped in Go, with some parts fully ported to enable easier
callbacks. In addition, the demo and terrain-generation tool were also fully ported to
serve as examples on how to use the library.

## Progress

This migration is a work in progress.  At this point, `tcod` will produce a few
deprecation warnings, due to the fact that `libtcod-go` is using a few older 
functions no longer supported in `libtcod`, but should build and work fine. I'm 
also in the process of rewriting the `sample` applicaition to bring it in line with `libtcod`'s C version,
and remove those deprecation warnings.

## Links

* [libtcod GitHub Repo](https://github.com/libtcod/libtcod)
* [libtcod Releases](https://github.com/libtcod/libtcod/releases)  
* [libtcod 1.16.6 docs](https://libtcod.readthedocs.io/en/latest)
* [libtcod 1.6.4 docs](https://libtcod.github.io/docs/index2.html?c=true&cpp=false&cs=false&py=false&lua=false)
* [libtcod-go](https://github.com/afolmert/libtcod-go) 

## Installation

To build the bindings, you will need the `libtcod` library from 
[libtcod Releases](https://github.com/libtcod/libtcod/releases).  See the `sample`
project for an idea of how to build and run your project.  Right now, I'm building
this on macOS, but it should work on Linux and Windows in a similar way (I'll test
those when I get further along in the refactoring)

## Documentation

**Unfortunately the `libtcod-go` package code wasn't documented.  I'll rectify that
as well.**

## Differences between Go `tcod` and the C/C++ version of `libtcod`

The entire `libtcod` library has not and will not be implemented in Go, as there are
Go alternatives that don't require the overhead of a C call.

Original API parts missing from Go bindings:

* custom containers (TCOD_list_t)
* thread/mutexes functions
* SDL callback renderer
* Networking functions

## TODO

* Update the package to remove the deprecated `libtcod` function calls
* Properly document the code
* Refactor console creation, so it's more in line with the C version
* Add any new features from more recent versions of `libtcod`
* Support for some `libtcod` and `SDL` libraries, e.g. I tend to use `SDL`'s event
  handling instead of `libtcod`'s.

## Credits

These are [Adam Folmert's](https://github.com/afolmert) original credits, so it 
seemed fitting to leave them in place:

* Go Authors for the Go language
* Jice, Mingos and others for the libtcod library
* Chris Hamons for API design ideas in the libtcod-net bindings
* Felipe Bichued for comments and ideas
* Alex Ogier for patches updating to Go 1.2

I'd also like to add a thanks to Adam for his initial efforts.  He doesn't appear to
have been on GitHub for a while...I hope that's simply from lack of interest and 
nothing untoward.  
