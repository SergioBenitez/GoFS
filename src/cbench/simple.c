#include <stdio.h>
#include "benchmark.h"

void
OtC() {
  // should used /dev/shm, but no tmpfs on Mac
  FILE *file = fopen("test.out", "wb");
  fseek(file, 100, SEEK_SET);
  fputs("hello, world!", file);
  fclose(file);
}

int main() {
  benchmark("OpenThenClose", OtC, 2);
}
