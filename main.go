package main

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-api/controllers"
	"go-api/initializers"
	"go-api/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func isOnlyNumber(text string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(text)
}

func isOnlyAlphabet(text string) bool {
	re := regexp.MustCompile(`^[a-zA-Z .,!?]+$`)
	return re.MatchString(text)
}

func isOnlyLink(text string) bool {
	return strings.HasPrefix("my string", text)
}

func hasMaxLenght(text string, maxLenght int) bool {
	return len(text) <= maxLenght
}

type Attachment struct {
	gorm.Model
	AttachmentID       uint              `gorm:"column:attachment_id;primaryKey;autoIncrement;unique" json:"attachment_id"`
	LearningmaterialID uint              `json:"learning_material_id"`
	Learningmaterial   Learning_Material `gorm:"references:LearningmaterialID"`
	Source             string            `json:"source"`
	Description        string            `json:"description"`
}

func createAttachment(c *gin.Context) {
	var newAttachment Attachment
	if err := c.ShouldBindJSON(&newAttachment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var learning_material Learning_Material
	if err := db.First(&learning_material, newAttachment.LearningmaterialID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Learning Material with that ID not found"})
		return
	}

	db.Create(&newAttachment)
	c.JSON(http.StatusCreated, newAttachment)
}

func getAttachments(c *gin.Context) {
	var attachments []Attachment
	db.Preload("Learningmaterial").Find(&attachments)
	c.JSON(http.StatusOK, attachments)
}

func getAttachmentByLearningMaterialID(c *gin.Context) {
	id := c.Param("id")
	var attachment []Attachment
	if err := db.Preload("Learningmaterial").Find(&attachment, "learningmaterial_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attachment not found"})
		return
	}
	c.JSON(http.StatusOK, attachment)
}

func updateAttachment(c *gin.Context) {
	id := c.Param("id")
	var attachment Attachment
	if err := db.Preload("Learningmaterial").First(&attachment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attachment not found"})
		return
	}

	var input struct {
		Source      *string `json:"source"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData := map[string]interface{}{}

	if !isOnlyLink(*input.Source) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Source must be a link, starts with http:// or https://"})
		return
	}

	if input.Source != nil {
		updateData["source"] = *input.Source
	}
	if input.Description != nil {
		updateData["description"] = *input.Description
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	db.Model(&attachment).Updates(updateData)
	c.JSON(http.StatusOK, attachment)
}

// =============================

type Interest struct {
	gorm.Model
	InterestID    uint   `gorm:"column:interest_id;primaryKey;autoIncrement;unique" json:"interest_id"`
	Interest_name string `json:"interest_name"`
	Description   string `json:"description"`
}

func createInterest(c *gin.Context) {
	var newInterest Interest
	if err := c.ShouldBindJSON(&newInterest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&newInterest)
	c.JSON(http.StatusCreated, newInterest)
}

func getInterests(c *gin.Context) {
	var interests []Interest
	db.Find(&interests)
	c.JSON(http.StatusOK, interests)
}

func getInterestByID(c *gin.Context) {
	id := c.Param("id")
	var interest Interest
	if err := db.First(&interest, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interest not found"})
		return
	}

	c.JSON(http.StatusOK, interest)
}

func updateInterest(c *gin.Context) {
	id := c.Param("id")
	var interest Interest
	if err := db.First(&interest, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interest not found"})
		return
	}

	var input struct {
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData := map[string]interface{}{}

	if input.Description != nil {
		updateData["description"] = *input.Description
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	db.Model(&interest).Updates(updateData)
	c.JSON(http.StatusOK, interest)
}

// ==================================

type Learning_Material struct {
	gorm.Model
	LearningmaterialID uint    `gorm:"column:learning_material_id;primaryKey;autoIncrement;unique" json:"learning_material_id"`
	SubjectID          uint    `json:"subject_id"`
	Subject            Subject `gorm:"references:SubjectID"`
	Content            string  `json:"content"`
}

func createLearningMaterial(c *gin.Context) {
	var newLearningMaterial Learning_Material
	if err := c.ShouldBindJSON(&newLearningMaterial); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var subject Subject
	if err := db.First(&subject, newLearningMaterial.SubjectID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subject with that ID not found"})
		return
	}

	db.Create(&newLearningMaterial)
	c.JSON(http.StatusCreated, newLearningMaterial)
}

func getLearningMaterials(c *gin.Context) {
	var learningMaterial []Learning_Material
	db.Preload("Subject").Find(&learningMaterial)
	c.JSON(http.StatusOK, learningMaterial)
}

func getLearningMaterialByID(c *gin.Context) {
	id := c.Param("id")
	var learningMaterial Learning_Material
	if err := db.Preload("Subject").First(&learningMaterial, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Learning Material not found"})
		return
	}

	c.JSON(http.StatusOK, learningMaterial)
}

func getLearningMaterialBySubjectID(c *gin.Context) {
	id := c.Param("id")
	var learningMaterial []Learning_Material
	if err := db.Preload("Subject").Find(&learningMaterial, "subject_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Learning Material not found"})
		return
	}

	c.JSON(http.StatusOK, learningMaterial)
}

func updateLearningMaterial(c *gin.Context) {
	id := c.Param("id")
	var learningMaterial Learning_Material
	if err := db.Preload("Subject").First(&learningMaterial, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Learning Material not found"})
		return
	}

	var input struct {
		Content *string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData := map[string]interface{}{}

	if input.Content != nil {
		updateData["content"] = *input.Content
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	db.Model(&learningMaterial).Updates(updateData)
	c.JSON(http.StatusOK, learningMaterial)
}

func deleteLearningMaterial(c *gin.Context) {
	id := c.Param("id")
	var learningMaterial Learning_Material
	if err := db.First(&learningMaterial, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Learning Material not found"})
		return
	}

	db.Delete(&learningMaterial)
	c.JSON(http.StatusOK, gin.H{"message": "Learning Material deleted"})
}

// ==================================

type Placement_Test_Answer struct {
	gorm.Model
	PlacementtestanswerID uint           `gorm:"column:placement_test_answer_id;primaryKey;autoIncrement;unique" json:"placement_test_answer_id"`
	PlacementtestID       uint           `json:"placement_test_id"`
	Placementtest         Placement_Test `gorm:"references:PlacementtestID"`
	StudentID             uint           `json:"student_id"`
	Student               Student        `gorm:"references:StudentID"`
	Student_answer        string         `json:"student_answer"`
}

func createPlacementTestAnswer(c *gin.Context) {
	var newPlacementTestAnswer Placement_Test_Answer
	if err := c.ShouldBindJSON(&newPlacementTestAnswer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var placementTest Placement_Test
	if err := db.First(&placementTest, newPlacementTestAnswer.PlacementtestID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Placement Test with that ID not found"})
		return
	}

	var student Student
	if err := db.First(&student, newPlacementTestAnswer.PlacementtestID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Student with that ID not found"})
		return
	}

	db.Create(&newPlacementTestAnswer)
	c.JSON(http.StatusCreated, newPlacementTestAnswer)
}

func getPlacementTestAnswers(c *gin.Context) {
	var placementTestAnswer []Placement_Test_Answer
	db.Preload("Placement_Test").Find(&placementTestAnswer)
	c.JSON(http.StatusOK, placementTestAnswer)
}

func getPlacementTestAnswerByStudentID(c *gin.Context) {
	id := c.Param("id")
	var placementTestAnswer []Placement_Test_Answer
	if err := db.Preload("Student").Preload("Placementtest").Find(&placementTestAnswer, "student_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Placement Test Answer not found"})
		return
	}

	c.JSON(http.StatusOK, placementTestAnswer)
}

func getPlacementTestAnswerByPlacementTestID(c *gin.Context) {
	id := c.Param("id")
	var placementTestAnswer []Placement_Test_Answer
	if err := db.Preload("Student").Preload("Placementtest").Find(&placementTestAnswer, "placementtest_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Placement Test Answer not found"})
		return
	}

	c.JSON(http.StatusOK, placementTestAnswer)
}

func updatePlacementTestAnswer(c *gin.Context) {
	id := c.Param("id")
	var placementTestAnswer Placement_Test_Answer
	if err := db.Preload("Student").Preload("Placement_Test").First(&placementTestAnswer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Placement Test Answer not found"})
		return
	}

	var input struct {
		Student_answer *string `json:"student_answer"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData := map[string]interface{}{}

	if input.Student_answer != nil {
		updateData["student_answer"] = *input.Student_answer
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	db.Model(&placementTestAnswer).Updates(updateData)
	c.JSON(http.StatusOK, placementTestAnswer)
}

// ==============================

type Placement_Test_Result struct {
	gorm.Model
	PlacementtestresultID uint      `gorm:"column:placement_test_result_id;primaryKey;autoIncrement;unique" json:"placement_test_result_id"`
	StudentID             uint      `json:"student_id"`
	Student               Student   `gorm:"references:StudentID"`
	InterestID            uint      `json:"interest_id"`
	Interest              Interest  `gorm:"references:InterestID"`
	Score                 int       `json:"score"`
	Test_date             time.Time `json:"test_date"`
}

func createPlacementTestResult(c *gin.Context) {
	var newPlacementTestResult Placement_Test_Result
	if err := c.ShouldBindJSON(&newPlacementTestResult); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPlacementTestResult.Test_date = time.Now()

	var student Student
	if err := db.First(&student, newPlacementTestResult.StudentID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Student not found"})
		return
	}

	var interest Interest
	if err := db.First(&interest, newPlacementTestResult.StudentID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Interest not found"})
		return
	}

	db.Create(&newPlacementTestResult)
	c.JSON(http.StatusCreated, newPlacementTestResult)
}

func getPlacementTestResults(c *gin.Context) {
	var placementTestResult []Placement_Test_Result
	db.Preload("Student").Preload("Interest").Find(&placementTestResult)
	c.JSON(http.StatusOK, placementTestResult)
}

func getPlacementTestResultByStudentID(c *gin.Context) {
	id := c.Param("id")
	var placementTestResult []Placement_Test_Result
	if err := db.Preload("Student").Preload("Interest").Find(&placementTestResult, "student_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Placement Test Result not found"})
		return
	}

	c.JSON(http.StatusOK, placementTestResult)
}

func getPlacementTestResultByInterestID(c *gin.Context) {
	id := c.Param("id")
	var placementTestResult []Placement_Test_Result
	if err := db.Preload("Student").Preload("Interest").Find(&placementTestResult, "interest_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Placement Test Result not found"})
		return
	}

	c.JSON(http.StatusOK, placementTestResult)
}

// =======================

type Placement_Test struct {
	gorm.Model
	PlacementtestID uint     `gorm:"column:placement_test_id;primaryKey;autoIncrement;unique" json:"placement_test_id"`
	Question        string   `json:"question"`
	Correct_answer  string   `json:"correct_answer"`
	Option_a        string   `json:"option_a"`
	Option_b        string   `json:"option_b"`
	Option_c        string   `json:"option_c"`
	Option_d        string   `json:"option_d"`
	InterestID      uint     `json:"interest_id"`
	Interest        Interest `gorm:"references:InterestID"`
}

func createPlacementTest(c *gin.Context) {
	var newPlacementTest Placement_Test
	if err := c.ShouldBindJSON(&newPlacementTest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var interest Interest
	if err := db.First(&interest, newPlacementTest.InterestID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Interest with that ID not found"})
		return
	}

	db.Create(&newPlacementTest)
	c.JSON(http.StatusCreated, newPlacementTest)
}

func getPlacementTests(c *gin.Context) {
	var placementTests []Placement_Test
	db.Preload("Interest").Find(&placementTests)
	c.JSON(http.StatusOK, placementTests)
}

func updatePlacementTest(c *gin.Context) {
	id := c.Param("id")
	var placementTest Placement_Test
	if err := db.Preload("Interest").First(&placementTest, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Placement Test not found"})
		return
	}

	var input struct {
		Question       *string `json:"question"`
		Correct_answer *string `json:"correct_answer"`
		Option_a       *string `json:"option_a"`
		Option_b       *string `json:"option_b"`
		Option_c       *string `json:"option_c"`
		Option_d       *string `json:"option_d"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData := map[string]interface{}{}

	if input.Question != nil {
		updateData["question"] = *input.Question
	}
	if input.Correct_answer != nil {
		updateData["correct_answer"] = *input.Correct_answer
	}
	if input.Option_a != nil {
		updateData["option_a"] = *input.Option_a
	}
	if input.Option_b != nil {
		updateData["option_b"] = *input.Option_b
	}
	if input.Option_c != nil {
		updateData["option_c"] = *input.Option_c
	}
	if input.Option_d != nil {
		updateData["option_d"] = *input.Option_d
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	db.Model(&placementTest).Updates(updateData)
	c.JSON(http.StatusOK, placementTest)
}

// =================================

type Quiz_Answer struct {
	gorm.Model
	QuizanswerID   uint    `gorm:"column:quiz_answer_id;primaryKey;autoIncrement;unique" json:"quiz_answer_id"`
	QuizID         uint    `json:"quiz_id"`
	Quiz           Quiz    `gorm:"references:QuizID"`
	StudentID      uint    `json:"student_id"`
	Student        Student `gorm:"references:StudentID"`
	Student_answer string  `json:"student_answer"`
}

func createQuizAnswer(c *gin.Context) {
	var newQuizAnswer Quiz_Answer
	if err := c.ShouldBindJSON(&newQuizAnswer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var quiz Quiz
	if err := db.First(&quiz, newQuizAnswer.QuizID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Quiz with that ID not found"})
		return
	}

	var student Student
	if err := db.First(&student, newQuizAnswer.StudentID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Student with that ID not found"})
		return
	}

	db.Create(&newQuizAnswer)
	c.JSON(http.StatusCreated, newQuizAnswer)
}

func getQuizAnswers(c *gin.Context) {
	var quizAnswers []Quiz_Answer
	db.Preload("Quiz").Preload("Student").Find(&quizAnswers)
	c.JSON(http.StatusOK, quizAnswers)
}

func getQuizAnswerByStudentID(c *gin.Context) {
	id := c.Param("id")
	var quizAnswer []Quiz_Answer
	if err := db.Preload("Quiz").Preload("Student").Find(&quizAnswer, "student_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz Answer not found"})
		return
	}
	c.JSON(http.StatusOK, quizAnswer)
}

func getQuizAnswerByQuizID(c *gin.Context) {
	id := c.Param("id")
	var quizAnswer []Quiz_Answer
	if err := db.Preload("Quiz").Preload("Student").Find(&quizAnswer, "quiz_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz Answer not found"})
		return
	}
	c.JSON(http.StatusOK, quizAnswer)
}

// =============================

type Quiz_Result struct {
	gorm.Model
	QuizresultID uint      `gorm:"column:quiz_result_id;primaryKey;autoIncrement;unique" json:"quiz_result_id"`
	StudentID    uint      `json:"student_id"`
	Student      Student   `gorm:"references:StudentID"`
	SubjectID    uint      `json:"subject_id"`
	Subject      Subject   `gorm:"references:SubjectID"`
	Score        int       `json:"score"`
	Quiz_date    time.Time `json:"quiz_date"`
}

func createQuizResult(c *gin.Context) {
	var newQuizResult Quiz_Result
	if err := c.ShouldBindJSON(&newQuizResult); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newQuizResult.Quiz_date = time.Now()

	var student Student
	if err := db.First(&student, newQuizResult.StudentID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Student with that ID not found"})
		return
	}

	var subject Subject
	if err := db.First(&subject, newQuizResult.SubjectID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subject with that ID not found"})
		return
	}

	db.Create(&newQuizResult)
	c.JSON(http.StatusCreated, newQuizResult)
}

func getQuizResults(c *gin.Context) {
	var quizResults []Quiz_Result
	db.Preload("Student").Preload("Subject").Find(&quizResults)
	c.JSON(http.StatusOK, quizResults)
}

func getQuizResultByStudentID(c *gin.Context) {
	id := c.Param("id")
	var quizResult []Quiz_Result
	if err := db.Preload("Student").Preload("Subject").Find(&quizResult, "student_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz Result not found"})
		return
	}

	c.JSON(http.StatusOK, quizResult)
}

func getQuizResultBySubjectID(c *gin.Context) {
	id := c.Param("id")
	var quizResult []Quiz_Result
	if err := db.Preload("Student").Preload("Subject").Find(&quizResult, "subject_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz Result not found"})
		return
	}

	c.JSON(http.StatusOK, quizResult)
}

// ============================

type Quiz struct {
	gorm.Model
	QuizID         uint    `gorm:"column:quiz_id;primaryKey;autoIncrement;unique" json:"quiz_id"`
	SubjectID      uint    `json:"subject_id"`
	Subject        Subject `gorm:"references:SubjectID"`
	Question       string  `json:"question"`
	Correct_answer string  `json:"correct_answer"`
	Option_a       string  `json:"option_a"`
	Option_b       string  `json:"option_b"`
	Option_c       string  `json:"option_c"`
	Option_d       string  `json:"option_d"`
}

func createQuiz(c *gin.Context) {
	var newQuiz Quiz
	if err := c.ShouldBindJSON(&newQuiz); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var subject Subject
	if err := db.First(&subject, newQuiz.SubjectID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subject not found"})
		return
	}

	db.Create(&newQuiz)
	c.JSON(http.StatusCreated, newQuiz)
}

func getQuizs(c *gin.Context) {
	var quizs []Quiz
	db.Preload("Subject").Find(&quizs)
	c.JSON(http.StatusOK, quizs)
}

func getQuizBySubjectID(c *gin.Context) {
	id := c.Param("id")
	var quiz []Quiz
	if err := db.Preload("Subject").Find(&quiz, "subject_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}

	c.JSON(http.StatusOK, quiz)
}

func getQuizByID(c *gin.Context) {
	id := c.Param("id")
	var quiz Quiz
	if err := db.Preload("Subject").First(&quiz, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}

	c.JSON(http.StatusOK, quiz)
}

func updateQuiz(c *gin.Context) {
	id := c.Param("id")
	var quiz Quiz
	if err := db.Preload("Subject").First(&quiz, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}

	var input struct {
		Question       *string `json:"question"`
		Correct_answer *string `json:"correct_answer"`
		Option_a       *string `json:"option_a"`
		Option_b       *string `json:"option_b"`
		Option_c       *string `json:"option_c"`
		Option_d       *string `json:"option_d"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData := map[string]interface{}{}

	if input.Question != nil {
		updateData["question"] = *input.Question
	}
	if input.Correct_answer != nil {
		updateData["correct_answer"] = *input.Correct_answer
	}
	if input.Option_a != nil {
		updateData["option_a"] = *input.Option_a
	}
	if input.Option_b != nil {
		updateData["option_b"] = *input.Option_b
	}
	if input.Option_c != nil {
		updateData["option_c"] = *input.Option_c
	}
	if input.Option_d != nil {
		updateData["option_d"] = *input.Option_d
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	db.Model(&quiz).Updates(updateData)
	c.JSON(http.StatusOK, quiz)
}

// ================================

type Student struct {
	gorm.Model
	StudentID    uint     `gorm:"column:student_id;primaryKey;autoIncrement;unique" json:"student_id"`
	Phone_number string   `json:"phone_number"`
	Name         string   `json:"name"`
	Residence    string   `json:"residence"`
	InterestID   uint     `json:"interest_id"`
	Interest     Interest `gorm:"references:InterestID"`
}

func createStudent(c *gin.Context) {
	var newStudent struct {
		Student
		Password *string `json:"password"`
	}

	if err := c.ShouldBindJSON(&newStudent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newStudentData := Student{
		StudentID:    newStudent.StudentID,
		Phone_number: newStudent.Phone_number,
		Name:         newStudent.Name,
		Residence:    newStudent.Residence,
		InterestID:   newStudent.InterestID,
		Interest:     newStudent.Interest,
	}

	var interest Interest
	if err := db.First(&interest, newStudentData.InterestID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Interest with that ID not found"})
		return
	}

	db.Create(&newStudentData)
	controllers.Signup(c, newStudentData.StudentID, *newStudent.Password)

	c.JSON(http.StatusCreated, newStudentData)
}

func getStudents(c *gin.Context) {
	var students []Student
	db.Preload("Interest").Find(&students)
	c.JSON(http.StatusOK, students)
}

func getStudentByID(c *gin.Context) {
	id := c.Param("id")
	var student Student
	if err := db.Preload("Interest").First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, student)
}

func updateStudent(c *gin.Context) {
	id := c.Param("id")
	var student Student
	if err := db.Preload("Interest").First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	var input struct {
		Phone_number *string `json:"phone_number"`
		Residence    *string `json:"residence"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isOnlyNumber(*input.Phone_number) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone Number must be number"})
		return
	}

	if !isOnlyAlphabet(*input.Residence) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Residence must be alphabet"})
		return
	}

	updateData := map[string]interface{}{}

	if input.Phone_number != nil {
		updateData["phone_number"] = *input.Phone_number
	}

	if input.Residence != nil {
		updateData["residence"] = *input.Residence
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	db.Model(&student).Updates(updateData)
	c.JSON(http.StatusOK, student)
}

func deleteStudent(c *gin.Context) {
	id := c.Param("id")
	var student Student
	if err := db.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	db.Delete(&student)
	c.JSON(http.StatusOK, gin.H{"message": "Student deleted"})
}

// ===================================

type Subject struct {
	gorm.Model
	SubjectID      uint     `gorm:"column:subject_id;primaryKey;autoIncrement;unique" json:"subject_id"`
	Subject_name   string   `json:"subject_name"`
	Description    string   `json:"description"`
	InterestID     uint     `json:"interest_id"`
	Interest       Interest `gorm:"references:InterestID"`
	PrerequisiteID uint     `json:"prerequisite_id"`
}

func createSubject(c *gin.Context) {
	var newSubject Subject
	if err := c.ShouldBindJSON(&newSubject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var interest Interest
	if err := db.First(&interest, newSubject.InterestID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Interest with that ID not found"})
		return
	}

	db.Create(&newSubject)
	c.JSON(http.StatusCreated, newSubject)
}

func getSubjects(c *gin.Context) {
	var subjects []Subject
	db.Preload("Interest").Find(&subjects)
	c.JSON(http.StatusOK, subjects)
}

func getSubjectByInterestID(c *gin.Context) {
	id := c.Param("id")
	var subject []Subject
	if err := db.Preload("Interest").Find(&subject, "interest_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subject not found"})
		return
	}
	c.JSON(http.StatusOK, subject)
}

func getSubjectByID(c *gin.Context) {
	id := c.Param("id")
	var subject Subject
	if err := db.Preload("Interest").First(&subject, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subject not found"})
		return
	}
	c.JSON(http.StatusOK, subject)
}

func updateSubject(c *gin.Context) {
	id := c.Param("id")
	var subject Subject
	if err := db.Preload("Interest").First(&subject, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subject not found"})
		return
	}

	var input struct {
		Description     *string `json:"description"`
		Prerequisite_id *uint   `json:"prerequisite_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateData := map[string]interface{}{}

	if input.Description != nil {
		updateData["description"] = *input.Description
	}
	if input.Prerequisite_id != nil {
		updateData["prerequisite_id"] = *input.Prerequisite_id
	}

	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	db.Model(&subject).Updates(updateData)
	c.JSON(http.StatusOK, subject)
}

// ====================================

type Subject_Joined struct {
	gorm.Model
	SubjectjoinedID uint      `gorm:"column:subject_joined_id;primaryKey;autoIncrement;unique" json:"subject_joined_id"`
	StudentID       uint      `json:"student_id"`
	Student         Student   `gorm:"references:StudentID"`
	SubjectID       uint      `json:"subject_id"`
	Subject         Subject   `gorm:"references:SubjectID"`
	Date_joined     time.Time `json:"date_joined"`
}

func createSubjectJoined(c *gin.Context) {
	var newSubjectJoined Subject_Joined
	if err := c.ShouldBindJSON(&newSubjectJoined); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newSubjectJoined.Date_joined = time.Now()

	var student Student
	if err := db.First(&student, newSubjectJoined.StudentID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Student with that ID not found"})
		return
	}

	var subject Subject
	if err := db.First(&subject, newSubjectJoined.SubjectID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Subject with that ID not found"})
		return
	}

	db.Create(&newSubjectJoined)
	c.JSON(http.StatusCreated, newSubjectJoined)
}

func getSubjectJoineds(c *gin.Context) {
	var subjectJoineds []Subject_Joined
	db.Preload("Student").Preload("Subject").Find(&subjectJoineds)
	c.JSON(http.StatusOK, subjectJoineds)
}

func getSubjectJoindedByStudentID(c *gin.Context) {
	id := c.Param("id")
	var subjectJoined []Subject_Joined
	if err := db.Preload("Student").Preload("Subject").Find(&subjectJoined, "student_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subject Joined not found"})
		return
	}

	c.JSON(http.StatusOK, subjectJoined)
}

func getSubjectJoinedBySubjectID(c *gin.Context) {
	id := c.Param("id")
	var subjectJoined []Subject_Joined
	if err := db.Preload("Student").Preload("Subject").Find(&subjectJoined, "subject_id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subject Joined not found"})
		return
	}

	c.JSON(http.StatusOK, subjectJoined)
}

// ====================================

var db *gorm.DB

func Handler(w http.ResponseWriter, r *http.Request) {
	dsn := "postgresql://neondb_owner:npg_m5ljwUdP8fFh@ep-mute-thunder-a4w86qpm-pooler.us-east-1.aws.neon.tech/finpro_api?sslmode=require"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&Interest{})
	db.AutoMigrate(&Student{})
	db.AutoMigrate(&Subject{})
	db.AutoMigrate(&Subject_Joined{})
	db.AutoMigrate(&Placement_Test{})
	db.AutoMigrate(&Placement_Test_Answer{}) // error
	db.AutoMigrate(&Learning_Material{})
	db.AutoMigrate(&Attachment{}) // error
	db.AutoMigrate(&Placement_Test_Result{})
	db.AutoMigrate(&Quiz_Result{})
	db.AutoMigrate(&Quiz{})
	db.AutoMigrate(&Quiz_Answer{})

	router := gin.Default()

	router.POST("/login", controllers.Login)
	router.GET("/validate", middleware.RequireAuth, controllers.Validate)

	router.POST("/interest", createInterest)
	router.GET("/interest", getInterests)
	router.GET("/interest/:id", getInterestByID)
	router.PUT("/interest/:id", updateInterest)

	router.POST("/subject-joined", createSubjectJoined)
	router.GET("/subject-joined", getSubjectJoineds)
	router.GET("/subject-joined/by-student/:id", getSubjectJoindedByStudentID)
	router.GET("/subject-joined/by-subject/:id", getSubjectJoinedBySubjectID)

	router.POST("/student", createStudent)
	router.GET("/student", getStudents)
	router.GET("/student/:id", getStudentByID)
	router.PUT("/student/:id", updateStudent)
	router.DELETE("/student/:id", deleteStudent)

	router.POST("/placement-test", createPlacementTest)
	router.GET("/placement-test", getPlacementTests)
	router.PUT("/placement-test/:id", updatePlacementTest)

	router.POST("/subject", createSubject)
	router.GET("/subject", getSubjects)
	router.GET("/subject/:id", getSubjectByID)
	router.GET("/subject/by-interest/:id", getSubjectByInterestID)
	router.PUT("/subject/:id", updateSubject)

	router.POST("/quiz", createQuiz)
	router.GET("/quiz", getQuizs)
	router.GET("/quiz/by-subject/:id", getQuizBySubjectID)
	router.GET("/quiz/:id", getQuizByID)
	router.PUT("/quiz/:id", updateQuiz)

	router.POST("/learning-material", createLearningMaterial)
	router.GET("/learning-material", getLearningMaterials)
	router.GET("/learning-material/:id", getLearningMaterialByID)
	router.GET("/learning-material/by-subject/:id", getLearningMaterialBySubjectID)
	router.PUT("/learning-material/:id", updateLearningMaterial)
	router.DELETE("/learning-material/:id", deleteLearningMaterial)

	router.POST("/attachment", createAttachment)
	router.GET("/attachment", getAttachments)
	router.GET("/attachment/by-learning-material/:id", getAttachmentByLearningMaterialID)
	router.PUT("/attachment/:id", updateAttachment)

	router.POST("/placement-test-answer", createPlacementTestAnswer)
	router.GET("/placement-test-answer", getPlacementTestAnswers)
	router.GET("/placement-test-answer/by-student/:id", getPlacementTestAnswerByStudentID)
	router.GET("/placement-test-answer/by-placement-test/:id", getPlacementTestAnswerByPlacementTestID)
	router.PUT("/placement-test-answer/:id", updatePlacementTestAnswer)

	router.POST("/quiz-answer", createQuizAnswer)
	router.GET("/quiz-answer", getQuizAnswers)
	router.GET("/quiz-answer/by-student/:id", getQuizAnswerByStudentID)
	router.GET("/quiz-answer/by-quiz/:id", getQuizAnswerByQuizID)

	router.POST("/quiz-result", createQuizResult)
	router.GET("/quiz-result", getQuizResults)
	router.GET("/quiz-result/by-student/:id", getQuizResultByStudentID)
	router.GET("/quiz-result/by-subject/:id", getQuizResultBySubjectID)

	router.POST("/placement-test-result", createPlacementTestResult)
	router.GET("/placement-test-result", getPlacementTestResults)
	router.GET("/placement-test-result/by-student/:id", getPlacementTestResultByStudentID)
	router.GET("/placement-test-result/by-interest/:id", getPlacementTestResultByInterestID)

	router.Run(":8080")
}
