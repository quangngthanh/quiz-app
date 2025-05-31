package repository

import (
	"quiz-app/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuizRepository interface {
	CreateQuiz(quiz *model.QuizSession) error
	GetQuiz(id uuid.UUID) (*model.QuizSession, error)
	CreateUser(user *model.User) error
	GetUser(id uuid.UUID) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	SaveAnswer(answer *model.UserAnswer) error
	GetUserScore(userID, quizID uuid.UUID) (int, error)
	GetParticipants(quizID uuid.UUID) ([]model.Participant, error)
	AddParticipant(participant *model.Participant) error
}

type quizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) QuizRepository {
	return &quizRepository{db: db}
}

func (r *quizRepository) CreateQuiz(quiz *model.QuizSession) error {
	return r.db.Create(quiz).Error
}

func (r *quizRepository) GetQuiz(id uuid.UUID) (*model.QuizSession, error) {
	var quiz model.QuizSession
	err := r.db.Preload("Questions").Where("id = ?", id).First(&quiz).Error
	return &quiz, err
}

func (r *quizRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *quizRepository) GetUser(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *quizRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *quizRepository) SaveAnswer(answer *model.UserAnswer) error {
	return r.db.Create(answer).Error
}

func (r *quizRepository) GetUserScore(userID, quizID uuid.UUID) (int, error) {
	var totalScore int
	err := r.db.Model(&model.UserAnswer{}).
		Select("COALESCE(SUM(CASE WHEN is_correct THEN q.points ELSE 0 END), 0)").
		Joins("JOIN questions q ON user_answers.question_id = q.id").
		Where("user_answers.user_id = ? AND q.quiz_session_id = ?", userID, quizID).
		Scan(&totalScore).Error
	return totalScore, err
}

func (r *quizRepository) GetParticipants(quizID uuid.UUID) ([]model.Participant, error) {
	var participants []model.Participant
	err := r.db.Raw(`
        SELECT 
            u.id as user_id,
            u.username,
            ? as quiz_id,
            COALESCE(scores.score, 0) as score,
            u.created_at as joined_at
        FROM users u
        LEFT JOIN (
            SELECT 
                ua.user_id,
                SUM(CASE WHEN ua.is_correct THEN q.points ELSE 0 END) as score
            FROM user_answers ua
            JOIN questions q ON ua.question_id = q.id
            WHERE q.quiz_session_id = ?
            GROUP BY ua.user_id
        ) scores ON u.id = scores.user_id
        WHERE u.id IN (
            SELECT DISTINCT ua.user_id 
            FROM user_answers ua 
            JOIN questions q ON ua.question_id = q.id 
            WHERE q.quiz_session_id = ?
        )
        ORDER BY score DESC
    `, quizID, quizID, quizID).Scan(&participants).Error

	return participants, err
}

func (r *quizRepository) AddParticipant(participant *model.Participant) error {
	// This is handled implicitly when user submits first answer
	return nil
}
