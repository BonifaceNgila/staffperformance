package main

import (
	"html/template"
	"log"
	"net/http"
)

var templates *template.Template

func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	// Initialize database
	if err := InitDB(); err != nil {
		log.Fatal("Database initialization failed:", err)
	}
	defer db.Close()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Public routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/register", registrationHandler)

	// Protected routes
	http.HandleFunc("/dashboard", RequireAuth(dashboardHandler))

	// Main menu routes
	http.HandleFunc("/tasks", RequireAuth(tasksHandler))
	http.HandleFunc("/reports", RequireAuth(reportsHandler))
	http.HandleFunc("/objectives", RequireAuth(objectivesPageHandler))

	// Task routes
	http.HandleFunc("/tasks/new", RequireAuth(newTaskHandler))
	http.HandleFunc("/tasks/edit", RequireAuth(editTaskHandler))
	http.HandleFunc("/tasks/delete", RequireAuth(deleteTaskHandler))

	// Objective routes
	http.HandleFunc("/objectives/new", RequireAuth(newObjectiveHandler))
	http.HandleFunc("/objectives/edit", RequireAuth(editObjectiveHandler))
	http.HandleFunc("/objectives/delete", RequireAuth(deleteObjectiveHandler))

	// Expected Outcome routes
	http.HandleFunc("/outcomes/new", RequireAuth(newExpectedOutcomeHandler))
	http.HandleFunc("/outcomes/edit", RequireAuth(editExpectedOutcomeHandler))
	http.HandleFunc("/outcomes/delete", RequireAuth(deleteExpectedOutcomeHandler))

	// Activity routes
	http.HandleFunc("/activities/new", RequireAuth(newActivityHandler))
	http.HandleFunc("/activities/edit", RequireAuth(editActivityHandler))
	http.HandleFunc("/activities/delete", RequireAuth(deleteActivityHandler))

	// Staff management routes (Admin only)
	http.HandleFunc("/staff", RequireAuth(staffListHandler))
	http.HandleFunc("/staff/new", RequireAuth(newStaffHandler))
	http.HandleFunc("/staff/edit", RequireAuth(editStaffHandler))
	http.HandleFunc("/staff/delete", RequireAuth(deleteStaffHandler))

	// Supervisor routes
	http.HandleFunc("/supervisor/dashboard", RequireAuth(supervisorDashboardHandler))
	http.HandleFunc("/supervisor/staff", RequireAuth(viewStaffReportHandler))

	// Comment routes
	http.HandleFunc("/comments/new", RequireAuth(addCommentHandler))
	http.HandleFunc("/comments/delete", RequireAuth(deleteCommentHandler))

	log.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Check if user is already logged in
	user, err := GetSession(r)
	if err == nil && user != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	err = templates.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := GetUserByUsername(username)
	if err != nil || user.Password != password {
		log.Printf("Login failed for user: %s", username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Set session
	err = SetSession(w, r, user)
	if err != nil {
		log.Println("Error setting session:", err)
		http.Error(w, "Login error", http.StatusInternalServerError)
		return
	}

	log.Printf("Login successful for user: %s", username)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
