# Load Balancing algorithms in go

## Round Robin

A uniform/weighted round robin load balancer. Servers are assigned weights, determining selection frequency. Higher weights are favored more. If weights are equal, servers are chosen uniformly.

## Power of Two Choices

A power of two choices (P2C) load balancer. Servers are chosen uniformly at random, and the least loaded server is chosen. This load may be determined by any metric, such as CPU usage, memory usage, or latency.
