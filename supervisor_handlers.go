package main

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
)

// Supervisor dashboard - view all supervised staff
func supervisorDashboardHandler(w http.ResponseWriter, r *http.Request) {
	session, err := GetSessionValues(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, ok := session["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	currentUser, err := GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Only supervisors and admins can access
	if currentUser.Role != RoleSupervisor && currentUser.Role != RoleAdmin {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Get staff under this supervisor
	var staff []User
	if currentUser.Role == RoleAdmin {
		// Admins see all staff
		staff, err = GetAllUsers()
	} else {
		// Supervisors see only their staff
		staff, err = GetStaffBySupervisor(userID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate performance for each staff member
	for i := range staff {
		objectives, err := GetObjectivesByUserID(staff[i].ID)
		if err != nil {
			continue
		}

		var totalPerformance float64
		var count int
		for _, obj := range objectives {
			activities, err := GetActivitiesByObjectiveID(obj.ID)
			if err != nil {
				continue
			}

			var objTotal float64
			for _, act := range activities {
				objTotal += float64(act.ProgressPercentage)
			}

			if len(activities) > 0 {
				totalPerformance += objTotal / float64(len(activities))
				count++
			}
		}

		if count > 0 {
			staff[i].OverallPerformance = totalPerformance / float64(count)
		}
	}

	data := SupervisorDashboardData{
		Username: currentUser.Username,
		Role:     string(currentUser.Role),
		Staff:    staff,
	}

	tmpl := template.Must(template.ParseFiles("templates/supervisor_dashboard.html"))
	tmpl.Execute(w, data)
}

// View individual staff report
func viewStaffReportHandler(w http.ResponseWriter, r *http.Request) {
	session, err := GetSessionValues(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, ok := session["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	currentUser, err := GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Only supervisors and admins can access
	if currentUser.Role != RoleSupervisor && currentUser.Role != RoleAdmin {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	staffIDStr := r.URL.Query().Get("id")
	staffID, err := strconv.Atoi(staffIDStr)
	if err != nil {
		http.Error(w, "Invalid staff ID", http.StatusBadRequest)
		return
	}

	staff, err := GetUserByID(staffID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify supervisor relationship (unless admin)
	if currentUser.Role == RoleSupervisor {
		if staff.SupervisorID == nil || *staff.SupervisorID != userID {
			http.Error(w, "Access denied - not your supervisee", http.StatusForbidden)
			return
		}
	}

	// Get objectives with comments
	objectives, err := GetObjectivesByUserID(staffID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var objectivesWithComments []ObjectiveWithComments
	for _, obj := range objectives {
		outcomes, err := GetExpectedOutcomesByObjectiveID(obj.ID)
		if err != nil {
			continue
		}

		activities, err := GetActivitiesByObjectiveID(obj.ID)
		if err != nil {
			continue
		}

		comments, err := GetCommentsByObjective(obj.ID)
		if err != nil {
			comments = []Comment{} // Empty if error
		}

		objWithComments := ObjectiveWithComments{
			Objective:  obj,
			Outcomes:   outcomes,
			Activities: activities,
			Comments:   comments,
		}
		objectivesWithComments = append(objectivesWithComments, objWithComments)
	}

	// Get tasks
	tasks, err := GetTasksByUserID(staffID)
	if err != nil {
		tasks = []Task{} // Empty if error
	}

	data := struct {
		Username   string
		Staff      *User
		Objectives []ObjectiveWithComments
		Tasks      []Task
	}{
		Username:   currentUser.Username,
		Staff:      staff,
		Objectives: objectivesWithComments,
		Tasks:      tasks,
	}

	tmpl := template.Must(template.ParseFiles("templates/staff_report.html"))
	tmpl.Execute(w, data)
}

// Add comment handler
func addCommentHandler(w http.ResponseWriter, r *http.Request) {
	session, err := GetSessionValues(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, ok := session["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	currentUser, err := GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Only supervisors and admins can add comments
	if currentUser.Role != RoleSupervisor && currentUser.Role != RoleAdmin {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		objectiveIDStr := r.FormValue("objective_id")
		activityIDStr := r.FormValue("activity_id")
		commentText := r.FormValue("comment_text")
		staffIDStr := r.FormValue("staff_id")

		objectiveID, err := strconv.Atoi(objectiveIDStr)
		if err != nil {
			http.Error(w, "Invalid objective ID", http.StatusBadRequest)
			return
		}

		var activityID *int
		if activityIDStr != "" {
			id, err := strconv.Atoi(activityIDStr)
			if err == nil {
				activityID = &id
			}
		}

		comment := &Comment{
			ObjectiveID: &objectiveID,
			ActivityID:  activityID,
			UserID:      userID,
			CommentText: commentText,
			CreatedAt:   time.Now(),
		}

		if err := CreateComment(comment); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect back to staff report
		http.Redirect(w, r, "/supervisor/staff?id="+staffIDStr, http.StatusSeeOther)
		return
	}

	// GET - show comment form
	objectiveIDStr := r.URL.Query().Get("objective_id")
	activityIDStr := r.URL.Query().Get("activity_id")
	staffIDStr := r.URL.Query().Get("staff_id")

	objectiveID, err := strconv.Atoi(objectiveIDStr)
	if err != nil {
		http.Error(w, "Invalid objective ID", http.StatusBadRequest)
		return
	}

	objective, err := GetObjectiveByID(objectiveID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var activity *Activity
	if activityIDStr != "" {
		actID, err := strconv.Atoi(activityIDStr)
		if err == nil {
			activity, _ = GetActivityByID(actID)
		}
	}

	data := CommentFormData{
		Username:  currentUser.Username,
		Objective: objective,
		Activity:  activity,
		StaffID:   staffIDStr,
	}

	tmpl := template.Must(template.ParseFiles("templates/comment_form.html"))
	tmpl.Execute(w, data)
}

// Delete comment handler
func deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	session, err := GetSessionValues(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, ok := session["userID"].(int)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	currentUser, err := GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Only supervisors and admins can delete comments
	if currentUser.Role != RoleSupervisor && currentUser.Role != RoleAdmin {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	commentIDStr := r.URL.Query().Get("id")
	staffIDStr := r.URL.Query().Get("staff_id")

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	if err := DeleteComment(commentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect back to staff report
	http.Redirect(w, r, "/supervisor/staff?id="+staffIDStr, http.StatusSeeOther)
}
