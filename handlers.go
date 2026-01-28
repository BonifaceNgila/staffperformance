package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

// Dashboard handler
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get summary statistics
	objectives, _ := GetObjectivesByUserID(user.ID)
	tasks, _ := GetTasksByUserID(user.ID)
	completedTasks, pendingTasks, _ := GetTaskCountsByStatus(user.ID)

	var totalPerformance float64
	for _, obj := range objectives {
		totalPerformance += obj.Performance
	}
	avgPerformance := 0.0
	if len(objectives) > 0 {
		avgPerformance = totalPerformance / float64(len(objectives))
	}

	data := struct {
		User               User
		TotalObjectives    int
		TotalTasks         int
		CompletedTasks     int
		PendingTasks       int
		AveragePerformance float64
	}{
		User:               *user,
		TotalObjectives:    len(objectives),
		TotalTasks:         len(tasks),
		CompletedTasks:     completedTasks,
		PendingTasks:       pendingTasks,
		AveragePerformance: avgPerformance,
	}

	err = templates.ExecuteTemplate(w, "dashboard.html", data)
	if err != nil {
		log.Println("Template error:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// Objective handlers
func newObjectiveHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		description := r.FormValue("description")
		startDate := r.FormValue("start_date")
		endDate := r.FormValue("end_date")
		visibility := ObjectiveVisibility(r.FormValue("visibility"))
		status := ObjectiveStatus(r.FormValue("status"))
		category := ObjectiveCategory(r.FormValue("category"))
		categoryOther := r.FormValue("category_other")
		weightStr := r.FormValue("weight")

		start, _ := time.Parse("2006-01-02", startDate)
		end, _ := time.Parse("2006-01-02", endDate)
		weight, _ := strconv.ParseFloat(weightStr, 64)

		obj := &Objective{
			UserID:        user.ID,
			Title:         title,
			Description:   description,
			StartDate:     start,
			EndDate:       end,
			Visibility:    visibility,
			Status:        status,
			Category:      category,
			CategoryOther: categoryOther,
			Weight:        weight,
		}

		err = CreateObjective(obj)
		if err != nil {
			log.Println("Error creating objective:", err)
			http.Error(w, "Error creating objective", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	data := ObjectiveFormData{
		User:   *user,
		IsEdit: false,
	}

	err = templates.ExecuteTemplate(w, "objective_form.html", data)
	if err != nil {
		log.Println("Template error:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func editObjectiveHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid objective ID", http.StatusBadRequest)
		return
	}

	obj, err := GetObjectiveByID(id)
	if err != nil {
		http.Error(w, "Objective not found", http.StatusNotFound)
		return
	}

	// Verify ownership
	if obj.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if r.Method == http.MethodPost {
		obj.Title = r.FormValue("title")
		obj.Description = r.FormValue("description")
		startDate := r.FormValue("start_date")
		endDate := r.FormValue("end_date")
		obj.Visibility = ObjectiveVisibility(r.FormValue("visibility"))
		obj.Status = ObjectiveStatus(r.FormValue("status"))
		obj.Category = ObjectiveCategory(r.FormValue("category"))
		obj.CategoryOther = r.FormValue("category_other")
		weightStr := r.FormValue("weight")

		obj.StartDate, _ = time.Parse("2006-01-02", startDate)
		obj.EndDate, _ = time.Parse("2006-01-02", endDate)
		obj.Weight, _ = strconv.ParseFloat(weightStr, 64)

		err = UpdateObjective(obj)
		if err != nil {
			log.Println("Error updating objective:", err)
			http.Error(w, "Error updating objective", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	data := ObjectiveFormData{
		User:      *user,
		Objective: obj,
		IsEdit:    true,
	}

	err = templates.ExecuteTemplate(w, "objective_form.html", data)
	if err != nil {
		log.Println("Template error:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func deleteObjectiveHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid objective ID", http.StatusBadRequest)
		return
	}

	obj, err := GetObjectiveByID(id)
	if err != nil {
		http.Error(w, "Objective not found", http.StatusNotFound)
		return
	}

	// Verify ownership
	if obj.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	err = DeleteObjective(id)
	if err != nil {
		log.Println("Error deleting objective:", err)
		http.Error(w, "Error deleting objective", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Expected Outcome handlers
func newExpectedOutcomeHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	objIDStr := r.URL.Query().Get("objective_id")
	objID, err := strconv.Atoi(objIDStr)
	if err != nil {
		http.Error(w, "Invalid objective ID", http.StatusBadRequest)
		return
	}

	obj, err := GetObjectiveByID(objID)
	if err != nil {
		http.Error(w, "Objective not found", http.StatusNotFound)
		return
	}

	// Verify ownership
	if obj.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		description := r.FormValue("description")

		outcome := &ExpectedOutcome{
			ObjectiveID: objID,
			Title:       title,
			Description: description,
		}

		err = CreateExpectedOutcome(outcome)
		if err != nil {
			log.Println("Error creating expected outcome:", err)
			http.Error(w, "Error creating expected outcome", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	data := ExpectedOutcomeFormData{
		User:      *user,
		Objective: *obj,
		IsEdit:    false,
	}

	err = templates.ExecuteTemplate(w, "expected_outcome_form.html", data)
	if err != nil {
		log.Println("Template error:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func editExpectedOutcomeHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid expected outcome ID", http.StatusBadRequest)
		return
	}

	outcome, err := GetExpectedOutcomeByID(id)
	if err != nil {
		http.Error(w, "Expected outcome not found", http.StatusNotFound)
		return
	}

	obj, err := GetObjectiveByID(outcome.ObjectiveID)
	if err != nil {
		http.Error(w, "Objective not found", http.StatusNotFound)
		return
	}

	// Verify ownership
	if obj.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if r.Method == http.MethodPost {
		outcome.Title = r.FormValue("title")
		outcome.Description = r.FormValue("description")

		err = UpdateExpectedOutcome(outcome)
		if err != nil {
			log.Println("Error updating expected outcome:", err)
			http.Error(w, "Error updating expected outcome", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	data := ExpectedOutcomeFormData{
		User:            *user,
		Objective:       *obj,
		ExpectedOutcome: outcome,
		IsEdit:          true,
	}

	err = templates.ExecuteTemplate(w, "expected_outcome_form.html", data)
	if err != nil {
		log.Println("Template error:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func deleteExpectedOutcomeHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid expected outcome ID", http.StatusBadRequest)
		return
	}

	outcome, err := GetExpectedOutcomeByID(id)
	if err != nil {
		http.Error(w, "Expected outcome not found", http.StatusNotFound)
		return
	}

	obj, err := GetObjectiveByID(outcome.ObjectiveID)
	if err != nil {
		http.Error(w, "Objective not found", http.StatusNotFound)
		return
	}

	// Verify ownership
	if obj.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	err = DeleteExpectedOutcome(id)
	if err != nil {
		log.Println("Error deleting expected outcome:", err)
		http.Error(w, "Error deleting expected outcome", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Activity handlers
func newActivityHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	outcomeIDStr := r.URL.Query().Get("outcome_id")
	outcomeID, err := strconv.Atoi(outcomeIDStr)
	if err != nil {
		http.Error(w, "Invalid expected outcome ID", http.StatusBadRequest)
		return
	}

	outcome, err := GetExpectedOutcomeByID(outcomeID)
	if err != nil {
		http.Error(w, "Expected outcome not found", http.StatusNotFound)
		return
	}

	obj, err := GetObjectiveByID(outcome.ObjectiveID)
	if err != nil {
		http.Error(w, "Objective not found", http.StatusNotFound)
		return
	}

	// Verify ownership
	if obj.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		description := r.FormValue("description")
		category := ActivityCategory(r.FormValue("category"))
		progressStr := r.FormValue("progress_percentage")
		implementationLevel := r.FormValue("implementation_level")

		progress, _ := strconv.ParseFloat(progressStr, 64)

		activity := &Activity{
			ExpectedOutcomeID:   outcomeID,
			Title:               title,
			Description:         description,
			Category:            category,
			ProgressPercentage:  progress,
			ImplementationLevel: implementationLevel,
		}

		err = CreateActivity(activity)
		if err != nil {
			log.Println("Error creating activity:", err)
			http.Error(w, "Error creating activity", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	categories := []ActivityCategory{
		CategoryDaily,
		CategoryWeekly,
		CategoryMonthly,
		CategoryQuarterly,
		CategoryBiannually,
		CategoryAnnually,
	}

	data := ActivityFormData{
		User:            *user,
		Objective:       *obj,
		ExpectedOutcome: *outcome,
		Categories:      categories,
		IsEdit:          false,
	}

	err = templates.ExecuteTemplate(w, "activity_form.html", data)
	if err != nil {
		log.Println("Template error:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func editActivityHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid activity ID", http.StatusBadRequest)
		return
	}

	activity, err := GetActivityByID(id)
	if err != nil {
		http.Error(w, "Activity not found", http.StatusNotFound)
		return
	}

	outcome, err := GetExpectedOutcomeByID(activity.ExpectedOutcomeID)
	if err != nil {
		http.Error(w, "Expected outcome not found", http.StatusNotFound)
		return
	}

	obj, err := GetObjectiveByID(outcome.ObjectiveID)
	if err != nil {
		http.Error(w, "Objective not found", http.StatusNotFound)
		return
	}

	// Verify ownership
	if obj.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if r.Method == http.MethodPost {
		activity.Title = r.FormValue("title")
		activity.Description = r.FormValue("description")
		activity.Category = ActivityCategory(r.FormValue("category"))
		progressStr := r.FormValue("progress_percentage")
		activity.ImplementationLevel = r.FormValue("implementation_level")

		activity.ProgressPercentage, _ = strconv.ParseFloat(progressStr, 64)

		err = UpdateActivity(activity)
		if err != nil {
			log.Println("Error updating activity:", err)
			http.Error(w, "Error updating activity", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	categories := []ActivityCategory{
		CategoryDaily,
		CategoryWeekly,
		CategoryMonthly,
		CategoryQuarterly,
		CategoryBiannually,
		CategoryAnnually,
	}

	data := ActivityFormData{
		User:            *user,
		Objective:       *obj,
		ExpectedOutcome: *outcome,
		Activity:        activity,
		Categories:      categories,
		IsEdit:          true,
	}

	err = templates.ExecuteTemplate(w, "activity_form.html", data)
	if err != nil {
		log.Println("Template error:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func deleteActivityHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid activity ID", http.StatusBadRequest)
		return
	}

	activity, err := GetActivityByID(id)
	if err != nil {
		http.Error(w, "Activity not found", http.StatusNotFound)
		return
	}

	outcome, err := GetExpectedOutcomeByID(activity.ExpectedOutcomeID)
	if err != nil {
		http.Error(w, "Expected outcome not found", http.StatusNotFound)
		return
	}

	obj, err := GetObjectiveByID(outcome.ObjectiveID)
	if err != nil {
		http.Error(w, "Objective not found", http.StatusNotFound)
		return
	}

	// Verify ownership
	if obj.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	err = DeleteActivity(id)
	if err != nil {
		log.Println("Error deleting activity:", err)
		http.Error(w, "Error deleting activity", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Logout handler
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	err := ClearSession(w, r)
	if err != nil {
		log.Println("Error clearing session:", err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
