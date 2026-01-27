package main

import "time"

// User represents a staff member
type User struct {
	ID        int
	Username  string
	Password  string
	FullName  string
	Email     string
	CreatedAt time.Time
}

// Objective represents a performance objective
type Objective struct {
	ID          int
	UserID      int
	Title       string
	Description string
	StartDate   time.Time
	EndDate     time.Time
	CreatedAt   time.Time
	Performance float64 // Calculated mean percentage
}

// ExpectedOutcome represents an expected outcome for an objective
type ExpectedOutcome struct {
	ID          int
	ObjectiveID int
	Title       string
	Description string
	CreatedAt   time.Time
}

// ActivityCategory represents the frequency category of an activity
type ActivityCategory string

const (
	CategoryDaily      ActivityCategory = "Daily"
	CategoryWeekly     ActivityCategory = "Weekly"
	CategoryMonthly    ActivityCategory = "Monthly"
	CategoryQuarterly  ActivityCategory = "Quarterly"
	CategoryBiannually ActivityCategory = "Biannually"
	CategoryAnnually   ActivityCategory = "Annually"
)

// TaskPriority represents task priority levels
type TaskPriority string

const (
	PriorityLow    TaskPriority = "Low"
	PriorityMedium TaskPriority = "Medium"
	PriorityHigh   TaskPriority = "High"
	PriorityUrgent TaskPriority = "Urgent"
)

// TaskStatus represents task completion status
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "Pending"
	TaskStatusInProgress TaskStatus = "In Progress"
	TaskStatusCompleted  TaskStatus = "Completed"
	TaskStatusOnHold     TaskStatus = "On Hold"
)

// Task represents a standalone task (separate from activities)
type Task struct {
	ID          int
	UserID      int
	Title       string
	Description string
	Priority    TaskPriority
	Status      TaskStatus
	DueDate     time.Time
	CreatedAt   time.Time
	CompletedAt *time.Time
}

// Activity represents a task or activity
type Activity struct {
	ID                  int
	ExpectedOutcomeID   int
	Title               string
	Description         string
	Category            ActivityCategory
	ProgressPercentage  float64
	ImplementationLevel string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// ViewModels for templates

type DashboardData struct {
	User       User
	Objectives []ObjectiveWithOutcomes
}

type ObjectiveWithOutcomes struct {
	Objective        Objective
	ExpectedOutcomes []ExpectedOutcomeWithActivities
}

type ExpectedOutcomeWithActivities struct {
	ExpectedOutcome ExpectedOutcome
	Activities      []Activity
}

type ObjectiveFormData struct {
	User      User
	Objective *Objective
	IsEdit    bool
}

type ExpectedOutcomeFormData struct {
	User            User
	Objective       Objective
	ExpectedOutcome *ExpectedOutcome
	IsEdit          bool
}

type ActivityFormData struct {
	User            User
	Objective       Objective
	ExpectedOutcome ExpectedOutcome
	Activity        *Activity
	Categories      []ActivityCategory
	IsEdit          bool
}

type TaskListData struct {
	User       User
	Tasks      []Task
	Priorities []TaskPriority
	Statuses   []TaskStatus
}

type TaskFormData struct {
	User       User
	Task       *Task
	Priorities []TaskPriority
	Statuses   []TaskStatus
	IsEdit     bool
}

type ReportData struct {
	User               User
	Objectives         []ObjectiveWithOutcomes
	Tasks              []Task
	TotalObjectives    int
	CompletedTasks     int
	PendingTasks       int
	AveragePerformance float64
}
