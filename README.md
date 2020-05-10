# Port scanner in go

This port scanner is a multi-threaded go program where we span 'ulimit -n' (maximum open file descriptors) number of go routines and we scan the ip within a port range.
