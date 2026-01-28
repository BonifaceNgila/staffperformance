package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

// Tasks management handlers
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tasks, err := GetTasksByUserID(user.ID)
	if err != nil {
		log.Println("Error fetching tasks:", err)
		http.Error(w, "Error loading tasks", http.StatusInternalServerError)
		return
	}

	priorities := []TaskPriority{PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent}
	statuses := []TaskStatus{TaskStatusPending, TaskStatusInProgress, TaskStatusCompleted, TaskStatusOnHold}

	data := TaskListData{
		User:       *user,
		Tasks:      tasks,
		Priorities: priorities,
		Statuses:   statuses,
	}

	err = templates.ExecuteTemplate(w, "tasks.html", data)
	if err != nil {
		log.Println("Template error:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func newTaskHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		description := r.FormValue("description")
		priority := TaskPriority(r.FormValue("priority"))
		status := TaskStatus(r.FormValue("status"))
		dueDateStr := r.FormValue("due_date")
		expectedOutcomeIDStr := r.FormValue("expected_outcome_id")
		completionPercentageStr := r.FormValue("completion_percentage")

		dueDate, _ := time.Parse("2006-01-02", dueDateStr)

		task := &Task{
			UserID:      user.ID,
			Title:       title,
			Description: description,
			Priority:    priority,
			Status:      status,
			DueDate:     dueDate,
		}

		// Parse expected outcome ID if provided
		if expectedOutcomeIDStr != "" {
			expectedOutcomeID, err := strconv.Atoi(expectedOutcomeIDStr)
			if err == nil {
				task.ExpectedOutcomeID = &expectedOutcomeID
			}
		}

		// Parse completion percentage
		if completionPercentageStr != "" {
			completionPercentage, err := strconv.ParseFloat(completionPercentageStr, 64)
			if err == nil && completionPercentage >= 0 && completionPercentage <= 100 {
				task.CompletionPercentage = completionPercentage
			}
		}

		// Set completed time if status is completed
		if status == TaskStatusCompleted {
			now := time.Now()
			task.CompletedAt = &now
			task.CompletionPercentage = 100 // Auto-set to 100% when completed
		}

		err = CreateTask(task)
		if err != nil {
			log.Println("Error creating task:", err)
			http.Error(w, "Error creating task", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
		return
	}

	// Get user's objectives with expected outcomes for the dropdown
	objectives, err := GetObjectivesWithOutcomes(user.ID)
	if err != nil {
		log.Println("Error fetching objectives:", err)
		http.Error(w, "Error loading objectives", http.StatusInternalServerError)
		return
	}

	priorities := []TaskPriority{PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent}
	statuses := []TaskStatus{TaskStatusPending, TaskStatusInProgress, TaskStatusCompleted, TaskStatusOnHold}
	taskTypes := []TaskType{TaskTypePersonal, TaskTypeServiceRequest, TaskTypeStaffAssignment, TaskTypeResponse}

	data := TaskFormData{
		User:       *user,
		Objectives: objectives,
		Priorities: priorities,
		Statuses:   statuses,
		TaskTypes:  taskTypes,
		IsEdit:     false,
	}

	err = templates.ExecuteTemplate(w, "task_form.html", data)
	if err != nil {
		log.Println("Template error:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := GetTaskByID(id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Verify ownership
	if task.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if r.Method == http.MethodPost {
		task.Title = r.FormValue("title")
		task.Description = r.FormValue("description")
		task.Priority = TaskPriority(r.FormValue("priority"))
		newStatus := TaskStatus(r.FormValue("status"))
		dueDateStr := r.FormValue("due_date")
		expectedOutcomeIDStr := r.FormValue("expected_outcome_id")
		completionPercentageStr := r.FormValue("completion_percentage")

		task.DueDate, _ = time.Parse("2006-01-02", dueDateStr)

		// Parse expected outcome ID if provided
		if expectedOutcomeIDStr != "" {
			expectedOutcomeID, err := strconv.Atoi(expectedOutcomeIDStr)
			if err == nil {
				task.ExpectedOutcomeID = &expectedOutcomeID
			}
		} else {
			task.ExpectedOutcomeID = nil
		}

		// Parse completion percentage
		if completionPercentageStr != "" {
			completionPercentage, err := strconv.ParseFloat(completionPercentageStr, 64)
			if err == nil && completionPercentage >= 0 && completionPercentage <= 100 {
				task.CompletionPercentage = completionPercentage
			}
		}

		// Update completed time if status changed to completed
		if newStatus == TaskStatusCompleted && task.Status != TaskStatusCompleted {
			now := time.Now()
			task.CompletedAt = &now
			task.CompletionPercentage = 100 // Auto-set to 100% when completed
		} else if newStatus != TaskStatusCompleted {
			task.CompletedAt = nil
		}

		task.Status = newStatus

		err = UpdateTask(task)
		if err != nil {
			log.Println("Error updating task:", err)
			http.Error(w, "Error updating task", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/tasks", http.StatusSeeOther)
		return
	}

	// Get user's objectives with expected outcomes for the dropdown
	objectives, err := GetObjectivesWithOutcomes(user.ID)
	if err != nil {
		log.Println("Error fetching objectives:", err)
		http.Error(w, "Error loading objectives", http.StatusInternalServerError)
		return
	}

	priorities := []TaskPriority{PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent}
	statuses := []TaskStatus{TaskStatusPending, TaskStatusInProgress, TaskStatusCompleted, TaskStatusOnHold}
	taskTypes := []TaskType{TaskTypePersonal, TaskTypeServiceRequest, TaskTypeStaffAssignment, TaskTypeResponse}

	data := TaskFormData{
		User:       *user,
		Task:       task,
		Objectives: objectives,
		Priorities: priorities,
		Statuses:   statuses,
		TaskTypes:  taskTypes,
		IsEdit:     true,
	}

	err = templates.ExecuteTemplate(w, "task_form.html", data)
	if err != nil {
		log.Println("Template error:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := GetTaskByID(id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Verify ownership
	if task.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	err = DeleteTask(id)
	if err != nil {
		log.Println("Error deleting task:", err)
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/tasks", http.StatusSeeOther)
}

// Reports handler
func reportsHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get objectives with full data
	objectives, err := GetObjectivesByUserID(user.ID)
	if err != nil {
		log.Println("Error fetching objectives:", err)
		http.Error(w, "Error loading reports", http.StatusInternalServerError)
		return
	}

	var objectivesWithOutcomes []ObjectiveWithOutcomes
	var totalPerformance float64
	for _, obj := range objectives {
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
			outcomesWithActivities = append(outcomesWithActivities, ExpectedOutcomeWithActivities{
				ExpectedOutcome: outcome,
				Activities:      activities,
			})
		}

		objectivesWithOutcomes = append(objectivesWithOutcomes, ObjectiveWithOutcomes{
			Objective:        obj,
			ExpectedOutcomes: outcomesWithActivities,
		})

		totalPerformance += obj.Performance
	}

	// Get tasks
	tasks, err := GetTasksByUserID(user.ID)
	if err != nil {
		log.Println("Error fetching tasks:", err)
		tasks = []Task{}
	}

	completedTasks, pendingTasks, _ := GetTaskCountsByStatus(user.ID)

	avgPerformance := 0.0
	if len(objectives) > 0 {
		avgPerformance = totalPerformance / float64(len(objectives))
	}

	data := ReportData{
		User:               *user,
		Objectives:         objectivesWithOutcomes,
		Tasks:              tasks,
		TotalObjectives:    len(objectives),
		CompletedTasks:     completedTasks,
		PendingTasks:       pendingTasks,
		AveragePerformance: avgPerformance,
	}

	err = templates.ExecuteTemplate(w, "reports.html", data)
	if err != nil {
		log.Println("Template error:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// Objectives page (separate from dashboard)
func objectivesPageHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	objectives, err := GetObjectivesByUserID(user.ID)
	if err != nil {
		log.Println("Error fetching objectives:", err)
		http.Error(w, "Error loading objectives", http.StatusInternalServerError)
		return
	}

	// Build dashboard data with full hierarchy
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

	data := DashboardData{
		User:       *user,
		Objectives: objectivesWithOutcomes,
	}

	err = templates.ExecuteTemplate(w, "objectives.html", data)
	if err != nil {
		log.Println("Template error:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
