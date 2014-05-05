function lib_exists() {
  local name=$1

  stat /usr/local/lib/lib${name}.* 2> /dev/null 1>&2
  local loc=$?

  stat /usr/lib/lib${name}.* 2> /dev/null 1>&2
  local sys=$?
  
  return $loc || $sys
}

lib_exists "check"
if [ $? -eq 1 ]; then
  echo "Please install the Check unit test library to run the tests."
  exit 1
fi

gcc main.c ../*.c -lcheck -o test && ./test && rm test
