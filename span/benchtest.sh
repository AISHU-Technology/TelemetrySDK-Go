#!/bin/bash

go test -benchmem -run=^$  -bench . span/benchmarks -v
