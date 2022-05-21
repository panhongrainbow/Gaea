#!/bin/sh

# sedReplace 为取代字符串
# sedReplace is to replace string
sedReplace () {
    sed -i "s/$2/$3/g" "$1"
}