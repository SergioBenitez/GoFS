export GOPATH=`pwd`
flag=$1
exp=$2

if [ `uname` != 'Linux' ]; then
  echo "Sorry! Please run this on Linux."
  exit 1
fi

if [ "$flag" == '-prof' ]; then
  if [ -z "$exp" ] || [ "$exp" == " " ]; then
    echo "Must supply benchmark regular expression."
    echo "Example: $0 $flag OWsC$"
    exit 1
  fi

  cpu=`mktemp -t cpuprofXXXXX`
  mem=`mktemp -t memprofXXXXX`
  echo "Profiling Go code..."
  go test -bench $exp -cpuprofile $cpu -memprofile $mem bench
  echo -e "CPU Profile: $cpu\nMEM Profile: $mem"
else
  gores=`mktemp -t goresXXXXX`
  echo "Running Go benchmarks...(output at $gores)"
  go test -bench . bench > $gores 2> /dev/null
  echo -e "Done. See $gores for output.\n"

  cres=`mktemp -t cresXXXXX`
  echo "Running C benchmarks...(output at $cres)"
  gcc -std=gnu99 src/cbench/*.c && ./a.out > $cres
  rm a.out
  echo -e "Done. See $cres for output.\n"

  cleanc=`mktemp -t cleancXXXXX`
  cat $cres | grep Real | awk '{ print $3 }' | cut -dn -f1 > $cleanc

  cleango=`mktemp -t cleangoXXXXX`
  cat $gores | head -n -1 | tail -n +2 | awk '{ print $3 }' > $cleango

  headers=`mktemp -t headersXXXXX`
  cat $gores | head -n -1 | tail -n +2 | awk '{ print $1 }' > $headers

  echo -e "Temp files: $cleanc, $cleango, $headers\n"
  echo "Results: (Go | C | ratio)"
  paste -d"\t" $headers $cleango $cleanc | awk '{ print $1 "\t" $2 "\t" $3 "\t" $2 / $3 }'
fi
