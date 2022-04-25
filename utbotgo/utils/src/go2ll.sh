#!/bin/bash
Makefile='
all: $(LLVM_BITCODE_FILE)

$(LLVM_BITCODE_FILE): $(GO_FILES)
	WORK=$(WORK); $(value LLVM_GOC_COMMAND)
'
build_script=$(go build -n . 2>&1)
WORK=$1
output_file=$2
go_files_file=$3
llvm_goc_command=$( \
  echo "$build_script" | \
  grep -E "llvm-goc -c" | \
  sed "s/-O2/-O0/;s/-o [^ ]*/-o ${output_file//\//\\\/} -S -emit-llvm/;s/ [^ ]*_gomod_.go//")
make -f <(echo "$Makefile") \
  LLVM_BITCODE_FILE="$output_file" \
  GO_FILES="$(cat "$go_files_file")" \
  LLVM_GOC_COMMAND="$llvm_goc_command" \
  WORK="$WORK"
