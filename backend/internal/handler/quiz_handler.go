package handler

import (
	"net/http"
	"quiz-app/internal/model"
	"quiz-app/internal/service"

	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	quizService service.QuizService
}

func NewQuizHandler(quizService service.QuizService) *QuizHandler {
	return &QuizHandler{quizService: quizService}
}

func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	var req model.CreateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quiz, err := h.quizService.CreateQuiz(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"quiz_id": quiz.ID,
		"title":   quiz.Title,
		"status":  quiz.Status,
	})
}

func (h *QuizHandler) JoinQuiz(c *gin.Context) {
	quizID := c.Param("quiz_id")

	var req model.JoinQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.quizService.JoinQuiz(quizID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":  user.ID,
		"username": user.Username,
	})
}

func (h *QuizHandler) SubmitAnswer(c *gin.Context) {
	quizID := c.Param("quiz_id")
	userID := c.GetHeader("X-User-ID")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID required"})
		return
	}

	var req model.SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.quizService.SubmitAnswer(userID, quizID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *QuizHandler) GetLeaderboard(c *gin.Context) {
	quizID := c.Param("quiz_id")

	leaderboard, err := h.quizService.GetLeaderboard(quizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"leaderboard": leaderboard})
}

func (h *QuizHandler) GetQuiz(c *gin.Context) {
	quizID := c.Param("quiz_id")

	quiz, err := h.quizService.GetQuiz(quizID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}

	// Remove correct answers from response for security
	for i := range quiz.Questions {
		quiz.Questions[i].CorrectAnswer = ""
	}

	c.JSON(http.StatusOK, quiz)
}
