export interface Quiz {
    id: string;
    title: string;
    status: 'waiting' | 'active' | 'completed';
    questions: Question[];
    created_at: string;
  }
  
  export interface Question {
    id: string;
    question_text: string;
    options: string[];
    correct_answer?: string;
    points: number;
    order: number;
  }
  
  export interface User {
    id: string;
    username: string;
  }
  
  export interface LeaderboardEntry {
    user_id: string;
    username: string;
    score: number;
    rank: number;
  }
  
  export interface SubmitAnswerResponse {
    correct: boolean;
    new_score: number;
    points: number;
  }
  
  export interface WSMessage {
    type: string;
    data: any;
  }
  
  export interface LeaderboardUpdate {
    type: 'leaderboard_update';
    leaderboard: LeaderboardEntry[];
    updated_at: string;
  }