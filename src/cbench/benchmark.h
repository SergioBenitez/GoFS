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

void benchmark(char *name, void f(), double min_time);
