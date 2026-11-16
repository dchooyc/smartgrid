#!/usr/bin/env bash
# Generates five sample asset CSV files into the generated-csv folder, ready
# to upload: one valid and four with different kinds of problem. Each comes
# with an answer sheet saying what the app should report for it.
#
# Every run makes new files with new asset IDs, so they never clash with
# assets that are already stored. To choose the problems yourself, use
# scripts/generate-csv.sh.
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root/backend"

# generate TITLE ARGUMENTS...: generates one file of 3 substations and 48
# assets, broken as the arguments say.
generate() {
  echo "$1"
  shift
  go run ./cmd/gencsv -substations 3 -assets 48 -out "$root/generated-csv" "$@"
  echo
}

generate "1. Valid: every row is accepted"

generate "2. Field errors: missing names, unsupported statuses, bad dates, negative ratings" \
  -missing-name 2 -invalid-status 2 -invalid-date 2 -negative-rating 2

generate "3. Hierarchy errors: parents not found, wrong types of parent, own parent, circular parents" \
  -parent-not-found 2 -wrong-parent-type 2 -self-parent 1 -cycle 1

generate "4. Duplicates: repeated asset IDs (errors) and repeated serial numbers (warnings)" \
  -duplicate-id 2 -duplicate-serial 2

generate "5. Missing column: the whole file is rejected before its rows are checked" \
  -missing-columns 1
