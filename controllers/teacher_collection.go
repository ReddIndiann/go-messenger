package controllers

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/ReddIndiann/go-messanger/database"
	"github.com/ReddIndiann/go-messanger/model"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func RegisterTeacher() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var teacher model.Teacher
		if err := c.BodyParser(&teacher); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		// Parse the request body to get all fields
		var requestData map[string]interface{}
		if err := json.Unmarshal(c.Body(), &requestData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		// Validate required fields
		if teacher.FirstName == "" || teacher.LastName == "" || teacher.Email == "" ||
			teacher.Phone == "" || teacher.SchoolID.IsZero() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "First name, last name, email, phone, and school ID are required",
			})
		}

		// Check if school exists
		schoolCollection := database.GetCollection("schools")
		schoolResult := schoolCollection.FindOne(context.Background(), bson.M{"_id": teacher.SchoolID})
		if schoolResult.Err() == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "School not found with the provided ID",
			})
		}

		// Check if teacher with same email already exists
		collection := database.GetCollection("teachers")
		existingTeacher := collection.FindOne(context.Background(), bson.M{"email": teacher.Email})
		if existingTeacher.Err() == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Teacher with this email already exists",
			})
		}

		// Create new teacher document
		newTeacher := model.Teacher{
			ID:               primitive.NewObjectID(),
			SchoolID:         teacher.SchoolID,
			FirstName:        teacher.FirstName,
			LastName:         teacher.LastName,
			Email:            teacher.Email,
			Phone:            teacher.Phone,
			DateOfBirth:      teacher.DateOfBirth,
			Gender:           teacher.Gender,
			Address:          teacher.Address,
			Qualifications:   teacher.Qualifications,
			SubjectIDs:       teacher.SubjectIDs,
			DepartmentID:     teacher.DepartmentID,
			GradeLevels:      teacher.GradeLevels,
			Designation:      teacher.Designation,
			JoiningDate:      time.Now(), // Set joining date to current time
			Experience:       teacher.Experience,
			Salary:           teacher.Salary,
			Status:           "Active", // Default status
			EmergencyContact: teacher.EmergencyContact,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		// Insert the new teacher
		result, err := collection.InsertOne(c.Context(), newTeacher)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create teacher",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":    "success",
			"message":   "Teacher registered successfully",
			"teacherId": result.InsertedID,
			"info":      newTeacher,
		})
	}
}

func GetAllTeachers() fiber.Handler {
	return func(c *fiber.Ctx) error {
		collection := database.GetCollection("teachers")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Optional: Add pagination
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		skip := (page - 1) * limit

		// Optional: Add filtering
		filter := bson.M{}
		if search := c.Query("search"); search != "" {
			filter["$or"] = []bson.M{
				{"first_name": bson.M{"$regex": search, "$options": "i"}},
				{"last_name": bson.M{"$regex": search, "$options": "i"}},
				{"email": bson.M{"$regex": search, "$options": "i"}},
			}
		}

		// Get total count for pagination
		total, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error counting teachers",
			})
		}

		// Find teachers with pagination
		opts := options.Find().
			SetSkip(int64(skip)).
			SetLimit(int64(limit)).
			SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by newest first

		cursor, err := collection.Find(ctx, filter, opts)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error fetching teachers",
			})
		}
		defer cursor.Close(ctx)

		var teachers []model.Teacher
		if err = cursor.All(ctx, &teachers); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error parsing teachers",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"teachers": teachers,
				"pagination": fiber.Map{
					"total": total,
					"page":  page,
					"limit": limit,
					"pages": total/int64(limit) + 1,
				},
			},
		})
	}
}

func GetTeacherByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID format",
			})
		}

		collection := database.GetCollection("teachers")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var teacher model.Teacher
		err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&teacher)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Teacher not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error fetching teacher",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data":   teacher,
		})
	}
}

func UpdateTeacher() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID format",
			})
		}

		// Get the existing teacher first
		collection := database.GetCollection("teachers")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var existingTeacher model.Teacher
		err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingTeacher)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Teacher not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error fetching teacher",
			})
		}

		// Parse the update request
		var updateData map[string]interface{}
		if err := json.Unmarshal(c.Body(), &updateData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Prepare update fields, only including fields that are provided
		update := bson.M{"$set": bson.M{}}
		setMap := update["$set"].(bson.M)

		// Helper function to check if a field exists in the update data
		fieldExists := func(field string) bool {
			_, exists := updateData[field]
			return exists
		}

		// Update only the fields that are provided
		if fieldExists("first_name") {
			setMap["first_name"] = updateData["first_name"]
		}
		if fieldExists("last_name") {
			setMap["last_name"] = updateData["last_name"]
		}
		if fieldExists("email") {
			setMap["email"] = updateData["email"]
		}
		if fieldExists("phone") {
			setMap["phone"] = updateData["phone"]
		}
		if fieldExists("date_of_birth") {
			setMap["date_of_birth"] = updateData["date_of_birth"]
		}
		if fieldExists("gender") {
			setMap["gender"] = updateData["gender"]
		}
		if fieldExists("address") {
			setMap["address"] = updateData["address"]
		}
		if fieldExists("qualifications") {
			setMap["qualifications"] = updateData["qualifications"]
		}
		if fieldExists("subject_ids") {
			setMap["subject_ids"] = updateData["subject_ids"]
		}
		if fieldExists("department_id") {
			setMap["department_id"] = updateData["department_id"]
		}
		if fieldExists("grade_levels") {
			setMap["grade_levels"] = updateData["grade_levels"]
		}
		if fieldExists("designation") {
			setMap["designation"] = updateData["designation"]
		}
		if fieldExists("experience") {
			setMap["experience"] = updateData["experience"]
		}
		if fieldExists("salary") {
			setMap["salary"] = updateData["salary"]
		}
		if fieldExists("status") {
			setMap["status"] = updateData["status"]
		}
		if fieldExists("emergency_contact") {
			setMap["emergency_contact"] = updateData["emergency_contact"]
		}

		// Always update the updated_at field
		setMap["updated_at"] = time.Now()

		// Only perform update if there are fields to update
		if len(setMap) > 1 { // More than just updated_at
			result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Error updating teacher",
				})
			}

			if result.MatchedCount == 0 {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Teacher not found",
				})
			}
		}

		// Get the updated teacher
		var updatedTeacher model.Teacher
		err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedTeacher)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error fetching updated teacher",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Teacher updated successfully",
			"data":    updatedTeacher,
		})
	}
}

func DeleteTeacher() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID format",
			})
		}

		collection := database.GetCollection("teachers")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error deleting teacher",
			})
		}

		if result.DeletedCount == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Teacher not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "Teacher deleted successfully",
		})
	}
}
