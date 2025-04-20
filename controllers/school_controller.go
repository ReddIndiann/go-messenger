package controllers

import (
	"context"
	"encoding/json"
	"math"
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

func RegisterSchool() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var school model.School
		if err := c.BodyParser(&school); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		// fmt.Printf("Parsed school: %+v\n", school)
		var requestData map[string]interface{}
		if err := json.Unmarshal(c.Body(), &requestData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		if school.Name == "" || school.Email == "" || school.Phone == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "All fields are required",
			})
		}

		collection := database.GetCollection("schools")
		newSchool := model.School{
			ID:        primitive.NewObjectID(),
			Name:      school.Name,
			Email:     school.Email,
			Phone:     school.Phone,
			Logo:      school.Logo,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		result, err := collection.InsertOne(c.Context(), newSchool)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create school",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":   "success",
			"message":  "School created successfully",
			"schoolId": result.InsertedID,
			"info":     newSchool,
		})
	}

}
func GetAllSchool() fiber.Handler {
	return func(c *fiber.Ctx) error {
		collection := database.GetCollection("schools")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Optional: Add pagination
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		skip := (page - 1) * limit

		// Optional: Add filtering
		filter := bson.M{}
		if search := c.Query("search"); search != "" {
			filter["name"] = bson.M{"$regex": search, "$options": "i"}
		}

		// Get total count for pagination
		total, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error counting schools",
			})
		}

		// Find schools with pagination
		opts := options.Find().
			SetSkip(int64(skip)).
			SetLimit(int64(limit)).
			SetSort(bson.D{{Key: "created_at", Value: -1}}) // Sort by newest first

		cursor, err := collection.Find(ctx, filter, opts)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error fetching schools",
			})
		}
		defer cursor.Close(ctx)

		var schools []model.School
		if err = cursor.All(ctx, &schools); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error parsing schools",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data": fiber.Map{
				"schools": schools,
				"pagination": fiber.Map{
					"total": total,
					"page":  page,
					"limit": limit,
					"pages": math.Ceil(float64(total) / float64(limit)),
				},
			},
		})
	}
}

func GetSchoolByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID format",
			})
		}

		collection := database.GetCollection("schools")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var school model.School
		err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&school)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "School not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data":   school,
		})
	}
}

func UpdateSchool() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID format",
			})
		}

		// Get the existing school first
		collection := database.GetCollection("schools")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var existingSchool model.School
		err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingSchool)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "School not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error fetching school",
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
		if fieldExists("name") {
			setMap["name"] = updateData["name"]
		}
		if fieldExists("email") {
			setMap["email"] = updateData["email"]
		}
		if fieldExists("phone") {
			setMap["phone"] = updateData["phone"]
		}
		if fieldExists("logo") {
			setMap["logo"] = updateData["logo"]
		}
		if fieldExists("verified") {
			setMap["verified"] = updateData["verified"]
		}

		// Always update the updated_at field
		setMap["updated_at"] = time.Now()

		// Only perform update if there are fields to update
		if len(setMap) > 1 { // More than just updated_at
			result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Error updating school",
				})
			}

			if result.MatchedCount == 0 {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "School not found",
				})
			}
		}

		// Get the updated school
		var updatedSchool model.School
		err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedSchool)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error fetching updated school",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "School updated successfully",
			"data":    updatedSchool,
		})
	}
}

func DeleteSchool() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
		}

		collection := database.GetCollection("schools")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error deleting school"})
		}

		if result.DeletedCount == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "School not found"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "School deleted successfully"})
	}
}
