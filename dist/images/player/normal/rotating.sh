#!/bin/bash

# Define the source image
SOURCE="jump0.png"

# Loop to create 8 rotated images
for i in {1..8}; do
  # Calculate the rotation angle
  ANGLE=$((i * 45))

  # Define the output filename
  OUTPUT="jump$i.png"

  # Rotate the image using ImageMagick 7's magick command
  magick "$SOURCE" -rotate $ANGLE "$OUTPUT"
done

echo "Rotation complete. Images saved as jump1.png to jump8.jpg."

