GOPATH=/code/

gores=`tempfile`
echo "Running Go benchmarks...(output at $gores)"
go test -bench . bench > $gores 2> /dev/null
echo -e "Done. See $gores for output.\n"

cres=`tempfile`
echo "Running C benchmarks...(output at $cres)"
gcc -std=gnu99 src/cbench/*.c && ./a.out > $cres
rm a.out
echo -e "Done. See $cres for output.\n"

cleanc=`tempfile`
cat $cres | grep Real | awk '{ print $3 }' | cut -dn -f1 > $cleanc

cleango=`tempfile`
cat $gores | head -n -1 | tail -n +2 | awk '{ print $3 }' > $cleango

headers=`tempfile`
echo -e "Temp files: $cleanc, $cleango, $headers\n"
echo "Results: (Go/C ratio)"
cat $gores | head -n -1 | tail -n +2 | awk '{ print $1 }' > $headers
paste -d"\t" $headers $cleango $cleanc | awk '{ print $1 "\t" $2 / $3 }'