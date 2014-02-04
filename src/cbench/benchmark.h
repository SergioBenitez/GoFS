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

typedef void (bench_func)(Benchmark *);

void benchmark(char *name, bench_func f, double min_time);
void bench_pause(Benchmark *b);
void bench_resume(Benchmark *b);
