#!/bin/sh

############################################################
# set_default_font_color 为设置默认字体颜色
# set_default_font_color is to set default font color
#
# parameter 1: red(int)
# parameter 2: green(int)
# parameter 3: blue(int)
#
set_default_font_color (){
  echo -ne "\033]10;#676767\007"
}

# 先写obase，再去写ibase
hex_to_rgb () {
    hex=$(echo "#676767" | sed 's/#//g')

    a=$(echo "$hex" | cut -c-2)
    b=$(echo "$hex" | cut -c3-4)
    c=$(echo "$hex" | cut -c5-6)

    echo "$a" "$b" "$c"

    r=$(echo "obase=10; ibase=16; $a" | bc)
    g=$(echo "obase=10; ibase=16; $b" | bc)
    b=$(echo "obase=10; ibase=16; $c" | bc)

    echo "$r" "$g" "$b"

    return 0
}

rgb_to_hex () {
    r=$(echo "obase=16; ibase=10; $1" | bc)
    g=$(echo "obase=16; ibase=10; $2" | bc)
    b=$(echo "obase=16; ibase=10; $3" | bc)

    echo "$r" "$g" "$b"
    
    return 0
}