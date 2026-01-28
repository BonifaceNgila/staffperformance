package main

import "time"

// UserRole represents user roles in the system
type UserRole string

const (
	RoleAdmin      UserRole = "Admin"
	RoleSupervisor UserRole = "Supervisor"
	RoleStaff      UserRole = "Staff"
)

// User represents a staff member
type User struct {
	ID                 int
	Username           string
	Password           string
	FullName           string
	Email              string
	Role               UserRole
	SupervisorID       *int
	Department         string
	Position           string
	CreatedAt          time.Time
	SupervisorName     string  // For display purposes
	OverallPerformance float64 // For supervisor dashboard
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

// TaskType represents the type of task
type TaskType string

const (
	TaskTypePersonal        TaskType = "Personal"
	TaskTypeServiceRequest  TaskType = "Service Request"
	TaskTypeStaffAssignment TaskType = "Staff Assignment"
	TaskTypeResponse        TaskType = "Response"
)

// Task represents a standalone task linked to an expected outcome
type Task struct {
	ID                   int
	ExpectedOutcomeID    *int // Links task to an expected outcome
	UserID               int
	AssignedToID         *int
	Title                string
	Description          string
	Priority             TaskPriority
	Status               TaskStatus
	TaskType             TaskType
	RequestedBy          string
	DueDate              time.Time
	CreatedAt            time.Time
	CompletedAt          *time.Time
	CompletionPercentage float64 // 0-100, indicates how much of the task is completed
	AssignedToUser       *User
	// For display purposes
	ObjectiveTitle       string
	ExpectedOutcomeTitle string
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
	Tasks           []Task // Tasks linked to this expected outcome
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

// Comment represents a comment on an objective or activity
type Comment struct {
	ID          int
	ObjectiveID *int
	ActivityID  *int
	UserID      int
	CommentText string
	CreatedAt   time.Time
	Username    string // For displaying commenter name
	UserRole    string // For displaying commenter role
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
	Objectives []ObjectiveWithOutcomes // For selecting expected outcome
	Priorities []TaskPriority
	Statuses   []TaskStatus
	TaskTypes  []TaskType
	IsEdit     bool
}

// TaskWithContext includes task details with associated objective and expected outcome
type TaskWithContext struct {
	Task            Task
	Objective       *Objective
	ExpectedOutcome *ExpectedOutcome
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
type StaffListData struct {
	Username string
	Staff    []User
}

type StaffFormData struct {
	Username    string
	Staff       *User
	Supervisors []User
	IsEdit      bool
}

type SupervisorDashboardData struct {
	Username string
	Role     string
	Staff    []User
}

type ObjectiveWithComments struct {
	Objective  Objective
	Outcomes   []ExpectedOutcome
	Activities []Activity
	Comments   []Comment
}

type CommentFormData struct {
	Username  string
	Objective *Objective
	Activity  *Activity
	StaffID   string
}
