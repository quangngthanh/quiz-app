import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import QuizCreator from './components/QuizCreator';
import QuizParticipant from './components/QuizParticipant';
import Leaderboard from './components/Leaderboard';

function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/" element={<QuizCreator />} />
          <Route path="/quiz/:quiz_id/play" element={<QuizParticipant />} />
          <Route path="/quiz/:quiz_id/board" element={<Leaderboard />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App