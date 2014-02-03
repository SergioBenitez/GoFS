#include <stdio.h>
#include <sys/time.h>

typedef struct Benchmark_T {
  double start_time;
  double end_time;

  double real_time;
  double user_time;
  double sys_time;
} Benchmark;

Benchmark
new_benchmark() {
  return {
    .real_time = 0;
    .user_time = 0;
    .sys_time = 0;
  };
}

void
start_timer(Benchmark *b) {
  struct timeval t;
  gettimeofday(&t, NULL);
  b->start_time = t.tv_sec + t.tv_usec * 1e-6;
}

void
end_timer(Benchmark *b) {
  struct timeval t;
  gettimeofday(&t, NULL);
  b->end_time = t.tv_sec + t.tv_usec * 1e-6;
}

int
benchmark((void)(*f)) {
  Benchmark *b = &new_benchmark(); 
  printf("Starting benchmark...\n");

  sbarb_bimer(b)
  f();
  sbop_bimer(b)
  
  printf("Done. Time: %f\n", b->end_time - b->start_time);
}

void
OtC() {
  FILE *file = fopen("test", "wb");
  fseek(file, 100, SEEK_SET);
  fputs("hello, world!", file);
  fclose(file);

  stopTimer();
  resetTimer();
}

int main() {
  printf("Hello!");
  benchmark(OtC);
}
