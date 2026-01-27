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
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
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
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		priority TEXT NOT NULL,
		status TEXT NOT NULL,
		due_date DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		completed_at DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`

	_, err := db.Exec(schema)
	return err
}

func createDefaultUser() error {
	query := `INSERT INTO users (username, password, full_name, email) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, "admin", "admin123", "Administrator", "admin@example.com")
	return err
}

// User CRUD operations
func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	query := `SELECT id, username, password, full_name, email, created_at FROM users WHERE username = ?`
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.FullName, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByID(id int) (*User, error) {
	user := &User{}
	query := `SELECT id, username, password, full_name, email, created_at FROM users WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password, &user.FullName, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
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
	query := `INSERT INTO tasks (user_id, title, description, priority, status, due_date) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, task.UserID, task.Title, task.Description, task.Priority, task.Status, task.DueDate)
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
	query := `SELECT id, user_id, title, description, priority, status, due_date, created_at, completed_at FROM tasks WHERE user_id = ? ORDER BY due_date ASC, created_at DESC`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var completedAt sql.NullTime
		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.DueDate, &task.CreatedAt, &completedAt)
		if err != nil {
			return nil, err
		}
		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func GetTaskByID(id int) (*Task, error) {
	task := &Task{}
	var completedAt sql.NullTime
	query := `SELECT id, user_id, title, description, priority, status, due_date, created_at, completed_at FROM tasks WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Priority, &task.Status, &task.DueDate, &task.CreatedAt, &completedAt)
	if err != nil {
		return nil, err
	}
	if completedAt.Valid {
		task.CompletedAt = &completedAt.Time
	}
	return task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE tasks SET title = ?, description = ?, priority = ?, status = ?, due_date = ?, completed_at = ? WHERE id = ?`
	_, err := db.Exec(query, task.Title, task.Description, task.Priority, task.Status, task.DueDate, task.CompletedAt, task.ID)
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
