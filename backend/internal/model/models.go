package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Username  string    `json:"username" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at"`
}

type QuizSession struct {
	ID        uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title     string     `json:"title" gorm:"not null"`
	Status    string     `json:"status" gorm:"default:'waiting'"` // waiting, active, completed
	Questions []Question `json:"questions" gorm:"foreignKey:QuizSessionID"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt time.Time  `json:"expires_at"`
}

type Question struct {
	ID            uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	QuizSessionID uuid.UUID      `json:"quiz_session_id"`
	QuestionText  string         `json:"question_text" gorm:"not null"`
	Options       pq.StringArray `json:"options" gorm:"type:text[]"`
	CorrectAnswer string         `json:"correct_answer" gorm:"not null"`
	Points        int            `json:"points" gorm:"default:10"`
	Order         int            `json:"order"`
}

type UserAnswer struct {
	ID         uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid"`
	QuestionID uuid.UUID `json:"question_id" gorm:"type:uuid"`
	Answer     string    `json:"answer"`
	IsCorrect  bool      `json:"is_correct"`
	AnsweredAt time.Time `json:"answered_at"`
}

type Participant struct {
	UserID   uuid.UUID `json:"user_id"`
	QuizID   uuid.UUID `json:"quiz_id"`
	Username string    `json:"username"`
	Score    int       `json:"score"`
	JoinedAt time.Time `json:"joined_at"`
}

// Request/Response DTOs
type CreateQuizRequest struct {
	Title     string            `json:"title" binding:"required"`
	Questions []QuestionRequest `json:"questions" binding:"required,min=1"`
}

type QuestionRequest struct {
	QuestionText  string   `json:"question_text" binding:"required"`
	Options       []string `json:"options" binding:"required,min=2"`
	CorrectAnswer string   `json:"correct_answer" binding:"required"`
	Points        int      `json:"points"`
}

type JoinQuizRequest struct {
	Username string `json:"username" binding:"required"`
}

type SubmitAnswerRequest struct {
	QuestionID string `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
}

type SubmitAnswerResponse struct {
	Correct  bool `json:"correct"`
	NewScore int  `json:"new_score"`
	Points   int  `json:"points"`
}

type LeaderboardEntry struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Score    int       `json:"score"`
	Rank     int       `json:"rank"`
}

// WebSocket Messages
type WSMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type LeaderboardUpdate struct {
	Type        string             `json:"type"`
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
	UpdatedAt   time.Time          `json:"updated_at"`
}
