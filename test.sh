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

echo "travis_fold:start:Running tests"
echo "Creating database schema..."
"$EXE" createdb
echo "Destroying database schema..."
"$EXE" destroydb
echo "Recreating database schema for rest of tests..."
"$EXE" createdb
echo "Submitting several entries..."
"$EXE" submit --username="someone" --service="myriad"
"$EXE" submit --username="someone" --service="legion"
"$EXE" submit --username="someone" --service="grace"
"$EXE" submit --username="someone" --service="kathleen"
"$EXE" submit --username="someone" --service="thomas"
"$EXE" submit --username="someone" --service="michael"
echo "Listing..."
"$EXE" list
echo "Complete."
echo "travis_fold:end:Running tests"
