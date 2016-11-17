# Sprite-Locator

Ever download a sprite sheet from an online source such as https://www.spriters-resource.com/, but find yourself painstakingly separating out each and individual sprite from pre-made sprite sheets? 

Now you can use sprite-locator to automatically find bounding boxes for you inside your images. This can be used to directly import unevenly-spaced sprite sheets into your game, or automate image processing tools to separate out sprites from sprite sheets so you don't have to do it by hand.

Let's say we have the following sprite sheet: 

<img src="http://i.imgur.com/JseTzfw.png">

analyzing with sprite-locator produces the following:

```
{
  "sprites": [
  ...
    {
      "min": {
        "x": 222,
        "y": 114
      },
      "max": {
        "x": 239,
        "y": 141
      }
    },
    {
      "min": {
        "x": 223,
        "y": 42
      },
      "max": {
        "x": 240,
        "y": 68
      }
    },
    ...
  ]
}

```

which describes the bounding box for each sprite found in the sheet. here's a (rough) approximation of what these boxes represent:

<img src="http://i.imgur.com/ukQOsR3.png">

Build with `go build`

Run:

`./sprite-locator <sprite-sheet-file> <outfile>`

- sprite-locator works by using a [flood-fill algorithm](https://en.wikipedia.org/wiki/Flood_fill).
- sprite-locator picks the most commonly occurring color in a file as the "background color" and distinguishes sprite-pixels based on having a different color. 
- sprite-locator finds contiguous blocks of non-background-color pixels and groups them as sprites.
- after playing with sprite-locator, I found that a lof of small areas of pixels (usually related to shadows) get missed with this algorithm, so there is a configurable 'margin' that allows empty pixels to be included in the sprite bounding box algorithm.
- the margin is set to 4 by default , but can be overridden by setting the `PIXEL_MARGIN` environment variable (to an integer value). make sure that it's set to at least 1, so that the algorithm can search adjacent pixels.

I haven't calculated the runtime of this algorithm (or the way I implemented it here) but it works at a reasonable speed. If anyone feels like taking a look at the code to help me optimize, submit a PR and I'd be glad to merge.