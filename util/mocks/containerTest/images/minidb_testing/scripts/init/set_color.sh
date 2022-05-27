#!/bin/bash

############################################################
# set_default_font_color 为设置默认字体颜色
# set_default_font_color is to set default font color
#
# parameter 1: red(int)
# parameter 2: green(int)
# parameter 3: blue(int)
#
set_default_font_color (){
  # 先写obase，再去写ibase
  r_hex=$(echo "obase=16; ibase=10; $1" | bc)
  g_hex=$(echo "obase=16; ibase=10; $2" | bc)
  b_hex=$(echo "obase=16; ibase=10; $3" | bc)

  # 设置默认字体颜色 set default font color
  echo -n "\033]10;#$r_hex$g_hex$b_hex\007"

  return 0
}

############################################################
# hex_to_rgb_test 为测试 hex_to_rgb
#
hex_to_rgb_test () {
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