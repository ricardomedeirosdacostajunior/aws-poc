#!/bin/bash

. env/.aws.env

function run() {
	for f in "$1"*; do
		echo "executing $f"
		eval "bash $f"
	done
}

run "delete-tables"
