package service

import (
	"fmt"
	"log"
	"quiz-app/internal/model"
	"quiz-app/internal/repository"
	"time"

	"github.com/google/uuid"
)

type QuizService interface {
	CreateQuiz(req *model.CreateQuizRequest) (*model.QuizSession, error)
	JoinQuiz(quizID string, req *model.JoinQuizRequest) (*model.User, error)
	SubmitAnswer(userID, quizID string, req *model.SubmitAnswerRequest) (*model.SubmitAnswerResponse, error)
	GetLeaderboard(quizID string) ([]model.LeaderboardEntry, error)
	GetQuiz(quizID string) (*model.QuizSession, error)

	// New methods for cache management
	InvalidateQuizCache(quizID string) error
	InvalidateLeaderboardCache(quizID string) error
	WarmupCache(quizID string) error
}

type quizService struct {
	quizRepo  repository.QuizRepository
	redisRepo repository.RedisRepository
	wsService WebSocketService
}

func NewQuizService(quizRepo repository.QuizRepository, redisRepo repository.RedisRepository, wsService WebSocketService) QuizService {
	return &quizService{
		quizRepo:  quizRepo,
		redisRepo: redisRepo,
		wsService: wsService,
	}
}

func (s *quizService) CreateQuiz(req *model.CreateQuizRequest) (*model.QuizSession, error) {
	quiz := &model.QuizSession{
		ID:        uuid.New(),
		Title:     req.Title,
		Status:    "waiting",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Create questions
	for i, qReq := range req.Questions {
		question := model.Question{
			ID:            uuid.New(),
			QuizSessionID: quiz.ID,
			QuestionText:  qReq.QuestionText,
			Options:       qReq.Options,
			CorrectAnswer: qReq.CorrectAnswer,
			Points:        qReq.Points,
			Order:         i + 1,
		}
		if question.Points == 0 {
			question.Points = 10
		}
		quiz.Questions = append(quiz.Questions, question)
	}

	// Save to database
	if err := s.quizRepo.CreateQuiz(quiz); err != nil {
		return nil, err
	}

	// Cache in Redis
	s.redisRepo.SetQuizSession(quiz.ID.String(), quiz)

	return quiz, nil
}

func (s *quizService) JoinQuiz(quizID string, req *model.JoinQuizRequest) (*model.User, error) {
	// Check if quiz exists
	quizUUID, err := uuid.Parse(quizID)
	if err != nil {
		return nil, fmt.Errorf("invalid quiz ID")
	}

	_, err = s.quizRepo.GetQuiz(quizUUID)
	if err != nil {
		return nil, fmt.Errorf("quiz not found")
	}

	// Check if user already exists
	existingUser, err := s.quizRepo.GetUserByUsername(req.Username)
	if err == nil {
		return existingUser, nil
	}

	// Create new user
	user := &model.User{
		ID:        uuid.New(),
		Username:  req.Username,
		CreatedAt: time.Now(),
	}

	if err := s.quizRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *quizService) SubmitAnswer(userID, quizID string, req *model.SubmitAnswerRequest) (*model.SubmitAnswerResponse, error) {
	userUUID, _ := uuid.Parse(userID)
	quizUUID, _ := uuid.Parse(quizID)
	questionUUID, _ := uuid.Parse(req.QuestionID)

	// Get quiz and question
	quiz, err := s.GetQuiz(quizID)
	if err != nil {
		return nil, err
	}

	var question *model.Question
	for _, q := range quiz.Questions {
		if q.ID == questionUUID {
			question = &q
			break
		}
	}

	if question == nil {
		return nil, fmt.Errorf("question not found")
	}

	// Check if answer is correct
	isCorrect := question.CorrectAnswer == req.Answer
	points := 0
	if isCorrect {
		points = question.Points
	}

	// Save answer
	answer := &model.UserAnswer{
		ID:         uuid.New(),
		UserID:     userUUID,
		QuestionID: questionUUID,
		Answer:     req.Answer,
		IsCorrect:  isCorrect,
		AnsweredAt: time.Now(),
	}

	if err := s.quizRepo.SaveAnswer(answer); err != nil {
		return nil, err
	}

	log.Printf("🗑️ Invalidating cache due to new answer for quiz %s", quizID)

	go func() {
		if err := s.InvalidateLeaderboardCache(quizID); err != nil {
			log.Printf("⚠️ Failed to invalidate leaderboard cache: %v", err)
		}

	}()

	// Get updated score
	newScore, err := s.quizRepo.GetUserScore(userUUID, quizUUID)
	if err != nil {
		return nil, err
	}

	// Update leaderboard and broadcast if there are viewers
	go s.updateAndBroadcastLeaderboard(quizID)

	return &model.SubmitAnswerResponse{
		Correct:  isCorrect,
		NewScore: newScore,
		Points:   points,
	}, nil
}

func (s *quizService) updateAndBroadcastLeaderboard(quizID string) {
	if !s.wsService.HasLeaderboardViewers(quizID) {
		return
	}

	leaderboard, err := s.GetLeaderboard(quizID)
	if err != nil {
		return
	}
	// Broadcast to WebSocket viewers
	s.wsService.BroadcastLeaderboardUpdate(quizID, leaderboard)
}

func (s *quizService) GetLeaderboard(quizID string) ([]model.LeaderboardEntry, error) {
	quizUUID, err := uuid.Parse(quizID)
	if err != nil {
		return nil, err
	}

	cachedLeaderboard, err := s.redisRepo.GetLeaderboard(quizID)
	if err == nil && len(cachedLeaderboard) > 0 {
		return cachedLeaderboard, nil
	}

	participants, err := s.quizRepo.GetParticipants(quizUUID)
	if err != nil {
		return nil, err
	}

	var leaderboard []model.LeaderboardEntry
	for i, p := range participants {
		leaderboard = append(leaderboard, model.LeaderboardEntry{
			UserID:   p.UserID,
			Username: p.Username,
			Score:    p.Score,
			Rank:     i + 1,
		})
	}

	// Update Redis cache
	s.redisRepo.UpdateLeaderboard(quizID, participants)

	return leaderboard, nil
}

func (s *quizService) GetQuiz(quizID string) (*model.QuizSession, error) {
	quizUUID, err := uuid.Parse(quizID)
	if err != nil {
		return nil, err
	}
	cachedQuiz, err := s.redisRepo.GetQuizSession(quizID)
	if err == nil {
		return cachedQuiz, nil
	}
	quiz, err := s.quizRepo.GetQuiz(quizUUID)
	if err != nil {
		return nil, err
	}
	s.redisRepo.SetQuizSession(quizID, quiz)
	return quiz, nil
}

// ================================================================
// 4. CACHE MANAGEMENT METHODS
// ================================================================

func (s *quizService) InvalidateQuizCache(quizID string) error {
	key := fmt.Sprintf("quiz:%s", quizID)
	if err := s.redisRepo.DeleteKey(key); err != nil {
		return fmt.Errorf("failed to invalidate quiz cache: %w", err)
	}
	log.Printf("🗑️ Invalidated quiz cache for %s", quizID)
	return nil
}

func (s *quizService) InvalidateLeaderboardCache(quizID string) error {
	key := fmt.Sprintf("leaderboard:%s", quizID)
	if err := s.redisRepo.DeleteKey(key); err != nil {
		return fmt.Errorf("failed to invalidate leaderboard cache: %w", err)
	}
	log.Printf("🗑️ Invalidated leaderboard cache for %s", quizID)
	return nil
}

// ================================================================
// 5. CACHE WARMUP - Pre-load popular data
// ================================================================

func (s *quizService) WarmupCache(quizID string) error {
	log.Printf("🔥 Starting cache warmup for quiz %s", quizID)

	startTime := time.Now()

	// Warm up quiz data
	if _, err := s.GetQuiz(quizID); err != nil {
		return fmt.Errorf("failed to warm up quiz cache: %w", err)
	}

	// Warm up leaderboard
	if _, err := s.GetLeaderboard(quizID); err != nil {
		return fmt.Errorf("failed to warm up leaderboard cache: %w", err)
	}

	elapsed := time.Since(startTime)
	log.Printf("🔥 Cache warmup completed for quiz %s (took %v)", quizID, elapsed)

	return nil
}
