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

func RegisterSubject() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var subject model.SchoolSubject
		if err := c.BodyParser(&subject); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Invalid request format",
				"details": err.Error(),
			})
		}

		// Validate required fields
		if subject.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Subject name is required",
			})
		}
		if subject.SchoolID.IsZero() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "School ID is required",
			})
		}
		if subject.TeacherID.IsZero() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Teacher ID is required",
			})
		}

		// Validate optional fields if provided
		if subject.Grade != "" && (len(subject.Grade) > 2 || !isValidGrade(subject.Grade)) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid grade format. Grade should be a number between 1-12",
			})
		}
		if subject.Section != "" && (len(subject.Section) != 1 || !isValidSection(subject.Section)) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid section format. Section should be a single letter (A-Z)",
			})
		}
		if subject.Status != "" && subject.Status != "Active" && subject.Status != "Inactive" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid status. Status must be either 'Active' or 'Inactive'",
			})
		}

		// Check if school exists
		schoolCollection := database.GetCollection("schools")
		if err := schoolCollection.FindOne(context.Background(), bson.M{"_id": subject.SchoolID}).Err(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "School not found with the provided ID",
			})
		}

		// Check if teacher exists
		teacherCollection := database.GetCollection("teachers")
		if err := teacherCollection.FindOne(context.Background(), bson.M{"_id": subject.TeacherID}).Err(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Teacher not found with the provided ID",
			})
		}

		// Check if students exist
		if len(subject.StudentIDs) > 0 {
			studentCollection := database.GetCollection("students")
			for _, studentID := range subject.StudentIDs {
				if err := studentCollection.FindOne(context.Background(), bson.M{"_id": studentID}).Err(); err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "Student not found with ID: " + studentID.Hex(),
					})
				}
			}
		}

		// Create new subject
		newSubject := model.SchoolSubject{
			ID:          primitive.NewObjectID(),
			Name:        subject.Name,
			Description: subject.Description,
			SchoolID:    subject.SchoolID,
			TeacherID:   subject.TeacherID,
			StudentIDs:  subject.StudentIDs,
			Grade:       subject.Grade,
			Section:     subject.Section,
			Status:      "Active", // Default status
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Insert the new subject
		collection := database.GetCollection("subjects")
		result, err := collection.InsertOne(c.Context(), newSubject)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to create subject",
				"details": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":    "success",
			"message":   "Subject created successfully",
			"subjectId": result.InsertedID,
			"info":      newSubject,
		})
	}
}

// Helper function to validate grade
func isValidGrade(grade string) bool {
	if len(grade) > 2 {
		return false
	}
	// Check if grade is a number between 1-12
	num := 0
	for _, c := range grade {
		if c < '0' || c > '9' {
			return false
		}
		num = num*10 + int(c-'0')
	}
	return num >= 1 && num <= 12
}

// Helper function to validate section
func isValidSection(section string) bool {
	if len(section) != 1 {
		return false
	}
	c := section[0]
	return c >= 'A' && c <= 'Z'
}

func GetSubject() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID format",
			})
		}

		collection := database.GetCollection("subjects")
		var subject model.SchoolSubject
		err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&subject)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Subject not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error fetching subject",
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"subject": subject,
		})
	}
}

func UpdateSubject() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID format",
			})
		}

		var updateData map[string]interface{}
		if err := json.Unmarshal(c.Body(), &updateData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		// Remove any fields that shouldn't be updated
		delete(updateData, "_id")
		delete(updateData, "created_at")
		updateData["updated_at"] = time.Now()

		collection := database.GetCollection("subjects")
		result, err := collection.UpdateOne(
			context.Background(),
			bson.M{"_id": objectID},
			bson.M{"$set": updateData},
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update subject",
			})
		}

		if result.MatchedCount == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Subject not found",
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Subject updated successfully",
		})
	}
}

func DeleteSubject() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID format",
			})
		}

		collection := database.GetCollection("subjects")
		result, err := collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete subject",
			})
		}

		if result.DeletedCount == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Subject not found",
			})
		}

		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Subject deleted successfully",
		})
	}
}

func ListSubjects() fiber.Handler {
	return func(c *fiber.Ctx) error {
		collection := database.GetCollection("subjects")
		cursor, err := collection.Find(context.Background(), bson.M{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch subjects",
			})
		}
		defer cursor.Close(context.Background())

		var subjects []model.SchoolSubject
		if err := cursor.All(context.Background(), &subjects); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to decode subjects",
			})
		}

		return c.JSON(fiber.Map{
			"status":   "success",
			"subjects": subjects,
		})
	}
}
