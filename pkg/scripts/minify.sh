#!/bin/bash
bunx uglify-js tracker.js \
  --compress \
  sequences=true,dead_code=true,conditionals=true,booleans=true,unused=true,if_return=true,join_vars=true,drop_console=true \
  --mangle \
  toplevel=true \
  --output \
  comments=false,beautify=false \
  -o tracker.min.js
