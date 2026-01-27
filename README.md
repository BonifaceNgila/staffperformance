# Staff Performance Management System

A comprehensive web-based system for managing staff performance through a modular dashboard with Tasks, Objectives, and Reports.

## 🎯 Features Overview

### Main Dashboard (Home)
- **Visual Navigation Menu** - Easy access to all sections
- **Real-time Statistics** - Objectives, tasks, and performance metrics
- **Quick Actions** - Shortcuts to create tasks and objectives

### 1. Tasks Management ✓
**Standalone task tracking system**
- Create, edit, and delete tasks
- **Priority Levels:** Low, Medium, High, Urgent
- **Status Tracking:** Pending, In Progress, Completed, On Hold
- **Due Date Management:** Track deadlines
- **Interactive Filters:** Filter by status
- **Task Viewer:** Visual card-based interface with color-coded badges
- Automatic completion date tracking

### 2. Objectives Management 🎯
**Strategic performance tracking**
- Create, edit, and delete performance objectives
- Set objective title, description, start date, and end date
- Add multiple expected outcomes per objective
- Track activities for each outcome with:
  - **Categories:** Daily, Weekly, Monthly, Quarterly, Biannually, Annually
  - **Progress Percentage:** 0-100%
  - **Implementation Level:** Detailed status descriptions
- Automatic performance calculation (mean of all activities)
- Visual performance indicators

### 3. Reports & Analytics 📊
**Comprehensive performance reporting**
- Executive summary with key metrics
- Detailed objective performance breakdown
- Tasks summary with priorities and statuses
- Average performance calculation
- Print/PDF export capability
- Performance visualization

## Technical Stack

- **Backend:** Go (Golang)
- **Database:** SQLite (pure Go implementation - modernc.org/sqlite)
- **Session Management:** Gorilla Sessions
- **Frontend:** HTML Templates, CSS
- **Server:** Built-in Go HTTP server

## Installation & Setup

### Prerequisites
- Go 1.21 or higher
- No external dependencies (pure Go SQLite driver)

### Installation Steps

1. **Clone or navigate to the project directory:**
   ```bash
   cd c:\xampp\staffperformance
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the server:**
   ```bash
   go run .
   ```

4. **Access the application:**
   - Open your browser and go to: `http://localhost:8080`
   - Login with default credentials:
     - Username: `admin`
     - Password: `admin123`

## Project Structure

```
staffperformance/
├── main.go                        # Main application entry point and routes
├── models.go                      # Data models and structures
├── database.go                    # Database operations and CRUD functions
├── handlers.go                    # Objectives handlers
├── task_handlers.go              # Tasks and reports handlers
├── session.go                     # Session management
├── go.mod                         # Go module dependencies
├── README.md                      # This file
├── USER_GUIDE.md                 # Comprehensive user guide
├── WHATS_NEW.md                  # Latest updates and features
├── QUICKSTART.md                 # Quick start guide
├── static/
│   └── css/
│       └── style.css             # Application styles
└── templates/
    ├── login.html                # Login page
    ├── dashboard.html            # Main dashboard/home
    ├── tasks.html                # Task list viewer
    ├── task_form.html            # Task create/edit form
    ├── objectives.html           # Objectives view
    ├── objective_form.html       # Objective create/edit form
    ├── expected_outcome_form.html # Outcome create/edit form
    ├── activity_form.html        # Activity create/edit form
    └── reports.html              # Reports and analytics
```

## Database Schema

### Users Table
- id (Primary Key)
- username (Unique)
- password
- full_name
- email
- created_at

### Objectives Table
- id (Primary Key)
- user_id (Foreign Key → users.id)
- title
- description
- start_date
- end_date
- created_at

### Expected Outcomes Table
- id (Primary Key)
- objective_id (Foreign Key → objectives.id)
- title
- description
- created_at

### Activities Table
- id (Primary Key)
- expected_outcome_id (Foreign Key → expected_outcomes.id)
- title
- description
- category (Daily/Weekly/Monthly/Quarterly/Biannually/Annually)
- progress_percentage (0-100)
- implementation_level (Text description)
- created_at
- updated_at

### Tasks Table (New!)
- id (Primary Key)
- user_id (Foreign Key → users.id)
- title
- description
- priority (Low/Medium/High/Urgent)
- status (Pending/In Progress/Completed/On Hold)
- due_date
- created_at
- completed_at

## Usage Workflow

### Quick Start
1. **Login** to the system with your credentials
2. **View Dashboard** - See overview statistics
3. Choose your path:
   - **Tasks** for immediate action items
   - **Objectives** for strategic goals
   - **Reports** for performance review

### Task Management Workflow
1. Click **"Tasks"** from the main menu
2. Click **"+ New Task"** 
3. Fill in task details (title, description, priority, status, due date)
4. Click **"Create Task"**
5. Use filters to view tasks by status
6. Edit tasks to update status and track progress
7. Mark complete when done

### Objectives Workflow
1. Click **"Objectives"** from the main menu
2. Create an **Objective** with title, description, and date range
3. Add **Expected Outcomes** for each objective
4. Create **Activities/Tasks** for each expected outcome with:
   - Title and description
   - Category (frequency)
   - Progress percentage
   - Implementation level details
5. Update activity progress regularly
6. Monitor objective performance (automatically calculated)

### Reporting Workflow
1. Click **"Reports"** from the main menu
2. Review summary statistics
3. Analyze objective performance details
4. Check task completion status
5. Print or export reports as needed

## Key Concepts

### Performance Calculation
The objective performance is calculated using the formula:

```
Objective Performance = (Sum of all activity progress percentages) / (Number of activities)
```

This provides a clear, quantifiable measure of progress toward each objective.

### Hierarchical Structure
```
User
└── Objectives
    └── Expected Outcomes
        └── Activities/Tasks
            ├── Category
            ├── Progress %
            └── Implementation Level
```

## API Routes

### Public Routes
- `GET /` - Login page
- `POST /login` - Process login
- `GET /logout` - Logout

### Protected Routes (Require Authentication)
- `GET /dashboard` - Main dashboard/home

#### Tasks
- `GET /tasks` - Task list viewer
- `GET /tasks/new` - Create new task
- `GET /tasks/edit?id=X` - Edit task
- `GET /tasks/delete?id=X` - Delete task

#### Objectives
- `GET /objectives` - Objectives view
- `GET /objectives/new` - Create new objective
- `GET /objectives/edit?id=X` - Edit objective
- `GET /objectives/delete?id=X` - Delete objective

#### Expected Outcomes
- `GET /outcomes/new?objective_id=X` - Create new outcome
- `GET /outcomes/edit?id=X` - Edit outcome
- `GET /outcomes/delete?id=X` - Delete outcome

#### Activities
- `GET /activities/new?outcome_id=X` - Create new activity
- `GET /activities/edit?id=X` - Edit activity
- `GET /activities/delete?id=X` - Delete activity

#### Reports
- `GET /reports` - Performance reports and analytics

## Security Features

- Session-based authentication
- Password verification
- Route protection middleware
- Ownership verification for all CRUD operations
- HTTP-only cookies

## Future Enhancements

- User registration and profile management
- Password hashing (bcrypt)
- Role-based access control (Admin, Manager, Staff)
- Team collaboration features
- Task assignment to other users
- Email notifications for deadlines
- Calendar integration
- Performance analytics with charts/graphs
- Export reports to Excel/CSV
- Recurring tasks
- Task templates
- Comment threads on objectives and tasks
- File attachments
- Mobile app

---

**Version:** 2.0.0  
**Last Updated:** January 27, 2026  
**New Features:** Tasks Management, Reports Section, Modular Dashboard

## Related Documentation

- [USER_GUIDE.md](USER_GUIDE.md) - Comprehensive user guide
- [WHATS_NEW.md](WHATS_NEW.md) - Latest updates and changes
- [QUICKSTART.md](QUICKSTART.md) - Quick start guide
