package main

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
)

// Staff list handler - display all staff members
func staffListHandler(w http.ResponseWriter, r *http.Request) {
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
	if err != nil || currentUser.Role != RoleAdmin {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	staff, err := GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Load supervisor names for each staff member
	for i := range staff {
		if staff[i].SupervisorID != nil {
			supervisor, err := GetUserByID(*staff[i].SupervisorID)
			if err == nil {
				staff[i].SupervisorName = supervisor.Username
			}
		}
	}

	data := StaffListData{
		Username: currentUser.Username,
		Staff:    staff,
	}

	tmpl := template.Must(template.ParseFiles("templates/staff_list.html"))
	tmpl.Execute(w, data)
}

// New staff handler - show staff registration form
func newStaffHandler(w http.ResponseWriter, r *http.Request) {
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
	if err != nil || currentUser.Role != RoleAdmin {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		role := r.FormValue("role")
		department := r.FormValue("department")
		position := r.FormValue("position")
		supervisorIDStr := r.FormValue("supervisor_id")

		var supervisorID *int
		if supervisorIDStr != "" {
			id, err := strconv.Atoi(supervisorIDStr)
			if err == nil {
				supervisorID = &id
			}
		}

		user := &User{
			Username:     username,
			Password:     password,
			Role:         UserRole(role),
			Department:   department,
			Position:     position,
			SupervisorID: supervisorID,
			CreatedAt:    time.Now(),
		}

		if err := CreateUser(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/staff", http.StatusSeeOther)
		return
	}

	// Get supervisors and admins for dropdown
	supervisors, err := GetUsersByRole(RoleSupervisor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	admins, err := GetUsersByRole(RoleAdmin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allSupervisors := append(supervisors, admins...)

	data := StaffFormData{
		Username:    currentUser.Username,
		IsEdit:      false,
		Supervisors: allSupervisors,
	}

	tmpl := template.Must(template.ParseFiles("templates/staff_form.html"))
	tmpl.Execute(w, data)
}

// Edit staff handler - modify existing staff member
func editStaffHandler(w http.ResponseWriter, r *http.Request) {
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
	if err != nil || currentUser.Role != RoleAdmin {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	staffIDStr := r.URL.Query().Get("id")
	staffID, err := strconv.Atoi(staffIDStr)
	if err != nil {
		http.Error(w, "Invalid staff ID", http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		role := r.FormValue("role")
		department := r.FormValue("department")
		position := r.FormValue("position")
		supervisorIDStr := r.FormValue("supervisor_id")

		var supervisorID *int
		if supervisorIDStr != "" {
			id, err := strconv.Atoi(supervisorIDStr)
			if err == nil {
				supervisorID = &id
			}
		}

		user := &User{
			ID:           staffID,
			Username:     username,
			Role:         UserRole(role),
			Department:   department,
			Position:     position,
			SupervisorID: supervisorID,
		}

		// Only update password if provided
		if password != "" {
			user.Password = password
		}

		if err := UpdateUser(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/staff", http.StatusSeeOther)
		return
	}

	// Get staff member details
	staff, err := GetUserByID(staffID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get supervisors and admins for dropdown
	supervisors, err := GetUsersByRole(RoleSupervisor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	admins, err := GetUsersByRole(RoleAdmin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	allSupervisors := append(supervisors, admins...)

	data := StaffFormData{
		Username:    currentUser.Username,
		IsEdit:      true,
		Staff:       staff,
		Supervisors: allSupervisors,
	}

	tmpl := template.Must(template.ParseFiles("templates/staff_form.html"))
	tmpl.Execute(w, data)
}

// Delete staff handler
func deleteStaffHandler(w http.ResponseWriter, r *http.Request) {
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
	if err != nil || currentUser.Role != RoleAdmin {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	staffIDStr := r.URL.Query().Get("id")
	staffID, err := strconv.Atoi(staffIDStr)
	if err != nil {
		http.Error(w, "Invalid staff ID", http.StatusBadRequest)
		return
	}

	// Prevent deleting self
	if staffID == userID {
		http.Error(w, "Cannot delete your own account", http.StatusBadRequest)
		return
	}

	if err := DeleteUser(staffID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/staff", http.StatusSeeOther)
}

// Public registration handler - for staff self-registration
func registrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		department := r.FormValue("department")
		position := r.FormValue("position")

		user := &User{
			Username:   username,
			Password:   password,
			Role:       RoleStaff, // Default role for self-registration
			Department: department,
			Position:   position,
			CreatedAt:  time.Now(),
		}

		if err := CreateUser(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Auto-login after registration
		SetSession(w, r, user)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/registration.html"))
	tmpl.Execute(w, nil)
}
