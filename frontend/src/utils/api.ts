import axios from 'axios';

const API_BASE_URL = `http://localhost:${import.meta.env.VITE_PORT}/api`;

export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add user ID to requests if available
api.interceptors.request.use((config) => {
  const userId = localStorage.getItem('userId');
  if (userId) {
    config.headers['X-User-ID'] = userId;
  }
  return config;
});

export const quizAPI = {
  createQuiz: (data: any) => api.post('/quiz', data),
  getQuiz: (quizId: string) => api.get(`/quiz/${quizId}`),
  joinQuiz: (quizId: string, data: any) => api.post(`/quiz/${quizId}/join`, data),
  submitAnswer: (quizId: string, data: any) => api.post(`/quiz/${quizId}/answer`, data),
  getLeaderboard: (quizId: string) => api.get(`/quiz/${quizId}/leaderboard`),
};