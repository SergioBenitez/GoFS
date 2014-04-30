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
typedef void (bench_clean)(void);

void benchmark(char *name, bench_func f, bench_clean c, double min_time);
void bench_pause(Benchmark *b);
void bench_resume(Benchmark *b);
