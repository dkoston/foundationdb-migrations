#!/bin/bash

# Quick and dirty script to make sure all included migrations run


ALLOK=1

echo "1..3"

# Test that sql transaction migration worked

echo "SELECT id FROM sql;" > test.sql
ROW=$(psql -h postgresql-db -U test database < test.sql 2>&1)
if [[ $ROW == *"(1 row)"* ]]; then
  echo "ok 1 - sql with transactions"
else
  echo "not ok"
  ALLOK=0
fi


# Tests that sql no transaction migration worked

echo "SELECT id FROM post;" > test2.sql
ROW=$(psql -h postgresql-db -U test database < test.sql 2>&1)
if [[ $ROW == *"(1 row)"* ]]; then
  echo "ok 2 - sql no transactions"
else
  echo "not ok"
  ALLOK=0
fi

# Tests that JS migration worked

JS=$(cat  node_migration_complete 2>&1)
if [[ $JS == "1" ]]; then
  echo "ok 3 - javascript"
else
  echo "not ok"
  ALLOK=0
fi


if [ $ALLOK -eq 1 ]; then
    echo "ok"
    echo "All tests successful."
    echo "Files=1, Tests=3,  0 wallclock secs ( 0.04 usr  0.01 sys +  0.04 cusr  0.02 csys =  0.11 CPU)"
    echo "Result: PASS"
else
    echo "ERR!"
fi


