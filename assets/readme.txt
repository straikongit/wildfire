imagemagick:

convert tree8x8.png -resize 8x8
convert map1024x768.png -crop 64x48+497+461 map64.png
convert +append fire8x8small*.png fire.png

