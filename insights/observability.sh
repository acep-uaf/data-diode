#!/bin/bash

START=`date +%Y-%m-%d-%k:%m:%S`
COLLECTION=`date +%Y-%m-%d`
LOGFILE="insights/ACEP_CAMIO_SEDOID_RESULTS-${COLLECTION}.txt"

function log() {
    echo "diode: $@"
    echo "$@" >> $LOGFILE
}

TARGET_LOCATION="localhost"
VARIABLE_SIZE=(64 128 256 512 1024 2048) # packets
TEST_DURATION=10 # seconds

run_experiment() {
    local packet_size=$1
    echo ">> Running trial with $packet_size bytes"

    ping -c $((TEST_DURATION)) -s $packet_size $TARGET_LOCATION >> $LOGFILE

    sleep 5
}

for iteration in ${VARIABLE_SIZE[@]}; do
    run_experiment $iteration
done

echo ">> Experiment complete. Results: $LOGFILE"
