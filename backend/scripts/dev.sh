#!/usr/bin/env bash

set -e

HTTP_ADDR="${HTTP_ADDR:-:18080}" go run ./cmd/server
