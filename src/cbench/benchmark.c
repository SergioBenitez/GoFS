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

  if (real) *real = t.tv_sec + t.tv_usec * 1e-6;
  if (user) *user = r.ru_utime.tv_sec + r.ru_utime.tv_usec * 1e-6;
  if (sys) *sys = r.ru_stime.tv_sec + r.ru_stime.tv_usec * 1e-6;
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
bench_pause(Benchmark *b) {
  stop_timer(b);
  agg_timers(b);
}

void
bench_resume(Benchmark *b) {
  start_timer(b);
}

void
print_results(Benchmark *b) {
  printf("Reps: %d, Total time: %f\n\n", b->reps, b->real);

  /* printf("Totals:\n"); */
  /* printf("Real:\t%fs\n", b->real); */
  /* printf("User:\t%fs\n", b->user); */
  /* printf("Sys:\t%fs\n\n", b->sys); */

  double avg_real = b->real / b->reps;
  double avg_user = b->user / b->reps;
  double avg_sys = b->sys / b->reps;

  printf("Averages:\n");
  printf("Real:\t%6.5fs%15.2fns\n", avg_real, avg_real * 1e9);
  printf("User:\t%6.5fs%15.2fns\n", avg_user, avg_user * 1e9);
  printf("Sys:\t%6.5fs%15.2fns\n", avg_sys, avg_sys * 1e9);
}

void
benchmark(char *name, bench_func f, double min_time) {
  printf("------------------------------\n");
  printf("Running '%s'...", name);

  Benchmark b;
  reset_timer(&b);
  for (double time = 0; time < min_time; time = b.real) {
    bench_resume(&b);
    f(&b);
    bench_pause(&b);
  }

  printf("Done.\n");
  print_results(&b);
  printf("------------------------------\n\n");
}
