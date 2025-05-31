import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { quizAPI } from '../utils/api';
import type { Quiz, Question, User, SubmitAnswerResponse } from '../types/quiz';

const QuizParticipant: React.FC = () => {
  const { quiz_id } = useParams<{ quiz_id: string }>();
  const [quiz, setQuiz] = useState<Quiz | null>(null);
  const [user, setUser] = useState<User | null>(null);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [selectedAnswer, setSelectedAnswer] = useState<string>('');
  const [score, setScore] = useState(0);
  const [answeredQuestions, setAnsweredQuestions] = useState<Set<string>>(new Set());
  const [lastResult, setLastResult] = useState<SubmitAnswerResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [joining, setJoining] = useState(false);
  const [username, setUsername] = useState('');
  const [showJoinForm, setShowJoinForm] = useState(true);

  useEffect(() => {
    // Check if user is already logged in
    const savedUserId = localStorage.getItem('userId');
    const savedUsername = localStorage.getItem('username');
    
    if (savedUserId && savedUsername) {
      setUser({ id: savedUserId, username: savedUsername });
      setShowJoinForm(false);
      loadQuiz();
    }
  }, []);

  const loadQuiz = async () => {
    if (!quiz_id) return;
    
    try {
      const response = await quizAPI.getQuiz(quiz_id);
      setQuiz(response.data);
    } catch (error) {
      console.error('Failed to load quiz:', error);
    }
  };

  const joinQuiz = async () => {
    if (!quiz_id || !username.trim()) return;
    
    setJoining(true);
    try {
      const response = await quizAPI.joinQuiz(quiz_id, { username: username.trim() });
      const userData = {
        id: response.data.user_id,
        username: response.data.username,
      };
      
      setUser(userData);
      localStorage.setItem('userId', userData.id);
      localStorage.setItem('username', userData.username);
      setShowJoinForm(false);
      await loadQuiz();
    } catch (error) {
      console.error('Failed to join quiz:', error);
    } finally {
      setJoining(false);
    }
  };

  const submitAnswer = async () => {
    if (!quiz_id || !selectedAnswer || !currentQuestion) return;
    
    setLoading(true);
    try {
      const response = await quizAPI.submitAnswer(quiz_id, {
        question_id: currentQuestion.id,
        answer: selectedAnswer,
      });
      
      setLastResult(response.data);
      setScore(response.data.new_score);
      setAnsweredQuestions(prev => new Set([...prev, currentQuestion.id]));
      setSelectedAnswer('');
      
      // Auto advance to next question after 2 seconds
      setTimeout(() => {
        setLastResult(null);
        if (currentQuestionIndex < quiz!.questions.length - 1) {
          setCurrentQuestionIndex(prev => prev + 1);
        }
      }, 2000);
    } catch (error) {
      console.error('Failed to submit answer:', error);
    } finally {
      setLoading(false);
    }
  };

  const currentQuestion: Question | undefined = quiz?.questions[currentQuestionIndex];
  const isCompleted = answeredQuestions.size === quiz?.questions.length;

  if (showJoinForm) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="bg-white p-8 rounded-lg shadow-lg max-w-md w-full">
          <h1 className="text-2xl font-bold text-center mb-6">Join Quiz</h1>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Your Name
              </label>
              <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Enter your name..."
                onKeyPress={(e) => e.key === 'Enter' && joinQuiz()}
              />
            </div>
            <button
              onClick={joinQuiz}
              disabled={joining || !username.trim()}
              className="w-full px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:bg-gray-400"
            >
              {joining ? 'Joining...' : 'Join Quiz'}
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (!quiz) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-500 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading quiz...</p>
        </div>
      </div>
    );
  }

  if (isCompleted) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="bg-white p-8 rounded-lg shadow-lg max-w-md w-full text-center">
          <h1 className="text-3xl font-bold text-green-600 mb-4">Quiz Completed!</h1>
          <p className="text-xl text-gray-700 mb-2">Final Score:</p>
          <p className="text-4xl font-bold text-blue-600 mb-6">{score} points</p>
          <p className="text-gray-600">Thank you for participating, {user?.username}!</p>
          <div className="mt-6">
            <a
              href={`/quiz/${quiz_id}/board`}
              className="inline-block px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
            >
              View Leaderboard
            </a>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center">
      <div className="bg-white p-8 rounded-lg shadow-lg max-w-md w-full">
        <h1 className="text-2xl font-bold text-center mb-6">Quiz</h1>
        <div className="space-y-4">
          <div>
            <p className="text-gray-700 mb-2">Question {currentQuestionIndex + 1}:</p>
            <p className="text-lg font-medium">{currentQuestion?.question_text}</p>
          </div>
          <div>
            {currentQuestion?.options.map((option, index) => (
              <button
                key={index}
                onClick={() => setSelectedAnswer(option)}
                className={`w-full px-4 py-2 rounded-lg ${selectedAnswer === option ? 'bg-blue-500 text-white' : 'bg-gray-200 text-gray-700'
                }`}
              >
                {option}
              </button>
            ))}
          </div>
          <div>
            <button
              onClick={submitAnswer}
              disabled={!selectedAnswer || loading || answeredQuestions.has(currentQuestion?.id || '')}
              className="w-full px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 disabled:bg-gray-400"
            >
              {loading ? 'Submitting...' : answeredQuestions.has(currentQuestion?.id || '') ? 'Answered' : 'Submit Answer'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default QuizParticipant;