#include <stdio.h>
#include <sys/time.h>
#include <sys/resource.h>
#include "benchmark.h"

void
record_times(double *real, double *user, double *sys) {
  struct timeval t;
  struct rusage r;

  gettimeofday(&t, NULL);
  getrusage(RUSAGE_SELF, &r);

  *real = t.tv_sec + t.tv_usec * 1e-6;
  *user = r.ru_utime.tv_sec + r.ru_utime.tv_usec * 1e-6;
  *sys = r.ru_stime.tv_sec + r.ru_stime.tv_usec * 1e-6;
}

void
reset_timer(Benchmark *b) {
  b->real_start = b->real_end = 0;
  b->user_start = b->user_end = 0;
  b->sys_start = b->sys_end = 0;
  b->real = b->user = b->sys = b->reps = 0;
}

void
start_timer(Benchmark *b) {
  record_times(&b->real_start, &b->user_start, &b->sys_start);
}

void
stop_timer(Benchmark *b) {
  record_times(&b->real_end, &b->user_end, &b->sys_end);
}

void
agg_timers(Benchmark *b) {
  b->real += b->real_end - b->real_start;
  b->user += b->user_end - b->user_start;
  b->sys += b->sys_end - b->sys_start;
  b->reps++;
}

void
print_results(Benchmark *b) {
  printf("Reps: %d, Total time: %f\n\n", b->reps, b->real);

  printf("Totals:\n");
  printf("Real:\t%f\n", b->real);
  printf("User:\t%f\n", b->user);
  printf("Sys:\t%f\n\n", b->sys);

  printf("Averages:\n");
  printf("Real:\t%f\n", b->real / b->reps);
  printf("User:\t%f\n", b->user / b->reps);
  printf("Sys:\t%f\n", b->sys / b->reps);
}

void
benchmark(char *name, void f(), double min_time) {
  printf("Running '%s'...", name);

  Benchmark b;
  reset_timer(&b);
  for (double time = 0; time < min_time; time = b.real) {
    start_timer(&b);
    f();
    stop_timer(&b);
    agg_timers(&b);
  }

  printf("Done.\n");
  print_results(&b);
}
