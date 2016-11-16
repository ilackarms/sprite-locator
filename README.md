# Sprite-Locator

Sprite locator is designed to locate sprites in a heterogeneous sprite sheet. Requires Go to build.

Build with `go build`

Run:

`./sprite-locator <sprite-sheet-file> <outfile>`

This will produce a json file listing the bounding rectangles of all the discovered sprites. Note that a whitespace margin of 4px is hard-coded. You can override the 4px margin by setting env var `PIXEL_MARGIN=1` before running.

Note: currently only supports png files, but can be extended easily

There is probably more efficient way of doing this, but it is processor and memory efficient enough to satisfy me. Enjoy.