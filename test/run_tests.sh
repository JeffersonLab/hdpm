# usage: bash run_tests.sh $EVIO > status.txt && tail -n1 status.txt
# summary.txt summarizes result (passed/failed) of each test
# status.txt shows any failures; the last line indicates the status
echo "Test status"
bash tests.sh $1 >& summary.txt
grep -i 'failed' summary.txt
if test $? -ne 0; then
    echo "SUCCESS: All tests passed."
else
    echo "FAILURE: One or more tests failed."
fi
