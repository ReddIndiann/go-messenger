package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SchoolID      primitive.ObjectID `bson:"school_id" json:"school_id"`   // Reference to School
	TeacherID     primitive.ObjectID `bson:"teacher_id" json:"teacher_id"` // Reference to Teacher
	FirstName     string             `bson:"first_name" json:"first_name"`
	LastName      string             `bson:"last_name" json:"last_name"`
	Email         string             `bson:"email" json:"email"`
	Phone         string             `bson:"phone" json:"phone"`
	DateOfBirth   time.Time          `bson:"date_of_birth" json:"date_of_birth"`
	Gender        string             `bson:"gender" json:"gender"`
	Address       Address            `bson:"address" json:"address"`
	Grade         string             `bson:"grade" json:"grade"`     // e.g., "10th Grade"
	Section       string             `bson:"section" json:"section"` // e.g., "A", "B"
	RollNumber    string             `bson:"roll_number" json:"roll_number"`
	ParentDetails ParentDetails      `bson:"parent_details" json:"parent_details"`
	Status        string             `bson:"status" json:"status"` // Active, Inactive, Graduated, etc.
	CreatedAt     time.Time          `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at,omitempty" json:"updated_at"`
}

type Address struct {
	Street     string `bson:"street" json:"street"`
	City       string `bson:"city" json:"city"`
	State      string `bson:"state" json:"state"`
	Country    string `bson:"country" json:"country"`
	PostalCode string `bson:"postal_code" json:"postal_code"`
}

type ParentDetails struct {
	FatherName    string `bson:"father_name" json:"father_name"`
	FatherPhone   string `bson:"father_phone" json:"father_phone"`
	FatherEmail   string `bson:"father_email" json:"father_email"`
	MotherName    string `bson:"mother_name" json:"mother_name"`
	MotherPhone   string `bson:"mother_phone" json:"mother_phone"`
	MotherEmail   string `bson:"mother_email" json:"mother_email"`
	GuardianName  string `bson:"guardian_name" json:"guardian_name"` // If different from parents
	GuardianPhone string `bson:"guardian_phone" json:"guardian_phone"`
	GuardianEmail string `bson:"guardian_email" json:"guardian_email"`
}
