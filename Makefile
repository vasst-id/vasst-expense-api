#!/bin/bash

all:
	@make build && make run

build:
	@echo "Building vasst-expense-api project..."
	@go build -o vasst-expense-api ./cmd/api/
	@echo "Done."

run:
	@echo "Running vasst-expense-api binary..."
	@./vasst-expense-api