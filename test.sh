#!/bin/bash

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

echo "travis_fold:start:Setting up config"

tmpdir="$(mktemp -d)"

echo "travis_fold:start:Running tests"
echo "Creating database schema..."
"$EXE" createdb
echo "Creating blank dump for comparison..."
"$EXE" dumpjson >"$tmpdir/dump-blank.json"
echo "Destroying database schema..."
"$EXE" destroydb
echo "Recreating database schema for rest of tests..."
"$EXE" createdb
echo "Submitting several entries..."
# One for each cluster
"$EXE" submit --username="someone" --service="myriad"
"$EXE" submit --username="someone" --service="legion"
"$EXE" submit --username="someone" --service="grace"
"$EXE" submit --username="someone" --service="kathleen"
"$EXE" submit --username="someone" --service="thomas"
"$EXE" submit --username="someone" --service="michael"
echo "Submitting an invalid entry (invalid clustername)..."
if "$EXE" submit --username="someone" --service="XXXXXXX"; then
  echo "Entry should have failed, instead succeeded."
  false
fi
echo "Listing..."
"$EXE" list
echo "Testing dump and re-import..."
"$EXE" dumpjson >"$tmpdir/dump-before.json"
"$EXE" list >"$tmpdir/dump-before.list"
echo "  Destroying and recreating db..."
"$EXE" destroydb
"$EXE" createdb
echo "  Checking fresh blank dump matches old one..."
"$EXE" dumpjson >"$tmpdir/dump-blank-2.json"
diff -q "$tmpdir/dump-blank.json" "$tmpdir/dump-blank-2.json"
echo "  Reimporting dump..."
"$EXE" importjson <"$tmpdir/dump-before.json"
"$EXE" list >"$tmpdir/dump-after.list"
"$EXE" dumpjson >"$tmpdir/dump-after.json"
echo "  Comparing before and after data..."
diff -q "$tmpdir/dump-before.json" "$tmpdir/dump-after.json"
diff -q "$tmpdir/dump-before.list" "$tmpdir/dump-after.list"
echo "Complete."
echo "travis_fold:end:Running tests"
