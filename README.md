# Sprite-Locator

Ever download a sprite sheet from an online source such as https://www.spriters-resource.com/, but find yourself painstakingly separating out each and individual sprite from pre-made sprite sheets? 

sprite-locator generates json metadata describing bounding boxes for all of the sprites it finds in a given png file. 

Build with `go build`

Run:

`./sprite-locator <sprite-sheet-file> <outfile>`

sprite-locator works by using a [flood-fill algorithm](https://en.wikipedia.org/wiki/Flood_fill).
sprite-locator picks the most commonly occurring color in a file as the "background color" and distinguishes sprite-pixels based on having a different color. 
sprite-locator finds contiguous blocks of non-background-color pixels and groups them as sprites.
after playing with sprite-locator, I found that a lof of small areas of pixels (usually related to shadows) get missed with this algorithm, so there is a configurable 'margin' that allows empty pixels to be included in the sprite bounding box algorithm.
the margin is set to 4 by default , but can be overridden by setting the `PIXEL_MARGIN` environment variable (to an integer value). make sure that it's set to at least 1, so that the algorithm can search adjacent pixels.

I haven't calculated the runtime of this algorithm (or the way I implemented it here) but it works at a reasonable speed. If anyone feels like taking a look at the code to help me optimize, submit a PR and I'd be glad to merge.