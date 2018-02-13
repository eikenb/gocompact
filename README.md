Overview
--------

A simple formatter that works like goimports/gofmt except that it changes the
format to a more vertically compact layout. Similar to what you see in most
presentations/slides. You can use it to format code for such things or just if
you like being able to fit more code in your editor.

Used in conjunction with gofmt, you get to have the more compact layout while
editing then trigger gofmt upon saving. See the example vimrc for one possible
setup.

Installation
------------

    go get github.com/eikenb/gocompact


Changing
--------

The Makefile has targets to help keeping the printer/ code copied from the Go
root tree up to date. To update the printer/ code to the latest version you
run...

    make update-printer
    make toggle-patch
    (manually fix if patch won't apply)
    make refresh-patch
    (commit)

To add more changes to the format, the flow is...

    (make your changes)
    make refresh-patch
    (commit)
