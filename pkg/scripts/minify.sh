#!/bin/bash
bunx uglify-js --compress --mangle --toplevel --output comments=false,beautify=false -o tracker.min.js tracker.js
