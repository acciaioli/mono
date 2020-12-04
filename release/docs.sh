#!/bin/bash

# exit when any command fails
set -e

if [[ -z "${OUTFILE}" ]]; then
  echo "error: OUTFILE not set" >&2
  exit 1
fi

if ! [[ -x "$(command -v mono)" ]]; then
  echo "error: mono is not installed" >&2
  exit 1
fi

rm -f "${OUTFILE}"
echo "outfile: ${OUTFILE}"

declare -a CMDS=(
  "mono --version"
  "mono list"
  "mono build --clean"
  "mono build"
  "mono push"
  "mono list"
  "echo foo bar > examples/python-service/a-new-file.txt"
  "mono list"
  "mono build --clean"
  "mono build"
  "mono push"
  "mono list"
  )

{
  echo "### mono usage"
  echo ""
  echo "\`\`\`"
  echo ""
} >> "${OUTFILE}"

for cmd in "${CMDS[@]}"; do
  {
    echo "▶ $cmd"
    eval "$cmd"
    echo ""
  } >> "${OUTFILE}"
done

{
  echo "\`\`\`"
  echo ""
} >> "${OUTFILE}"

rm -f examples/python-service/a-new-file.txt
