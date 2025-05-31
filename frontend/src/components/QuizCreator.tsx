import React, { useState } from 'react';
import { quizAPI } from '../utils/api';

interface Question {
  question_text: string;
  options: string[];
  correct_answer: string;
  points: number;
}

const QuizCreator: React.FC = () => {
  const [title, setTitle] = useState('');
  const [questions, setQuestions] = useState<Question[]>([]);
  const [currentQuestion, setCurrentQuestion] = useState<Question>({
    question_text: '',
    options: ['', ''],
    correct_answer: '',
    points: 10,
  });
  const [createdQuizId, setCreatedQuizId] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const addOption = () => {
    setCurrentQuestion({
      ...currentQuestion,
      options: [...currentQuestion.options, ''],
    });
  };

  const updateOption = (index: number, value: string) => {
    const newOptions = [...currentQuestion.options];
    newOptions[index] = value;
    setCurrentQuestion({
      ...currentQuestion,
      options: newOptions,
    });
  };

  const removeOption = (index: number) => {
    if (currentQuestion.options.length > 2) {
      const newOptions = currentQuestion.options.filter((_, i) => i !== index);
      setCurrentQuestion({
        ...currentQuestion,
        options: newOptions,
      });
    }
  };

  const addQuestion = () => {
    if (
      currentQuestion.question_text.trim() &&
      currentQuestion.options.every(opt => opt.trim()) &&
      currentQuestion.correct_answer.trim()
    ) {
      setQuestions([...questions, currentQuestion]);
      setCurrentQuestion({
        question_text: '',
        options: ['', ''],
        correct_answer: '',
        points: 10,
      });
    }
  };

  const createQuiz = async () => {
    if (!title.trim() || questions.length === 0) return;

    setLoading(true);
    try {
      const response = await quizAPI.createQuiz({
        title,
        questions,
      });
      setCreatedQuizId(response.data.quiz_id);
    } catch (error) {
      console.error('Failed to create quiz:', error);
    } finally {
      setLoading(false);
    }
  };

  if (createdQuizId) {
    return (
      <div className="max-w-2xl mx-auto p-6 bg-white rounded-lg shadow-lg">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-green-600 mb-4">Quiz Created Successfully!</h2>
          <p className="text-gray-600 mb-6">Share these links with participants:</p>
          <div className="space-y-4">
            <div className="p-4 bg-blue-50 rounded-lg">
              <h3 className="font-semibold text-blue-800">Participant URL:</h3>
              <p className="text-blue-600 font-mono">{`${window.location.origin}/quiz/${createdQuizId}/play`}</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100 py-8">
      <div className="max-w-2xl mx-auto bg-white p-8 rounded-lg shadow-lg">
        <h1 className="text-2xl font-bold mb-6 text-center">Create a New Quiz</h1>
        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">Quiz Title</label>
          <input
            type="text"
            value={title}
            onChange={e => setTitle(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter quiz title..."
          />
        </div>

        <div className="mb-8">
          <h2 className="text-lg font-semibold mb-4">Add Question</h2>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">Question Text</label>
            <input
              type="text"
              value={currentQuestion.question_text}
              onChange={e => setCurrentQuestion({ ...currentQuestion, question_text: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg"
              placeholder="Enter question..."
            />
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">Options</label>
            {currentQuestion.options.map((option, idx) => (
              <div key={idx} className="flex items-center mb-2">
                <input
                  type="text"
                  value={option}
                  onChange={e => updateOption(idx, e.target.value)}
                  className="flex-1 px-3 py-2 border border-gray-300 rounded-lg"
                  placeholder={`Option ${idx + 1}`}
                />
                <button
                  type="button"
                  onClick={() => removeOption(idx)}
                  disabled={currentQuestion.options.length <= 2}
                  className="ml-2 px-2 py-1 bg-red-500 text-white rounded disabled:bg-gray-300"
                >
                  Remove
                </button>
              </div>
            ))}
            <button
              type="button"
              onClick={addOption}
              className="mt-2 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Add Option
            </button>
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">Correct Answer</label>
            <select
              value={currentQuestion.correct_answer}
              onChange={e => setCurrentQuestion({ ...currentQuestion, correct_answer: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg"
            >
              <option value="">Select correct answer</option>
              {currentQuestion.options.map((option, idx) => (
                <option key={idx} value={option}>
                  {option || `Option ${idx + 1}`}
                </option>
              ))}
            </select>
          </div>
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 mb-2">Points</label>
            <input
              type="number"
              min={1}
              value={currentQuestion.points}
              onChange={e => setCurrentQuestion({ ...currentQuestion, points: Number(e.target.value) })}
              className="w-32 px-3 py-2 border border-gray-300 rounded-lg"
            />
          </div>
          <button
            type="button"
            onClick={addQuestion}
            disabled={
              !currentQuestion.question_text.trim() ||
              currentQuestion.options.some(opt => !opt.trim()) ||
              !currentQuestion.correct_answer.trim()
            }
            className="px-6 py-2 bg-green-500 text-white rounded hover:bg-green-600 disabled:bg-gray-300"
          >
            Add Question
          </button>
        </div>

        {/* Danh sách câu hỏi đã thêm */}
        {questions.length > 0 && (
          <div className="mb-8">
            <h2 className="text-lg font-semibold mb-2">Questions Added</h2>
            <ul className="list-decimal list-inside space-y-2">
              {questions.map((q, idx) => (
                <li key={idx}>
                  <span className="font-medium">{q.question_text}</span> ({q.points} pts)
                </li>
              ))}
            </ul>
          </div>
        )}

        <button
          type="button"
          onClick={createQuiz}
          disabled={!title.trim() || questions.length === 0 || loading}
          className="w-full px-6 py-3 bg-blue-600 text-white rounded-lg font-bold hover:bg-blue-700 disabled:bg-gray-400"
        >
          {loading ? 'Creating...' : 'Create Quiz'}
        </button>
      </div>
    </div>
  );
};

export default QuizCreator;