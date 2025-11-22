#!/bin/bash

go test -bench=. -run=^# -count=5 -benchmem