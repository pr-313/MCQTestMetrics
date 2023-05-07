# MCQTestMetrics
This is a simple CLI app that provides an interface for a test taker to log their answers, per-question time spent, and cross check with user entered answer keys to get a full picture of one's performance.

![Example Image](https://user-images.githubusercontent.com/46706232/236655596-5440627f-16de-4cb0-a108-10df80ea1768.png)
![Example2](https://user-images.githubusercontent.com/46706232/236655658-eee5bc47-792a-4a87-b770-4ef246176f77.png)
![Example3](https://user-images.githubusercontent.com/46706232/236655659-6efaf2e4-2106-43f2-baa0-f82e4d747f95.png)

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
6. Results will also be dumped into corresponding .csv files for the user to check in MS Excel or any equivalent app.
