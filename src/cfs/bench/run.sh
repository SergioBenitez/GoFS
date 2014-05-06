gcc bench.c ../*.c ../../cbench/benchmark.c -O3 -std=gnu99 -o bench.out && ./bench.out
rm bench.out
