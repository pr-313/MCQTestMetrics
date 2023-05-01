#!/usr/bin/python3
from threading import Thread
import csv
import time
import os
import sys

timer_running = False

def get_time():
    return int(time.time())

def print_timer(time_limit, main_timer, start_time):
    global timer_running
    timer_running = True
    while timer_running:
        question_time_spent = int((get_time() - start_time))
        time_remaining = int((time_limit*60 - (get_time() - main_timer))/60)
        sys.stdout.write(f"\rTime spent for Question : {question_time_spent}s | Time remaining for Test: {time_remaining}m")
        sys.stdout.flush()
        time.sleep(1)

class Question:
    def __init__(self, idx=1, time_taken=0, selected_answer='', correct_answer='', is_correct=''):
        self.idx = idx
        self.time_taken = time_taken
        self.selected_answer = selected_answer
        self.correct_answer = correct_answer
        self.is_correct = is_correct

class GMATTest:
    def __init__(self, time_limit, result_file='result_file.csv', answer_key_file='answer_key_file.csv', start_idx=1, stop_idx=2):
        self.time_limit = time_limit
        self.questions = []
        self.answers = []
        self.result_file = result_file
        self.answer_key_file = answer_key_file
        self.start_idx = start_idx
        self.stop_idx = stop_idx
        self.num_questions = self.stop_idx - self.start_idx + 1

    def add_question(self, question):
            self.questions.append(question)

    def add_answer(self, question):
            self.answers.append(question)

    def start(self):
        global timer_running
        print(f"Starting test with {self.num_questions} questions.")
        main_timer = get_time()
        for i in range(self.start_idx, self.stop_idx + 1):
            start_time = get_time()
            timer_thread = Thread(target=print_timer, args=(self.time_limit, main_timer, start_time))
            question = Question(i)
            timer_thread.start()
            print()
            print(f"Question {i}:")
            answer = input().lower()
            while answer != 'quit' and answer not in ['a', 'b', 'c', 'd', 'e']:
                answer = input("Invalid input. Please enter A, B, C, D, or E, or 'quit' to exit: ").lower()
            if answer == 'quit':
                print(f"Exiting test. Blank answers will be assumed for unanswered questions.")
                break
            timer_thread.join()
            end_time = get_time()
            per_q_time = end_time - start_time
            question.selected_answer = answer
            question.time_taken = per_q_time
            self.questions.append(question)
            timer_running = False
            
        main_timer = get_time() - main_timer

        # Blank out any unanswered questions
        for i in range(self.start_idx+len(self.questions), self.stop_idx+1):
            question = Question(i, 0, '', '', '') # empty values and is_correct=False for unanswered questions
            self.questions.append(question)
        
        self.export_results(main_timer)

    def export_results(self, total_time):
        with open(self.result_file, 'w', newline='') as csvfile:
            writer = csv.writer(csvfile)
            writer.writerow(['Question #', 'Time (s)', 'Selected Answer', 'Correct Answer', 'Is Correct'])
            for i, question in enumerate(self.questions):
                writer.writerow([question.idx, question.time_taken, question.selected_answer, question.correct_answer, question.is_correct])
            print('Total time: '+ str(total_time))
        print(f"Results exported to {self.result_file}")

    def enter_answer_key(self):
        for i in range(self.start_idx, self.stop_idx + 1):
            print(f"Enter the answer for question {i} or type 'quit' to exit:")
            answer = input().lower()
            while answer != 'quit' and answer not in ['a', 'b', 'c', 'd', 'e']:
                answer = input("Invalid input. Please enter A, B, C, D, or E, or 'quit' to exit: ").lower()
            if answer == 'quit':
                print(f"Exiting answer key entry. Blank answers will be assumed for unanswered questions.")
                break
            answer = Question(i, 0, answer, answer, 'Yes') # the correct answer is set as the selected answer for answer key
            self.answers.append(answer)
        
        # Blank out any unanswered questions
        for i in range(self.start_idx+len(self.answers), self.stop_idx+1):
            answer = Question(i, 0, '', '', '') # empty values and is_correct=False for unanswered questions
            self.answers.append(answer)
        
        self.export_answer_key('answer_key.csv')

    def export_answer_key(self, filename):
        with open(filename, mode='w', newline='') as file:
            writer = csv.writer(file)
            writer.writerow(['Question #', 'Correct Answer'])
            for i, question in enumerate(self.answers):
                writer.writerow([question.idx, question.selected_answer])
        print(f"Answer key exported to {filename}")

    def load_csv_data(self):
        if not (os.path.exists(self.result_file)):
            print(f"Does not exist: {self.result_file}")
            return
        if not (os.path.exists(self.answer_key_file)):
            print(f"Does not exist: {self.result_file}")
            return
        with open(self.answer_key_file, newline='') as csvfile:
            reader = csv.DictReader(csvfile)
            for row in reader:
                self.add_answer(Question(int(row['Question #']), 0, row['Correct Answer'], row['Correct Answer'], 'Yes'))
        with open(self.result_file, newline='') as csvfile:
            reader = csv.DictReader(csvfile)
            self.results = {}
            for row in reader:
                self.add_question(Question(int(row['Question #']), int(row['Time (s)']), row['Selected Answer'], row['Correct Answer'], row['Is Correct']))

    def evaluate(self):
        self.load_csv_data()
        with open(self.result_file, 'w', newline='') as csvfile:
            writer = csv.writer(csvfile)
            writer.writerow(['Question #', 'Time (s)', 'Selected Answer', 'Correct Answer', 'Is Correct'])
            for j, answer in enumerate(self.answers):
                for i, question in enumerate(self.questions):
                    if (question.idx != answer.idx): continue
                    if (question.selected_answer == answer.correct_answer):
                        is_correct = 'Yes'
                    else: 
                        is_correct = 'No'
                    question.is_correct = is_correct
                    question.correct_answer = answer.correct_answer
                    writer.writerow([question.idx, question.time_taken, question.selected_answer, question.correct_answer, is_correct])
        print(f"Results exported to {self.result_file}.")

    def pretty_print_results(self):
        import csv
        """
        Prints the results csv file in a pretty format
        """
        if not (os.path.exists(self.result_file)):
            print(f"Does not exist: {self.result_file}")
            return
        with open(self.result_file, mode='r') as csv_file:
            csv_reader = csv.DictReader(csv_file)
            print("Question # | Time (s) | Selected Answer | Correct Answer | Is Correct")
            print("---------------------------------------------------------------------")
            for row in csv_reader:
                print(f"{row['Question #']:>10} | {row['Time (s)']:>8} | {row['Selected Answer']:>15} | {row['Correct Answer']:>14} | {row['Is Correct']:>10}")



if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser(description='GMAT test-taking app')
    parser.add_argument('--start_idx'   , type=int , default=1          , help='Questions start index')
    parser.add_argument('--stop_idx'    , type=int , default=2          , help='Questions stop index')
    parser.add_argument('--dur'         , type=int , default=5          , help='Duration of test in minutes')
    parser.add_argument('--key'         , default='answer_key_file.csv' , help='the file to write the test results to')
    parser.add_argument('--result_file' , default='result_file.csv'    , help='the file to write the test results to')
    parser.add_argument('--mode'        , choices=['test', 'key', 'eval'] , default='' , help='the mode of operation for the app')
    args = parser.parse_args()

    if args.mode == 'test':
        test = GMATTest(args.dur, args.result_file, args.key, args.start_idx, args.stop_idx)
        test.start()
    elif args.mode == 'key':
        test = GMATTest(0, args.result_file, args.key, args.start_idx, args.stop_idx)
        test.enter_answer_key()
    elif args.mode == 'eval':
        test = GMATTest(0, args.result_file, args.key, args.start_idx, args.stop_idx)
        test.evaluate()
        test.pretty_print_results()
    else:
        test = GMATTest(0, args.result_file, args.key, args.start_idx, args.stop_idx)
        test.pretty_print_results()

