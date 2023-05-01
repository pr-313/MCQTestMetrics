import csv
import time

def get_time():
    return int(time.time())

class Question:
    def __init__(self, time_taken=0, selected_answer='', correct_answer='', is_correct=''):
        self.time_taken = time_taken
        self.selected_answer = selected_answer
        self.correct_answer = correct_answer
        self.is_correct = is_correct

class GMATTest:
    def __init__(self, num_questions, time_limit, result_file='result_file.csv', answer_key_file='answer_key_file.csv'):
        self.num_questions = num_questions
        self.time_limit = time_limit
        self.questions = []
        self.answers = []
        self.result_file = result_file
        self.answer_key_file = answer_key_file

    def add_question(self, question):
            self.questions.append(question)

    def add_answer(self, question):
            self.answers.append(question)

    def start(self):
        print(f"Starting test with {self.num_questions} questions.")
        main_timer = get_time()
        for i in range(1, self.num_questions + 1):
            start_time = get_time()
            question = Question()
            print(f"Question {i}:")
            # print(f"Select answer: (A, B, C, D, E)")
            answer = input().lower()
            while answer != 'quit' and answer not in ['a', 'b', 'c', 'd', 'e']:
                answer = input("Invalid input. Please enter A, B, C, D, or E, or 'quit' to exit: ").lower()
            if answer == 'quit':
                print(f"Exiting test. Blank answers will be assumed for unanswered questions.")
                break
            end_time = get_time()
            per_q_time = end_time - start_time
            question.selected_answer = answer
            question.time_taken = per_q_time
            self.questions.append(question)
        main_timer = get_time() - main_timer

        # Blank out any unanswered questions
        for i in range(len(self.questions), self.num_questions):
            question = Question(0, '', '', '') # empty values and is_correct=False for unanswered questions
            self.questions.append(question)
        
        self.export_results(main_timer)

    def export_results(self, total_time):
        with open(self.result_file, 'w', newline='') as csvfile:
            writer = csv.writer(csvfile)
            writer.writerow(['Question #', 'Time (s)', 'Selected Answer', 'Correct Answer', 'Is Correct'])
            for i, question in enumerate(self.questions):
                writer.writerow([i+1, question.time_taken, question.selected_answer, question.correct_answer, question.is_correct])
            writer.writerow(['Total time:', total_time])
        print(f"Results exported to {self.result_file}")

    def load_csv_data(self):
        with open(self.answer_key_file, newline='') as csvfile:
            reader = csv.DictReader(csvfile)
            for row in reader:
                self.add_answer(Question(-1, row['Correct Answer'], row['Correct Answer'], 'Yes'))
        with open(self.result_file, newline='') as csvfile:
            reader = csv.DictReader(csvfile)
            self.results = {}
            for row in reader:
                self.add_question(Question(int(row['Time (s)']), row['Selected Answer'], row['Correct Answer'], row['Is Correct']))

    def enter_answer_key(self):
        for i in range(1, self.num_questions + 1):
            print(f"Enter the answer for question {i} or type 'quit' to exit:")
            answer = input().lower()
            while answer != 'quit' and answer not in ['a', 'b', 'c', 'd', 'e']:
                answer = input("Invalid input. Please enter A, B, C, D, or E, or 'quit' to exit: ").lower()
            if answer == 'quit':
                print(f"Exiting answer key entry. Blank answers will be assumed for unanswered questions.")
                break
            question = Question(0, answer, answer, 'Yes') # the correct answer is set as the selected answer for answer key
            self.questions.append(question)
        
        # Blank out any unanswered questions
        for i in range(len(self.questions), self.num_questions):
            question = Question(0, '', '', '') # empty values and is_correct=False for unanswered questions
            self.questions.append(question)
        
        self.export_answer_key('answer_key.csv')

    def export_answer_key(self, filename):
        with open(filename, mode='w', newline='') as file:
            writer = csv.writer(file)
            writer.writerow(['Question', 'Correct Answer'])
            for i, question in enumerate(self.questions, start=1):
                writer.writerow([f'Q{i}', question.selected_answer])
        print(f"Answer key exported to {filename}")

    def evaluate(self):
        self.load_csv_data()
        with open(self.result_file, 'w', newline='') as csvfile:
            writer = csv.writer(csvfile)
            writer.writerow(['Question #', 'Time (s)', 'Selected Answer', 'Correct Answer', 'Is Correct?'])
            for i, question in enumerate(self.questions):
                is_correct = (question.selected_answer == self.answers[i].correct_answer)
                question.is_correct = is_correct
                writer.writerow([i+1, question.time_taken, question.selected_answer, question.correct_answer, is_correct])
        print(f"Results exported to {self.result_file}.")


if __name__ == '__main__':
    import argparse
    parser = argparse.ArgumentParser(description='GMAT test-taking app')
    parser.add_argument('--q_num'       , default='5'              , help='Number of Questions')
    parser.add_argument('--dur'         , default='5'              , help='Duration of test in minutes')
    parser.add_argument('--key'         , default='answer_key.csv' , help='the file to write the test results to')
    parser.add_argument('--mode'        , choices=['test', 'key', 'eval'] , default='test' , help='the mode of operation for the app')
    parser.add_argument('--result_file' , default='results.csv'    , help='the file to write the test results to')
    args = parser.parse_args()

    if args.mode == 'test':
        num_questions = int(args.q_num)
        time_limit = int(args.dur)
        test = GMATTest(num_questions, time_limit)
        test.start()
    elif args.mode == 'key':
        num_questions = int(args.q_num)
        test = GMATTest(num_questions, 0)
        test.enter_answer_key()
    elif args.mode == 'eval':
        test = GMATTest(0, 0)
        test.evaluate()

