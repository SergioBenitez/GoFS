#include <stdio.h>
#include <stdlib.h>
#include "benchmark.h"

const int NUM = 100;

int
ceil_div(int x, int y) {
  return (x + y - 1) / y;
}

/**
 * A new FILE*[] is allocated and returned. Caller must free it.
 */
FILE **
open_many(int n) {
  int len = ceil_div(n, 26);
  
  char *end = ".out";
  char *filename = (char *)malloc(len + 5); // + 5 (4 from .out, 1 for \0)
  for (int i = 0; i < 5; ++i) filename[len + i] = end[i];

  // should used /dev/shm, but no tmpfs on Mac
  for (int i = 0; i < len; ++i) filename[i] = '@';

  FILE **files = (FILE **)malloc(n * sizeof(FILE*));
  for (int i = 0; i < n; ++i) {
    filename[i / 26] += 1;
    files[i] = fopen(filename, "wb");
  }

  return files;
}

void
OtC() {
  FILE **files = open_many(NUM);
  for (int i = 0; i < NUM; ++i) {
    fclose(files[i]);
  }
  free(files);
}

int main() {
  benchmark("OpenThenClose", OtC, 3);
}
