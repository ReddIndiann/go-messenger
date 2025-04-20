package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SchoolSubject struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name        string               `bson:"name" json:"name"`
	Description string               `bson:"description" json:"description"`
	SchoolID    primitive.ObjectID   `bson:"school_id" json:"school_id"`     // Reference to School
	TeacherID   primitive.ObjectID   `bson:"teacher_id" json:"teacher_id"`   // Reference to Teacher
	StudentIDs  []primitive.ObjectID `bson:"student_ids" json:"student_ids"` // References to Students
	Grade       string               `bson:"grade" json:"grade"`             // e.g., "10th Grade"
	Section     string               `bson:"section" json:"section"`         // e.g., "A", "B"
	Status      string               `bson:"status" json:"status"`           // e.g., "Active", "Inactive"
	CreatedAt   time.Time            `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at,omitempty" json:"updated_at"`
}
