package controllers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ReddIndiann/go-messanger/database"
	"github.com/ReddIndiann/go-messanger/model"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterStudent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var student model.Student
		if err := c.BodyParser(&student); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid request",
			})
		}

		// Parse the request body to get all fields
		var requestData map[string]interface{}
		if err := json.Unmarshal(c.Body(), &requestData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid request format",
			})
		}

		// Validate required fields
		if student.FirstName == "" || student.LastName == "" || student.Email == "" ||
			student.Phone == "" || student.SchoolID.IsZero() || student.Grade == "" ||
			student.Section == "" || student.RollNumber == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "First name, last name, email, phone, school ID, grade, section, and roll number are required",
			})
		}

		// Check if school exists
		schoolCollection := database.GetCollection("schools")
		schoolResult := schoolCollection.FindOne(context.Background(), bson.M{"_id": student.SchoolID})
		if schoolResult.Err() == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "School not found with the provided ID",
			})
		}

		// Check if teacher exists if teacher_id is provided
		if !student.TeacherID.IsZero() {
			teacherCollection := database.GetCollection("teachers")
			teacherResult := teacherCollection.FindOne(context.Background(), bson.M{"_id": student.TeacherID})
			if teacherResult.Err() == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "error",
					"message": "Teacher not found with the provided ID",
				})
			}
		}

		// Check if student with same email already exists
		collection := database.GetCollection("students")
		existingStudent := collection.FindOne(context.Background(), bson.M{"email": student.Email})
		if existingStudent.Err() == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Student with this email already exists",
			})
		}

		// Create new student document
		newStudent := model.Student{
			ID:            primitive.NewObjectID(),
			SchoolID:      student.SchoolID,
			TeacherID:     student.TeacherID,
			FirstName:     student.FirstName,
			LastName:      student.LastName,
			Email:         student.Email,
			Phone:         student.Phone,
			DateOfBirth:   student.DateOfBirth,
			Gender:        student.Gender,
			Address:       student.Address,
			Grade:         student.Grade,
			Section:       student.Section,
			RollNumber:    student.RollNumber,
			ParentDetails: student.ParentDetails,
			Status:        "Active", // Default status
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// Insert the new student
		result, err := collection.InsertOne(c.Context(), newStudent)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to create student",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":    "success",
			"message":   "Student registered successfully",
			"studentId": result.InsertedID,
			"info":      newStudent,
		})
	}
}

func GetStudent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		studentID := c.Params("id")
		objID, err := primitive.ObjectIDFromHex(studentID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid student ID",
			})
		}

		collection := database.GetCollection("students")
		var student model.Student
		err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&student)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status":  "error",
					"message": "Student not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch student",
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"student": student,
		})
	}
}

func UpdateStudent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		studentID := c.Params("id")
		objID, err := primitive.ObjectIDFromHex(studentID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid student ID",
			})
		}

		var updateData map[string]interface{}
		if err := json.Unmarshal(c.Body(), &updateData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid request format",
			})
		}

		// Remove any fields that shouldn't be updated
		delete(updateData, "_id")
		delete(updateData, "created_at")
		updateData["updated_at"] = time.Now()

		// If school_id is being updated, validate the new school exists
		if schoolID, ok := updateData["school_id"]; ok {
			schoolObjID, err := primitive.ObjectIDFromHex(schoolID.(string))
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "error",
					"message": "Invalid school ID format",
				})
			}
			schoolCollection := database.GetCollection("schools")
			if err := schoolCollection.FindOne(context.Background(), bson.M{"_id": schoolObjID}).Err(); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "error",
					"message": "School not found",
				})
			}
		}

		// If teacher_id is being updated, validate the new teacher exists
		if teacherID, ok := updateData["teacher_id"]; ok {
			teacherObjID, err := primitive.ObjectIDFromHex(teacherID.(string))
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "error",
					"message": "Invalid teacher ID format",
				})
			}
			teacherCollection := database.GetCollection("teachers")
			if err := teacherCollection.FindOne(context.Background(), bson.M{"_id": teacherObjID}).Err(); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "error",
					"message": "Teacher not found",
				})
			}
		}

		collection := database.GetCollection("students")
		result, err := collection.UpdateOne(
			context.Background(),
			bson.M{"_id": objID},
			bson.M{"$set": updateData},
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to update student",
			})
		}

		if result.MatchedCount == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Student not found",
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Student updated successfully",
		})
	}
}

func DeleteStudent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		studentID := c.Params("id")
		objID, err := primitive.ObjectIDFromHex(studentID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid student ID",
			})
		}

		collection := database.GetCollection("students")
		result, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to delete student",
			})
		}

		if result.DeletedCount == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Student not found",
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Student deleted successfully",
		})
	}
}

func ListStudents() fiber.Handler {
	return func(c *fiber.Ctx) error {
		collection := database.GetCollection("students")
		cursor, err := collection.Find(context.Background(), bson.M{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch students",
			})
		}
		defer cursor.Close(context.Background())

		var students []model.Student
		if err := cursor.All(context.Background(), &students); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to decode students",
			})
		}

		return c.JSON(fiber.Map{
			"status":   "success",
			"students": students,
		})
	}
}
