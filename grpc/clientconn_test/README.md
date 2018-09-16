# gRPC client connection test
How to maximize the client concurrency.

## Methods
1. One connection per request
1. Only one client, one connection
1. Fixed-size connection pool & Round Robin Usage
Use a fixed-size connection pool, fetch a connection round-robin.
1. Fixed-size Connection pool
If connection pool has enough connections, take it from pool, otherwise create a new connection.
The connection pool has fixed max capacity, release unused connection to pool when it's not full.

## Performance Comparison(Local Computer)
**Hardware**
MacBook Pro (15-inch, 2016)
Processor 2.7 GHz Intel Core i7
Memory 16GB 2133 MHz LPDDR3

**Press tool, client and server run on the same machine**

### Server just says hello to client
1. One client, one connection
    ```bash
    $ wrk -t2 -c100 -d10s http://localhost:10099/performance
    Running 10s test @ http://localhost:10099/performance
      2 threads and 100 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency     2.85ms    1.24ms  22.24ms   80.47%
        Req/Sec    15.86k     2.86k   21.58k    76.50%
      316674 requests in 10.04s, 38.66MB read
    Requests/sec:  31532.58
    Transfer/sec:      3.85MB
    ```
1. Fixed-size Connection pool
Take connection from pool first, otherwise create a new connection.
    ```bash
    wrk -t2 -c100 -d10s http://localhost:10099/performance
    Running 10s test @ http://localhost:10099/performance
      2 threads and 100 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    36.86ms   34.35ms 334.72ms   83.30%
        Req/Sec     1.58k   660.95     3.89k    67.86%
      31127 requests in 10.04s, 3.80MB read
    Requests/sec:   3100.89
    Transfer/sec:    387.61KB
    ```

### Server sleeps 0.5s and says hello to client

1. 5 client and round robin
    ```bash
    wrk -t2 -c100 -d10s http://localhost:10099/performance
    Running 10s test @ http://localhost:10099/performance
      2 threads and 100 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency   503.43ms    1.95ms 511.62ms   64.35%
        Req/Sec   130.60    119.12   485.00     80.85%
      2000 requests in 10.10s, 250.00KB read
    Requests/sec:    198.04
    Transfer/sec:     24.75KB
    ```

1. fixed connection pool, and expansion
    ```bash
    wrk -t2 -c100 -d10s http://localhost:10099/performance
    Running 10s test @ http://localhost:10099/performance
      2 threads and 100 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency   504.44ms    5.14ms 529.46ms   92.58%
        Req/Sec   119.40     95.83   470.00     85.11%
      1901 requests in 10.07s, 237.62KB read
    Requests/sec:    188.80
    Transfer/sec:     23.60KB
    ```

1. one connection per request
    ```bash
    wrk -t2 -c100 -d10s http://localhost:10099/performance
    Running 10s test @ http://localhost:10099/performance
      2 threads and 100 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency   509.30ms    4.85ms 529.65ms   73.74%
        Req/Sec   188.81    188.12   494.00     71.23%
      1900 requests in 10.06s, 237.50KB read
    Requests/sec:    188.89
    Transfer/sec:     23.61KB
    ```
