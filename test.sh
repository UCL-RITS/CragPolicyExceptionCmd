#!/bin/bash

echo "travis_fold:start:Setting up config"

TRAVIS_MYSQL_PORT="$(grep -e '^port=' "$HOME/.my.cnf" | sed -e 's/port=//')"

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
"$EXE" createdb
"$EXE" destroydb
"$EXE" createdb
"$EXE" submit --user="someone" --service="myriad"
"$EXE" submit --user="someone" --service="legion"
"$EXE" submit --user="someone" --service="grace"
"$EXE" submit --user="someone" --service="kathleen"
"$EXE" submit --user="someone" --service="thomas"
"$EXE" submit --user="someone" --service="michael"
"$EXE" list
echo "travis_fold:end:Running tests"
