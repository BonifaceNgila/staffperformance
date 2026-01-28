package main

import (
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

// InitDB initializes the SQLite database
func InitDB() error {
	var err error
	db, err = sql.Open("sqlite", "./staffperformance.db")
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	if err = createTables(); err != nil {
		return err
	}

	// Run migrations
	if err = runMigrations(); err != nil {
		log.Println("Migration warning:", err)
	}

	// Create default admin user if not exists
	if err = createDefaultUser(); err != nil {
		log.Println("Default user already exists or error:", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

func createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		full_name TEXT NOT NULL,
		email TEXT,
		role TEXT NOT NULL DEFAULT 'Staff',
		supervisor_id INTEGER,
		department TEXT,
		position TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (supervisor_id) REFERENCES users(id) ON DELETE SET NULL
	);

	CREATE TABLE IF NOT EXISTS objectives (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		start_date DATETIME,
		end_date DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS expected_outcomes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		objective_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (objective_id) REFERENCES objectives(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS activities (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		expected_outcome_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		category TEXT NOT NULL,
		progress_percentage REAL DEFAULT 0,
		implementation_level TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (expected_outcome_id) REFERENCES expected_outcomes(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		expected_outcome_id INTEGER,
		user_id INTEGER NOT NULL,
		assigned_to_id INTEGER,
		title TEXT NOT NULL,
		description TEXT,
		priority TEXT NOT NULL,
		status TEXT NOT NULL,
		task_type TEXT NOT NULL DEFAULT 'Personal',
		requested_by TEXT,
		due_date DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		completed_at DATETIME,
		completion_percentage REAL DEFAULT 0,
		FOREIGN KEY (expected_outcome_id) REFERENCES expected_outcomes(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (assigned_to_id) REFERENCES users(id) ON DELETE SET NULL
	);
	`

	_, err := db.Exec(schema)
	return err
}

func createDefaultUser() error {
	query := `INSERT INTO users (username, password, full_name, email, role, department, position) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, "admin", "admin123", "Administrator", "admin@example.com", "Admin", "Management", "System Administrator")
	return err
}

// runMigrations handles database schema migrations
func runMigrations() error {
	// Migration: Add expected_outcome_id and completion_percentage to tasks table if they don't exist
	var columnExists int
	err := db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('tasks') WHERE name='expected_outcome_id'`).Scan(&columnExists)
	if err != nil {
		return err
	}

	if columnExists == 0 {
		// Add expected_outcome_id column
		_, err = db.Exec(`ALTER TABLE tasks ADD COLUMN expected_outcome_id INTEGER REFERENCES expected_outcomes(id) ON DELETE CASCADE`)
		if err != nil {
			return err
		}
		log.Println("Migration: Added expected_outcome_id column to tasks table")
	}

	// Check for completion_percentage column
	err = db.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('tasks') WHERE name='completion_percentage'`).Scan(&columnExists)
	if err != nil {
		return err
	}

	if columnExists == 0 {
		// Add completion_percentage column
		_, err = db.Exec(`ALTER TABLE tasks ADD COLUMN completion_percentage REAL DEFAULT 0`)
		if err != nil {
			return err
		}
		log.Println("Migration: Added completion_percentage column to tasks table")
	}

	return nil
}

// User CRUD operations
func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	var supervisorID sql.NullInt64
	query := `SELECT id, username, password, full_name, email, role, supervisor_id, department, position, created_at FROM users WHERE username = ?`
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.FullName, &user.Email, &user.Role, &supervisorID, &user.Department, &user.Position, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	if supervisorID.Valid {
		supervisorIDInt := int(supervisorID.Int64)
		user.SupervisorID = &supervisorIDInt
	}
	return user, nil
}

func GetUserByID(id int) (*User, error) {
	user := &User{}
	var supervisorID sql.NullInt64
	query := `SELECT id, username, password, full_name, email, role, supervisor_id, department, position, created_at FROM users WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password, &user.FullName, &user.Email, &user.Role, &supervisorID, &user.Department, &user.Position, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	if supervisorID.Valid {
		supervisorIDInt := int(supervisorID.Int64)
		user.SupervisorID = &supervisorIDInt
	}
	return user, nil
}

// Staff management functions
func CreateUser(user *User) error {
	query := `INSERT INTO users (username, password, full_name, email, role, supervisor_id, department, position) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, user.Username, user.Password, user.FullName, user.Email, user.Role, user.SupervisorID, user.Department, user.Position)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)
	return nil
}

func UpdateUser(user *User) error {
	// If password is empty, don't update it
	if user.Password == "" {
		query := `UPDATE users SET username = ?, role = ?, supervisor_id = ?, department = ?, position = ? WHERE id = ?`
		_, err := db.Exec(query, user.Username, user.Role, user.SupervisorID, user.Department, user.Position, user.ID)
		return err
	}

	// Update with password
	query := `UPDATE users SET username = ?, password = ?, role = ?, supervisor_id = ?, department = ?, position = ? WHERE id = ?`
	_, err := db.Exec(query, user.Username, user.Password, user.Role, user.SupervisorID, user.Department, user.Position, user.ID)
	return err
}

func DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

func GetAllUsers() ([]User, error) {
	query := `SELECT id, username, password, full_name, email, role, supervisor_id, department, position, created_at FROM users ORDER BY full_name ASC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var supervisorID sql.NullInt64
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.FullName, &user.Email, &user.Role, &supervisorID, &user.Department, &user.Position, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		if supervisorID.Valid {
			supervisorIDInt := int(supervisorID.Int64)
			user.SupervisorID = &supervisorIDInt
		}
		users = append(users, user)
	}
	return users, nil
}

func GetUsersByRole(role UserRole) ([]User, error) {
	query := `SELECT id, username, password, full_name, email, role, supervisor_id, department, position, created_at FROM users WHERE role = ? ORDER BY full_name ASC`
	rows, err := db.Query(query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var supervisorID sql.NullInt64
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.FullName, &user.Email, &user.Role, &supervisorID, &user.Department, &user.Position, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		if supervisorID.Valid {
			supervisorIDInt := int(supervisorID.Int64)
			user.SupervisorID = &supervisorIDInt
		}
		users = append(users, user)
	}
	return users, nil
}

func GetStaffBySupervisor(supervisorID int) ([]User, error) {
	query := `SELECT id, username, password, full_name, email, role, supervisor_id, department, position, created_at FROM users WHERE supervisor_id = ? ORDER BY full_name ASC`
	rows, err := db.Query(query, supervisorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var supID sql.NullInt64
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.FullName, &user.Email, &user.Role, &supID, &user.Department, &user.Position, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		if supID.Valid {
			supervisorIDInt := int(supID.Int64)
			user.SupervisorID = &supervisorIDInt
		}
		users = append(users, user)
	}
	return users, nil
}

// Objective CRUD operations
func CreateObjective(obj *Objective) error {
	query := `INSERT INTO objectives (user_id, title, description, start_date, end_date) VALUES (?, ?, ?, ?, ?)`
	result, err := db.Exec(query, obj.UserID, obj.Title, obj.Description, obj.StartDate, obj.EndDate)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	obj.ID = int(id)
	return nil
}

func GetObjectivesByUserID(userID int) ([]Objective, error) {
	query := `SELECT id, user_id, title, description, start_date, end_date, created_at FROM objectives WHERE user_id = ? ORDER BY created_at DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var objectives []Objective
	for rows.Next() {
		var obj Objective
		err := rows.Scan(&obj.ID, &obj.UserID, &obj.Title, &obj.Description, &obj.StartDate, &obj.EndDate, &obj.CreatedAt)
		if err != nil {
			return nil, err
		}
		// Calculate performance
		obj.Performance = CalculateObjectivePerformance(obj.ID)
		objectives = append(objectives, obj)
	}
	return objectives, nil
}

func GetObjectiveByID(id int) (*Objective, error) {
	obj := &Objective{}
	query := `SELECT id, user_id, title, description, start_date, end_date, created_at FROM objectives WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&obj.ID, &obj.UserID, &obj.Title, &obj.Description, &obj.StartDate, &obj.EndDate, &obj.CreatedAt)
	if err != nil {
		return nil, err
	}
	obj.Performance = CalculateObjectivePerformance(obj.ID)
	return obj, nil
}

func UpdateObjective(obj *Objective) error {
	query := `UPDATE objectives SET title = ?, description = ?, start_date = ?, end_date = ? WHERE id = ?`
	_, err := db.Exec(query, obj.Title, obj.Description, obj.StartDate, obj.EndDate, obj.ID)
	return err
}

func DeleteObjective(id int) error {
	query := `DELETE FROM objectives WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// ExpectedOutcome CRUD operations
func CreateExpectedOutcome(outcome *ExpectedOutcome) error {
	query := `INSERT INTO expected_outcomes (objective_id, title, description) VALUES (?, ?, ?)`
	result, err := db.Exec(query, outcome.ObjectiveID, outcome.Title, outcome.Description)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	outcome.ID = int(id)
	return nil
}

func GetExpectedOutcomesByObjectiveID(objectiveID int) ([]ExpectedOutcome, error) {
	query := `SELECT id, objective_id, title, description, created_at FROM expected_outcomes WHERE objective_id = ? ORDER BY created_at ASC`
	rows, err := db.Query(query, objectiveID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var outcomes []ExpectedOutcome
	for rows.Next() {
		var outcome ExpectedOutcome
		err := rows.Scan(&outcome.ID, &outcome.ObjectiveID, &outcome.Title, &outcome.Description, &outcome.CreatedAt)
		if err != nil {
			return nil, err
		}
		outcomes = append(outcomes, outcome)
	}
	return outcomes, nil
}

func GetExpectedOutcomeByID(id int) (*ExpectedOutcome, error) {
	outcome := &ExpectedOutcome{}
	query := `SELECT id, objective_id, title, description, created_at FROM expected_outcomes WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&outcome.ID, &outcome.ObjectiveID, &outcome.Title, &outcome.Description, &outcome.CreatedAt)
	if err != nil {
		return nil, err
	}
	return outcome, nil
}

func UpdateExpectedOutcome(outcome *ExpectedOutcome) error {
	query := `UPDATE expected_outcomes SET title = ?, description = ? WHERE id = ?`
	_, err := db.Exec(query, outcome.Title, outcome.Description, outcome.ID)
	return err
}

func DeleteExpectedOutcome(id int) error {
	query := `DELETE FROM expected_outcomes WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// Activity CRUD operations
func CreateActivity(activity *Activity) error {
	query := `INSERT INTO activities (expected_outcome_id, title, description, category, progress_percentage, implementation_level) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, activity.ExpectedOutcomeID, activity.Title, activity.Description, activity.Category, activity.ProgressPercentage, activity.ImplementationLevel)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	activity.ID = int(id)
	return nil
}

func GetActivitiesByExpectedOutcomeID(outcomeID int) ([]Activity, error) {
	query := `SELECT id, expected_outcome_id, title, description, category, progress_percentage, implementation_level, created_at, updated_at FROM activities WHERE expected_outcome_id = ? ORDER BY created_at ASC`
	rows, err := db.Query(query, outcomeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []Activity
	for rows.Next() {
		var activity Activity
		err := rows.Scan(&activity.ID, &activity.ExpectedOutcomeID, &activity.Title, &activity.Description, &activity.Category, &activity.ProgressPercentage, &activity.ImplementationLevel, &activity.CreatedAt, &activity.UpdatedAt)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

func GetActivityByID(id int) (*Activity, error) {
	activity := &Activity{}
	query := `SELECT id, expected_outcome_id, title, description, category, progress_percentage, implementation_level, created_at, updated_at FROM activities WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&activity.ID, &activity.ExpectedOutcomeID, &activity.Title, &activity.Description, &activity.Category, &activity.ProgressPercentage, &activity.ImplementationLevel, &activity.CreatedAt, &activity.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return activity, nil
}

func UpdateActivity(activity *Activity) error {
	query := `UPDATE activities SET title = ?, description = ?, category = ?, progress_percentage = ?, implementation_level = ?, updated_at = ? WHERE id = ?`
	_, err := db.Exec(query, activity.Title, activity.Description, activity.Category, activity.ProgressPercentage, activity.ImplementationLevel, time.Now(), activity.ID)
	return err
}

func DeleteActivity(id int) error {
	query := `DELETE FROM activities WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// Calculate objective performance based on mean of all activities
func CalculateObjectivePerformance(objectiveID int) float64 {
	query := `
		SELECT AVG(a.progress_percentage)
		FROM activities a
		INNER JOIN expected_outcomes eo ON a.expected_outcome_id = eo.id
		WHERE eo.objective_id = ?
	`
	var performance sql.NullFloat64
	err := db.QueryRow(query, objectiveID).Scan(&performance)
	if err != nil || !performance.Valid {
		return 0.0
	}
	return performance.Float64
}

// Get activities by objective ID (for performance calculation)
func GetActivitiesByObjectiveID(objectiveID int) ([]Activity, error) {
	query := `
		SELECT a.id, a.expected_outcome_id, a.title, a.description, a.category, a.progress_percentage, a.implementation_level, a.created_at, a.updated_at
		FROM activities a
		INNER JOIN expected_outcomes eo ON a.expected_outcome_id = eo.id
		WHERE eo.objective_id = ?
		ORDER BY a.created_at ASC
	`
	rows, err := db.Query(query, objectiveID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []Activity
	for rows.Next() {
		var activity Activity
		err := rows.Scan(&activity.ID, &activity.ExpectedOutcomeID, &activity.Title, &activity.Description, &activity.Category, &activity.ProgressPercentage, &activity.ImplementationLevel, &activity.CreatedAt, &activity.UpdatedAt)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, nil
}

// Task CRUD operations
func CreateTask(task *Task) error {
	query := `INSERT INTO tasks (expected_outcome_id, user_id, title, description, priority, status, due_date, assigned_to_id, task_type, requested_by, completion_percentage, completed_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, task.ExpectedOutcomeID, task.UserID, task.Title, task.Description, task.Priority, task.Status, task.DueDate, task.AssignedToID, task.TaskType, task.RequestedBy, task.CompletionPercentage, task.CompletedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = int(id)
	return nil
}

func GetTasksByUserID(userID int) ([]Task, error) {
	query := `SELECT id, expected_outcome_id, user_id, title, description, priority, status, due_date, created_at, completed_at, assigned_to_id, task_type, requested_by, completion_percentage FROM tasks WHERE user_id = ? ORDER BY due_date ASC, created_at DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var completedAt sql.NullTime
		var assignedToID sql.NullInt64
		var expectedOutcomeID sql.NullInt64
		err := rows.Scan(&task.ID, &expectedOutcomeID, &task.UserID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.DueDate, &task.CreatedAt, &completedAt, &assignedToID, &task.TaskType, &task.RequestedBy, &task.CompletionPercentage)
		if err != nil {
			return nil, err
		}
		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}
		if assignedToID.Valid {
			assignedID := int(assignedToID.Int64)
			task.AssignedToID = &assignedID
		}
		if expectedOutcomeID.Valid {
			outcomeID := int(expectedOutcomeID.Int64)
			task.ExpectedOutcomeID = &outcomeID
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func GetTaskByID(id int) (*Task, error) {
	task := &Task{}
	var completedAt sql.NullTime
	var assignedToID sql.NullInt64
	var expectedOutcomeID sql.NullInt64
	query := `SELECT id, expected_outcome_id, user_id, title, description, priority, status, due_date, created_at, completed_at, assigned_to_id, task_type, requested_by, completion_percentage FROM tasks WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&task.ID, &expectedOutcomeID, &task.UserID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.DueDate, &task.CreatedAt, &completedAt, &assignedToID, &task.TaskType, &task.RequestedBy, &task.CompletionPercentage)
	if err != nil {
		return nil, err
	}
	if completedAt.Valid {
		task.CompletedAt = &completedAt.Time
	}
	if assignedToID.Valid {
		assignedID := int(assignedToID.Int64)
		task.AssignedToID = &assignedID
	}
	if expectedOutcomeID.Valid {
		outcomeID := int(expectedOutcomeID.Int64)
		task.ExpectedOutcomeID = &outcomeID
	}
	return task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE tasks SET expected_outcome_id = ?, title = ?, description = ?, priority = ?, status = ?, due_date = ?, completed_at = ?, assigned_to_id = ?, task_type = ?, requested_by = ?, completion_percentage = ? WHERE id = ?`
	_, err := db.Exec(query, task.ExpectedOutcomeID, task.Title, task.Description, task.Priority, task.Status, task.DueDate, task.CompletedAt, task.AssignedToID, task.TaskType, task.RequestedBy, task.CompletionPercentage, task.ID)
	return err
}

func DeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

func GetTaskCountsByStatus(userID int) (completed, pending int, err error) {
	query := `SELECT 
		SUM(CASE WHEN status = 'Completed' THEN 1 ELSE 0 END) as completed,
		SUM(CASE WHEN status != 'Completed' THEN 1 ELSE 0 END) as pending
		FROM tasks WHERE user_id = ?`
	err = db.QueryRow(query, userID).Scan(&completed, &pending)
	return
}

// Get tasks assigned to a specific user
func GetTasksAssignedToUser(userID int) ([]Task, error) {
	query := `SELECT id, expected_outcome_id, user_id, title, description, priority, status, due_date, created_at, completed_at, assigned_to_id, task_type, requested_by, completion_percentage FROM tasks WHERE assigned_to_id = ? ORDER BY due_date ASC, created_at DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var completedAt sql.NullTime
		var assignedToID sql.NullInt64
		var expectedOutcomeID sql.NullInt64
		err := rows.Scan(&task.ID, &expectedOutcomeID, &task.UserID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.DueDate, &task.CreatedAt, &completedAt, &assignedToID, &task.TaskType, &task.RequestedBy, &task.CompletionPercentage)
		if err != nil {
			return nil, err
		}
		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}
		if assignedToID.Valid {
			assignedID := int(assignedToID.Int64)
			task.AssignedToID = &assignedID
		}
		if expectedOutcomeID.Valid {
			outcomeID := int(expectedOutcomeID.Int64)
			task.ExpectedOutcomeID = &outcomeID
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// Get all tasks (created by user OR assigned to user)
func GetAllUserTasks(userID int) ([]Task, error) {
	query := `SELECT id, expected_outcome_id, user_id, title, description, priority, status, due_date, created_at, completed_at, assigned_to_id, task_type, requested_by, completion_percentage FROM tasks WHERE user_id = ? OR assigned_to_id = ? ORDER BY due_date ASC, created_at DESC`
	rows, err := db.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var completedAt sql.NullTime
		var assignedToID sql.NullInt64
		var expectedOutcomeID sql.NullInt64
		err := rows.Scan(&task.ID, &expectedOutcomeID, &task.UserID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.DueDate, &task.CreatedAt, &completedAt, &assignedToID, &task.TaskType, &task.RequestedBy, &task.CompletionPercentage)
		if err != nil {
			return nil, err
		}
		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}
		if assignedToID.Valid {
			assignedID := int(assignedToID.Int64)
			task.AssignedToID = &assignedID
		}
		if expectedOutcomeID.Valid {
			outcomeID := int(expectedOutcomeID.Int64)
			task.ExpectedOutcomeID = &outcomeID
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
		}
		if assignedToID.Valid {
			assignedID := int(assignedToID.Int64)
			task.AssignedToID = &assignedID
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// Comment CRUD operations
func CreateComment(comment *Comment) error {
	query := `INSERT INTO comments (objective_id, activity_id, user_id, comment_text) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, comment.ObjectiveID, comment.ActivityID, comment.UserID, comment.CommentText)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	comment.ID = int(id)
	return nil
}

func GetCommentsByObjective(objectiveID int) ([]Comment, error) {
	query := `
		SELECT c.id, c.objective_id, c.activity_id, c.user_id, c.comment_text, c.created_at,
		       u.username, u.role
		FROM comments c
		INNER JOIN users u ON c.user_id = u.id
		WHERE c.objective_id = ? AND c.activity_id IS NULL
		ORDER BY c.created_at DESC
	`
	rows, err := db.Query(query, objectiveID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		var activityID sql.NullInt64
		err := rows.Scan(&comment.ID, &comment.ObjectiveID, &activityID, &comment.UserID, &comment.CommentText, &comment.CreatedAt, &comment.Username, &comment.UserRole)
		if err != nil {
			return nil, err
		}
		if activityID.Valid {
			actID := int(activityID.Int64)
			comment.ActivityID = &actID
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func GetCommentsByActivity(activityID int) ([]Comment, error) {
	query := `
		SELECT c.id, c.objective_id, c.activity_id, c.user_id, c.comment_text, c.created_at,
		       u.username, u.role
		FROM comments c
		INNER JOIN users u ON c.user_id = u.id
		WHERE c.activity_id = ?
		ORDER BY c.created_at DESC
	`
	rows, err := db.Query(query, activityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		var activityID sql.NullInt64
		err := rows.Scan(&comment.ID, &comment.ObjectiveID, &activityID, &comment.UserID, &comment.CommentText, &comment.CreatedAt, &comment.Username, &comment.UserRole)
		if err != nil {
			return nil, err
		}
		if activityID.Valid {
			actID := int(activityID.Int64)
			comment.ActivityID = &actID
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func DeleteComment(id int) error {
	query := `DELETE FROM comments WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}
// Get tasks by expected outcome ID
func GetTasksByExpectedOutcome(expectedOutcomeID int) ([]Task, error) {
	query := `SELECT id, expected_outcome_id, user_id, title, description, priority, status, due_date, created_at, completed_at, assigned_to_id, task_type, requested_by, completion_percentage FROM tasks WHERE expected_outcome_id = ? ORDER BY due_date ASC, created_at DESC`
	rows, err := db.Query(query, expectedOutcomeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var completedAt sql.NullTime
		var assignedToID sql.NullInt64
		var expectedOutcomeID sql.NullInt64
		err := rows.Scan(&task.ID, &expectedOutcomeID, &task.UserID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.DueDate, &task.CreatedAt, &completedAt, &assignedToID, &task.TaskType, &task.RequestedBy, &task.CompletionPercentage)
		if err != nil {
			return nil, err
		}
		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}
		if assignedToID.Valid {
			assignedID := int(assignedToID.Int64)
			task.AssignedToID = &assignedID
		}
		if expectedOutcomeID.Valid {
			outcomeID := int(expectedOutcomeID.Int64)
			task.ExpectedOutcomeID = &outcomeID
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// Get tasks by objective ID (joins through expected outcomes)
func GetTasksByObjective(objectiveID int) ([]Task, error) {
	query := `
		SELECT t.id, t.expected_outcome_id, t.user_id, t.title, t.description, t.priority, t.status, 
		       t.due_date, t.created_at, t.completed_at, t.assigned_to_id, t.task_type, t.requested_by, 
		       t.completion_percentage
		FROM tasks t
		INNER JOIN expected_outcomes eo ON t.expected_outcome_id = eo.id
		WHERE eo.objective_id = ?
		ORDER BY t.due_date ASC, t.created_at DESC
	`
	rows, err := db.Query(query, objectiveID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var completedAt sql.NullTime
		var assignedToID sql.NullInt64
		var expectedOutcomeID sql.NullInt64
		err := rows.Scan(&task.ID, &expectedOutcomeID, &task.UserID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.DueDate, &task.CreatedAt, &completedAt, &assignedToID, &task.TaskType, &task.RequestedBy, &task.CompletionPercentage)
		if err != nil {
			return nil, err
		}
		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}
		if assignedToID.Valid {
			assignedID := int(assignedToID.Int64)
			task.AssignedToID = &assignedID
		}
		if expectedOutcomeID.Valid {
			outcomeID := int(expectedOutcomeID.Int64)
			task.ExpectedOutcomeID = &outcomeID
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// Calculate objective performance based on tasks and activities
func CalculateObjectivePerformance(objectiveID int) (float64, error) {
	// Get all tasks for this objective
	tasks, err := GetTasksByObjective(objectiveID)
	if err != nil {
		return 0, err
	}

	// Get all activities for this objective
	activities, err := GetActivitiesByObjective(objectiveID)
	if err != nil {
		return 0, err
	}

	totalItems := len(tasks) + len(activities)
	if totalItems == 0 {
		return 0, nil
	}

	var totalPercentage float64

	// Sum task completion percentages
	for _, task := range tasks {
		if task.Status == TaskStatusCompleted {
			totalPercentage += 100
		} else {
			totalPercentage += task.CompletionPercentage
		}
	}

	// Sum activity progress percentages
	for _, activity := range activities {
		totalPercentage += activity.ProgressPercentage
	}

	return totalPercentage / float64(totalItems), nil
}

// Get activities by objective (helper for performance calculation)
func GetActivitiesByObjective(objectiveID int) ([]Activity, error) {
	query := `
		SELECT a.id, a.expected_outcome_id, a.title, a.description, a.category, 
		       a.progress_percentage, a.implementation_level, a.created_at, a.updated_at
		FROM activities a
		INNER JOIN expected_outcomes eo ON a.expected_outcome_id = eo.id
		WHERE eo.objective_id = ?
		ORDER BY a.created_at DESC
	`
	rows, err := db.Query(query, objectiveID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []Activity
	for rows.Next() {
		var activity Activity
		err := rows.Scan(&activity.ID, &activity.ExpectedOutcomeID, &activity.Title, &activity.Description,
			&activity.Category, &activity.ProgressPercentage, &activity.ImplementationLevel,
			&activity.CreatedAt, &activity.UpdatedAt)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}
	return activities, nil
}// GetObjectivesWithOutcomes retrieves objectives with their expected outcomes for a user
func GetObjectivesWithOutcomes(userID int) ([]ObjectiveWithOutcomes, error) {
objectives, err := GetObjectivesByUserID(userID)
if err != nil {
return nil, err
}

var objectivesWithOutcomes []ObjectiveWithOutcomes
for _, obj := range objectives {
// Calculate and update objective performance
performance, err := CalculateObjectivePerformance(obj.ID)
if err == nil {
obj.Performance = performance
}

outcomes, err := GetExpectedOutcomesByObjectiveID(obj.ID)
if err != nil {
log.Println("Error fetching outcomes:", err)
continue
}

var outcomesWithActivities []ExpectedOutcomeWithActivities
for _, outcome := range outcomes {
activities, err := GetActivitiesByExpectedOutcomeID(outcome.ID)
if err != nil {
log.Println("Error fetching activities:", err)
continue
}

// Fetch tasks for this expected outcome
tasks, err := GetTasksByExpectedOutcome(outcome.ID)
if err != nil {
log.Println("Error fetching tasks for outcome:", err)
tasks = []Task{} // Empty list on error
}

outcomesWithActivities = append(outcomesWithActivities, ExpectedOutcomeWithActivities{
ExpectedOutcome: outcome,
Activities:      activities,
Tasks:           tasks,
})
}

objectivesWithOutcomes = append(objectivesWithOutcomes, ObjectiveWithOutcomes{
Objective:        obj,
ExpectedOutcomes: outcomesWithActivities,
})
}

return objectivesWithOutcomes, nil
}
