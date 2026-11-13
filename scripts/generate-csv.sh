#!/usr/bin/env bash
# Asks what kind of sample asset CSV to generate, then generates it into the
# generated-csv folder along with an answer sheet. Press Enter at a question
# to accept the answer shown in [brackets].
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root/backend"

# ask QUESTION DEFAULT: prints the answer, or the default if none is given.
ask() {
  local answer
  read -r -p "$1 [$2] " answer || true
  echo "${answer:-$2}"
}

substations="$(ask "How many substations?" 3)"
# Only suggest a total if the number of substations can be multiplied.
suggested=48
[[ "$substations" =~ ^[0-9]+$ ]] && suggested=$((substations * 16))
assets="$(ask "How many assets in total, substations included?" "$suggested")"

arguments=(-substations "$substations" -assets "$assets" -out "$root/generated-csv")

# A file with a column problem is rejected before its rows are looked at, so
# row errors are only offered when the columns are left alone.
echo
echo "Column problem?"
echo "  1) None"
echo "  2) Missing columns"
echo "  3) Duplicate columns"
echo "  4) Extra columns"
case "$(ask "Choose one" 1)" in
  2) arguments+=(-missing-columns "$(ask "How many columns to remove?" 1)") ;;
  3) arguments+=(-duplicate-columns "$(ask "How many columns to repeat?" 1)") ;;
  4) arguments+=(-extra-columns "$(ask "How many columns to add?" 1)") ;;
  *)
    # The generator lists its kinds of row error as "flag<TAB>title" lines.
    mapfile -t kinds < <(go run ./cmd/gencsv -list)

    echo
    echo "Row errors, in the order the app checks them:"
    for i in "${!kinds[@]}"; do
      printf "  %2d) %s\n" "$((i + 1))" "${kinds[$i]#*$'\t'}"
    done
    read -r -p "Numbers of the ones to include, separated by spaces [none] " chosen || true

    if [ -n "$chosen" ]; then
      count="$(ask "How many of each?" 3)"
      for number in $chosen; do
        if ! [[ "$number" =~ ^[0-9]+$ ]] || [ "$number" -lt 1 ] || [ "$number" -gt "${#kinds[@]}" ]; then
          echo "\"$number\" is not one of the numbers listed." >&2
          exit 1
        fi
        arguments+=("-${kinds[$((number - 1))]%%$'\t'*}" "$count")
      done
    fi
    ;;
esac

echo
go run ./cmd/gencsv "${arguments[@]}"
