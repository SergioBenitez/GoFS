#include <stdio.h>
#include <stdlib.h>
#include "benchmark.h"

#define UNUSED(x) (void)(x) 

const int NUM = 100;

int
ceil_div(int x, int y) {
  return (x + y - 1) / y;
}

char *
init_filename(int n, int end_len, char *end) {
  int len = ceil_div(n, 26);
  char *filename = (char *)malloc(len + end_len);

  for (int i = 0; i < end_len; ++i)
    filename[len + i] = end[i];

  // should used /dev/shm, but no tmpfs on Mac
  for (int i = 0; i < len; ++i)
    filename[i] = '@';
  
  return filename;
}

/**
 * A new FILE*[] is allocated and returned. Caller must free it.
 */
FILE **
open_many_c(Benchmark *b, int n, int (f) (FILE *, char *)) {
  // pausing for filename and fd array allocation
  bench_pause(b);

  // creating initial filename and fd array
  char *filename = init_filename(n, 5, ".out");
  FILE **files = (FILE **)malloc(n * sizeof(FILE*));
  
  // Done with allocations
  bench_resume(b);

  for (int i = 0; i < n; ++i) {
    filename[i / 26] += 1;
    files[i] = fopen(filename, "wb");
    if (f) f(files[i], filename);
  }

  // deallocating filename
  bench_pause(b);
  free(filename);
  bench_resume(b);

  return files;
}

void
unlink_all(Benchmark *b, int n) {
  // pausing for filename and fd array allocation
  bench_pause(b);

  // creating initial filename and fd array
  char *filename = init_filename(n, 5, ".out");
  FILE **files = (FILE **)malloc(n * sizeof(FILE*));
  
  // Done with allocations
  bench_resume(b);

  // unlinking
  for (int i = 0; i < n; ++i) {
    filename[i / 26] += 1;
    remove(filename);
  }
  
  // deallocating filename
  bench_pause(b);
  free(filename);
  bench_resume(b);
}

FILE **
open_many(Benchmark *b, int n) {
  return open_many_c(b, n, NULL);
}

void
close_all(FILE **files, int n) {
  for (int i = 0; i < n; ++i) fclose(files[i]);
}

int
help_close(FILE *f, char *name) {
  UNUSED(name);
  return fclose(f);
}

int
help_unlink(FILE *f, char *name) {
  UNUSED(f);
  return remove(name);
}

int
help_close_unlink(FILE *f, char *name) {
  help_close(f, name);
  return help_unlink(f, name);
}

void
OtC(Benchmark *b) {
  FILE **files = open_many(b, NUM);
  close_all(files, NUM);
  free(files);
}

void
OC(Benchmark *b) {
  FILE **files = open_many_c(b, NUM, help_close);
  free(files);
}

void
OtCtU(Benchmark *b) {
  FILE **files = open_many(b, NUM);
  close_all(files, NUM);
  unlink_all(b, NUM);
  free(files);
}

void
OCU(Benchmark *b) {
  FILE **files = open_many_c(b, NUM, help_close_unlink);
  free(files);
}

int main() {
  benchmark("Open-Close", OtC, 2);
  benchmark("OpenAndClose", OC, 2);
  benchmark("Open-Close-Unlink", OtCtU, 4);
  benchmark("OpenAndCloseAndUnlink", OCU, 4);
}
