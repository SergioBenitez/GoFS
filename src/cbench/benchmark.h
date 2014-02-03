typedef struct Benchmark_T {
  double real_start;
  double real_end;

  double user_start;
  double user_end;

  double sys_start;
  double sys_end;

  double real;
  double user;
  double sys;

  int reps;
} Benchmark;

typedef void (bench_func)();

void benchmark(char *name, bench_func f, double min_time);

// Would be nice to have these methods.
// void bench_pause(Benchmark *b);
// void bench_resume(Benchmark *b);

