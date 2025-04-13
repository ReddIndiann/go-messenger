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

		var school model.School
		if err := c.BodyParser(&school); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		collection := database.GetCollection("schools")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		update := bson.M{
			"$set": bson.M{
				"name":       school.Name,
				"email":      school.Email,
				"phone":      school.Phone,
				"logo":       school.Logo,
				"updated_at": time.Now(),
			},
		}

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

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "School updated successfully",
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
