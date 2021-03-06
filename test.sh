#!/bin/bash

function pb() {
  printf "\e[34m%s\e[0m\n" "$*"
}
function pr() {
  printf "\e[31m%s\e[0m\n" "$*"
}

echo "travis_fold:start:Setting up config"

TRAVIS_MYSQL_PORT="$(grep -e '^port =' "$HOME/.my.cnf" | sed -e 's/port = //')"

# The location of the binary is somewhere in /tmp after running install.sh
#  but we don't know where, since it's a mktemp directory.
EXE="$(find /tmp/tmp.* -name "exceptions")"

cat >"$HOME/.exceptions_db.conf" <<EOF
{
    "db_type": "mysql",
    "db_connection_string": "travis:@tcp(127.0.0.1:$TRAVIS_MYSQL_PORT)/service_exceptions_test"
}
EOF

echo "travis_fold:start:test_setup"

tmpdir="$(mktemp -d)"

echo "travis_fold:start:test_running"
set -o errexit
pb "Running tests..."
pb "Creating database schema..."
"$EXE" createdb
pb "Creating blank dump for comparison..."
"$EXE" dumpjson >"$tmpdir/dump-blank.json"
pb "Destroying database schema..."
"$EXE" destroydb
pb "Recreating database schema for rest of tests..."
"$EXE" createdb
pb "Submitting several entries..."
# One for each cluster
"$EXE" submit --username="someone" --service="myriad"
"$EXE" submit --username="someone" --service="legion"
"$EXE" submit --username="someone" --service="grace"
"$EXE" submit --username="someone" --service="kathleen"
"$EXE" submit --username="someone" --service="thomas"
"$EXE" submit --username="someone" --service="michael"
pb "Submitting an invalid entry (invalid clustername)..."
if "$EXE" submit --username="someone" --service="XXXXXXX"; then
  pr "Entry should have failed, instead succeeded."
  false
fi
echo "travis_fold:start:listing_dumps"
pb "Listing..."
"$EXE" list
pb "Printing dump..."
"$EXE" dumpjson
echo "travis_fold:end:listing_dumps"
pb "Testing dump and re-import..."
"$EXE" dumpjson >"$tmpdir/dump-before.json"
"$EXE" list >"$tmpdir/dump-before.list"
pb "  Destroying and recreating db..."
"$EXE" destroydb
"$EXE" createdb
pb "  Checking fresh blank dump matches old one..."
"$EXE" dumpjson >"$tmpdir/dump-blank-2.json"
if ! diff -q "$tmpdir/dump-blank.json" "$tmpdir/dump-blank-2.json"; then
  pr "blank dumps before and after should be the same, were different."
  pr "Full diff:"
  diff "$tmpdir/dump-blank.json" "$tmpdir/dump-blank-2.json"
fi
pb "  Reimporting dump..."
"$EXE" importjson <"$tmpdir/dump-before.json"
"$EXE" list >"$tmpdir/dump-after.list"
"$EXE" dumpjson >"$tmpdir/dump-after.json"
pb "  Comparing before and after data..."
if ! diff -q "$tmpdir/dump-before.json" "$tmpdir/dump-after.json"; then
  pr "dumps before and after should be the same, were different."
  pr "Full diff:"
  diff "$tmpdir/dump-before.json" "$tmpdir/dump-after.json"
fi
if ! diff -q "$tmpdir/dump-before.list" "$tmpdir/dump-after.list"; then
  pr "listings before and after should be the same, were different."
  pr "Full diff:"
  diff "$tmpdir/dump-before.list" "$tmpdir/dump-after.list"
fi
pb "  Submitting new entry to create different dump..."
"$EXE" submit --username="someone" --service="michael"
"$EXE" dumpjson >"$tmpdir/dump-after-different.json"
"$EXE" list >"$tmpdir/dump-after-different.list"
if diff -q "$tmpdir/dump-before.json" "$tmpdir/dump-after-different.json" >/dev/null; then
  pr "dumps before and after submitting new exception should be different, instead were the same."
  false
fi
if diff -q "$tmpdir/dump-before.list" "$tmpdir/dump-after-different.list" >/dev/null; then
  pr "listings before and after submitting new exception should be different, instead were the same."
  false
fi


# Okay, those were kind of simple.
# Now to test a workflow.
pb "Testing a sample workflow..."
function getprop() {
  grep "^ *$1 *|" \
    | sed -e 's/^[^|]*| //' \
          -e 's/ *$//' \
          -e 's/^ *//'
}
function checkprop() {
  local prop
  prop="$(
    "$EXE" info "$1" \
      | getprop "$2" \
    )"
  if [[ "$prop" != "$3" ]]; then
    echo "Failed: check for $2 in entry $1: expected \"$3\", got \"$prop\""
    return 1
  fi
}

"$EXE" destroydb
"$EXE" createdb
echo "TEST FILE" >"$tmpdir/test_file"
echo " Submitting..."
"$EXE" submit --username=BEEP123 --service=none --comment="ABCDEF" --type=special --submitted=2030-01-15 --starts=2030-01-31 --ends=2030-04-04 --form="$tmpdir/test_file"
echo " Checking username..."; checkprop 1 "Username"  "beep123" # Usernames should force lowercase
echo " Checking dates...";    checkprop 1 "Submitted" "2030-01-15"
                              checkprop 1 "Starts"    "2030-01-31" 
                              checkprop 1 "Ends"      "2030-04-04"
echo " Checking service...";  checkprop 1 "Service"   "none"
echo " Checking type...";     checkprop 1 "Type"      "special"
echo " Checking status...";   checkprop 1 "Status"    "undecided"
echo " Marking as approved..."
"$EXE" approve 1
echo " Marking as implemented..."
"$EXE" implemented 1
echo " Checking status updates..."
[[ "$("$EXE" info 1 | grep -c "Status Change")" == "3" ]]
[[ $("$EXE" info 1 | getprop "Status") == "implemented" ]]
echo " Adding comment..."
"$EXE" comment -c "MNOPQ" 1
echo " Marking as removed..."
"$EXE" remove 1
echo " Checking status updates..."
[[ "$("$EXE" info 1 | grep -c "Status Change")" == "4" ]]
[[ $("$EXE" info 1 | getprop "Status") == "removed" ]]
echo " Checking form attachment..."
"$EXE" form download-for 1
diff -q "test_file" "$tmpdir/test_file"
pb "Complete."
echo "travis_fold:end:test_running"
