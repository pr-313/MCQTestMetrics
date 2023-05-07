# MCQTestMetrics
This is a simple CLI app that provides an interface for a test taker to log their answers, per-question time spent, and cross check with user entered answer keys to get a full picture of one's performance.

To use this tool:
1. Clone this Git Repo
```sh
git clone https://github.com/pr-313/MCQTestMetrics.git && cd MCQTestMetrics
```
2. Check the help menu of the executable 
```sh
./bin/MCQTest -h

MCQTest

  Flags:
    -h --help       Displays help with available flag, subcommand, and positional value parameters.
    -s --startIdx   Start from Question Index (default: 1)
    -e --stopIdx    Stop at Question Index (default: 10)
    -t --dur        Duration of test (default: 10)
    -k --key        Answer Key mode
    -c --check      Check Responses Against Answer Key
```
3. Run the test in normal mode (not answer key mode) and write your test
4. After concluding, enter the answer key using the `-k` option
5. Finally check the entered answers against the entered key using `-c`
