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
pb "Listing..."
"$EXE" list
pb "Printing dump..."
"$EXE" dumpjson
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
pb "Complete."
echo "travis_fold:end:test_running"
