#!/bin/bash
set -e
CONTEXT=$1
ENTITY=$2
HDL=$3
DIR="/tmp/go-wbfbd/$CONTEXT/$ENTITY/python"
PYTHON_FILE="../../../tests/co-simulations/$CONTEXT/$ENTITY/python/tb.py"

export PYTHONPATH="$PYTHONPATH:$PWD/../../../tests/co-simulations/common/python/"
export PYTHONPATH="$PYTHONPATH:$PWD/../../wbfbd/$CONTEXT/$ENTITY/python/$HDL"

mkdir -p $DIR
if [ ! -f "$PYTHON_FILE" ]; then
	>&2 echo "$PYTHON_FILE not found"
	exit 1
fi
python3 $PYTHON_FILE \
	/tmp/go-wbfbd/$CONTEXT/$ENTITY/python_"$HDL" \
	/tmp/go-wbfbd/$CONTEXT/$ENTITY/"$HDL"_python \
	> "$DIR/$HDL.log" 2>&1 &
