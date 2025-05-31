import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { quizAPI } from '../utils/api';
import { useWebSocket } from '../hooks/useWebSocket';
import type { Quiz, LeaderboardEntry, LeaderboardUpdate } from '../types/quiz';

const Leaderboard: React.FC = () => {
  const { quiz_id } = useParams<{ quiz_id: string }>();
  const [quiz, setQuiz] = useState<Quiz | null>(null);
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [lastUpdate, setLastUpdate] = useState<Date>(new Date());

  // WebSocket connection for real-time updates
  const wsUrl = `ws://localhost:${import.meta.env.VITE_PORT}/ws/quiz/${quiz_id}/leaderboard`;
  const { isConnected, error } = useWebSocket(wsUrl, {
    onMessage: (message) => {
      if (message.type === 'leaderboard_update') {
        const update = message as unknown as LeaderboardUpdate;
        setLeaderboard(update.leaderboard);
        setLastUpdate(new Date(update.updated_at));
      }
    },
  });

  useEffect(() => {
    loadQuizAndLeaderboard();
  }, [quiz_id]);

  const loadQuizAndLeaderboard = async () => {
    if (!quiz_id) return;

    try {
      setLoading(true);
      
      // Load quiz info
      const quizResponse = await quizAPI.getQuiz(quiz_id);
      setQuiz(quizResponse.data);

      // Load initial leaderboard
      const leaderboardResponse = await quizAPI.getLeaderboard(quiz_id);
      setLeaderboard(leaderboardResponse.data.leaderboard || []);
      setLastUpdate(new Date());
    } catch (error) {
      console.error('Failed to load data:', error);
    } finally {
      setLoading(false);
    }
  };

  const refreshLeaderboard = async () => {
    if (!quiz_id) return;

    try {
      const response = await quizAPI.getLeaderboard(quiz_id);
      setLeaderboard(response.data.leaderboard || []);
      setLastUpdate(new Date());
    } catch (error) {
      console.error('Failed to refresh leaderboard:', error);
    }
  };

  const getRankColor = (rank: number) => {
    switch (rank) {
      case 1: return 'text-yellow-600 bg-yellow-100';
      case 2: return 'text-gray-600 bg-gray-100';
      case 3: return 'text-orange-600 bg-orange-100';
      default: return 'text-blue-600 bg-blue-100';
    }
  };

  const getRankIcon = (rank: number) => {
    switch (rank) {
      case 1: return '👑';
      case 2: return '🥈';
      case 3: return '🥉';
      default: return '📍';
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-indigo-100 via-purple-50 to-pink-100 flex items-center justify-center backdrop-blur-sm">
        <div className="text-center">
          <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-500 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading leaderboard...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-100 via-purple-50 to-pink-100">
      {/* Header */}
      <div className="bg-white shadow-sm">
        <div className="max-w-6xl mx-auto px-4 py-6">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-3xl font-bold text-gray-800">{quiz?.title} - Leaderboard</h1>
              <p className="text-gray-600 mt-2">
                Live rankings • {leaderboard?.length || 0} participants
              </p>
            </div>
            <div className="text-right">
              <div className="flex items-center space-x-4">
                <div className={`flex items-center space-x-2 px-3 py-1 rounded-full text-sm ${
                  isConnected ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
                }`}>
                  <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`}></div>
                  <span>{isConnected ? 'Live' : 'Disconnected'}</span>
                </div>
                <button
                  onClick={refreshLeaderboard}
                  className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
                >
                  Refresh
                </button>
              </div>
              <p className="text-sm text-gray-500 mt-2">
                Last updated: {lastUpdate.toLocaleTimeString()}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Leaderboard */}
      <div className="max-w-6xl mx-auto px-4 py-8">
        {leaderboard?.length === 0 ? (
          <div className="bg-white rounded-lg shadow-lg p-12 text-center">
            <div className="text-6xl mb-4">🎯</div>
            <h2 className="text-2xl font-bold text-gray-700 mb-2">No participants yet</h2>
            <p className="text-gray-600">Participants will appear here as they join and answer questions.</p>
            <div className="mt-6">
              <p className="text-sm text-gray-500 mb-2">Share the participant link:</p>
              <div className="p-3 bg-gray-100 rounded-lg font-mono text-sm">
                {`${window.location.origin}/quiz/${quiz_id}/play`}
              </div>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            {/* Top 3 Podium */}
            {leaderboard?.length >= 3 && (
              <div className="bg-white rounded-lg shadow-lg p-8 mb-8">
                <h2 className="text-2xl font-bold text-center text-gray-800 mb-8">Top 3</h2>
                <div className="flex justify-center items-end space-x-8">
                  {/* 2nd Place */}
                  {leaderboard[1] && (
                    <div className="text-center w-1/5">
                      <div className="bg-gray-200 rounded-lg p-6 mb-4" style={{ height: '140px' }}>
                        <div className="text-4xl mb-2">🥈</div>
                        <div className="font-bold text-gray-700">{leaderboard[1].username}</div>
                        <div className="text-xl font-bold text-gray-600">{leaderboard[1].score}</div>
                      </div>
                      <div className="text-lg font-bold text-gray-600">#2</div>
                    </div>
                  )}

                  {/* 1st Place */}
                  {leaderboard[0] && (
                    <div className="text-center w-1/5">
                      <div className="bg-yellow-200 rounded-lg p-6 mb-4" style={{ height: '170px' }}>
                        <div className="text-5xl mb-2">👑</div>
                        <div className="font-bold text-yellow-800">{leaderboard[0].username}</div>
                        <div className="text-2xl font-bold text-yellow-700">{leaderboard[0].score}</div>
                      </div>
                      <div className="text-xl font-bold text-yellow-600">#1</div>
                    </div>
                  )}

                  {/* 3rd Place */}
                  {leaderboard[2] && (
                    <div className="text-center w-1/5">
                      <div className="bg-orange-200 rounded-lg p-6 mb-4" style={{ height: '130px' }}>
                        <div className="text-3xl mb-2">🥉</div>
                        <div className="font-bold text-orange-700">{leaderboard[2].username}</div>
                        <div className="text-lg font-bold text-orange-600">{leaderboard[2].score}</div>
                      </div>
                      <div className="text-lg font-bold text-orange-600">#3</div>
                    </div>
                  )}
                </div>
              </div>
            )}

            {/* Full Leaderboard */}
            <div className="bg-white rounded-lg shadow-lg overflow-hidden">
              <div className="px-6 py-4 bg-gray-50 border-b">
                <h3 className="text-lg font-semibold text-gray-800">All Participants</h3>
              </div>
              <div className="divide-y divide-gray-200 max-h-[500px] overflow-y-auto">
                {leaderboard?.map((entry, index) => (
                  <div
                    key={entry.user_id}
                    className={`flex items-center justify-between p-6 hover:bg-gray-50 transition-colors ${
                      index < 3 ? 'bg-gradient-to-r from-gray-50 to-white' : ''
                    }`}
                  >
                    <div className="flex items-center space-x-4">
                      <div className={`w-12 h-12 rounded-full flex items-center justify-center font-bold text-lg ${getRankColor(entry.rank)}`}>
                        <span className="mr-1">{getRankIcon(entry.rank)}</span>
                        {entry.rank}
                      </div>
                      <div>
                        <h4 className="text-lg font-semibold text-gray-800">{entry.username}</h4>
                        <p className="text-sm text-gray-600">Participant</p>
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="text-2xl font-bold text-gray-800">{entry.score}</div>
                      <div className="text-sm text-gray-600">points</div>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="bg-white rounded-lg shadow p-6 text-center">
                <div className="text-3xl font-bold text-blue-600">{leaderboard?.length}</div>
                <div className="text-gray-600">Total Participants</div>
              </div>
              <div className="bg-white rounded-lg shadow p-6 text-center">
                <div className="text-3xl font-bold text-green-600">
                  {leaderboard?.length > 0 ? Math.max(...leaderboard.map(l => l.score)) : 0}
                </div>
                <div className="text-gray-600">Highest Score</div>
              </div>
              <div className="bg-white rounded-lg shadow p-6 text-center">
                <div className="text-3xl font-bold text-purple-600">
                  {leaderboard?.length > 0 ? Math.round(leaderboard.reduce((sum, l) => sum + l.score, 0) / leaderboard.length) : 0}
                </div>
                <div className="text-gray-600">Average Score</div>
              </div>
            </div>
          </div>
        )}

        {error && (
          <div className="mt-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded-lg">
            Connection error: {error}
          </div>
        )}
      </div>
    </div>
  );
};

export default Leaderboard;