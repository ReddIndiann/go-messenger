package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Teacher struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SchoolID         primitive.ObjectID `bson:"school_id" json:"school_id"` // Reference to School
	FirstName        string             `bson:"first_name" json:"first_name"`
	LastName         string             `bson:"last_name" json:"last_name"`
	Email            string             `bson:"email" json:"email"`
	Phone            string             `bson:"phone" json:"phone"`
	DateOfBirth      time.Time          `bson:"date_of_birth" json:"date_of_birth"`
	Gender           string             `bson:"gender" json:"gender"`
	Address          Address            `bson:"address" json:"address"`
	Qualifications   []Qualification    `bson:"qualifications" json:"qualifications"`
	SubjectIDs       primitive.ObjectID `bson:"subject_ids" json:"subject_id"`      // References to Subject collection
	DepartmentID     primitive.ObjectID `bson:"department_id" json:"department_id"` // Reference to Department collection
	GradeLevels      []string           `bson:"grade_levels" json:"grade_levels"`   // Grade levels they teach
	Designation      string             `bson:"designation" json:"designation"`     // e.g., "Senior Teacher", "Head of Department"
	JoiningDate      time.Time          `bson:"joining_date" json:"joining_date"`
	Experience       int                `bson:"experience" json:"experience"` // Years of experience
	Salary           float64            `bson:"salary" json:"salary"`
	Status           string             `bson:"status" json:"status"` // Active, On Leave, Resigned, etc.
	EmergencyContact EmergencyContact   `bson:"emergency_contact" json:"emergency_contact"`
	CreatedAt        time.Time          `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at,omitempty" json:"updated_at"`
}

// New Subject model
type Subject struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SchoolID    primitive.ObjectID `bson:"school_id" json:"school_id"` // Reference to School
	Name        string             `bson:"name" json:"name"`           // e.g., "Mathematics", "Physics"
	Code        string             `bson:"code" json:"code"`           // e.g., "MATH101", "PHY101"
	Description string             `bson:"description" json:"description"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty" json:"updated_at"`
}

// New Department model
type Department struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SchoolID    primitive.ObjectID `bson:"school_id" json:"school_id"` // Reference to School
	Name        string             `bson:"name" json:"name"`           // e.g., "Mathematics", "Science"
	Code        string             `bson:"code" json:"code"`           // e.g., "MATH", "SCI"
	Description string             `bson:"description" json:"description"`
	HeadID      primitive.ObjectID `bson:"head_id" json:"head_id"` // Reference to Teacher (Department Head)
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty" json:"updated_at"`
}

type Qualification struct {
	Degree         string `bson:"degree" json:"degree"` // e.g., "B.Ed", "M.Ed"
	Institution    string `bson:"institution" json:"institution"`
	YearCompleted  int    `bson:"year_completed" json:"year_completed"`
	Specialization string `bson:"specialization" json:"specialization"` // e.g., "Mathematics", "Physics"
}

type EmergencyContact struct {
	Name         string `bson:"name" json:"name"`
	Relationship string `bson:"relationship" json:"relationship"`
	Phone        string `bson:"phone" json:"phone"`
	Email        string `bson:"email" json:"email"`
	Address      string `bson:"address" json:"address"`
}
