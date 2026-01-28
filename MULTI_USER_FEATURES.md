# Staff Performance Management System - Multi-User Features

## Implementation Summary

### Features Implemented

#### 1. Multi-User System with Role-Based Access
- **User Roles**: Admin, Supervisor, Staff
- **Role-Based Permissions**:
  - Admin: Full access to all features including staff management
  - Supervisor: Can view and comment on assigned staff performance
  - Staff: Can manage their own objectives, tasks, and activities

#### 2. Staff Management (Admin Only)
- **Staff CRUD Operations**:
  - Create new staff members with role assignment
  - Edit existing staff (username, role, department, position, supervisor)
  - Delete staff members
  - View all staff in a table with filters
- **Routes**:
  - `/staff` - Staff list
  - `/staff/new` - Add new staff
  - `/staff/edit?id=X` - Edit staff
  - `/staff/delete?id=X` - Delete staff

#### 3. Supervisor Hierarchy
- Each staff member can be assigned to a supervisor
- Supervisors (and Admins) can:
  - View all their supervised staff
  - Monitor staff performance metrics
  - Access detailed staff reports
  - Add comments on objectives and activities

#### 4. Supervisor Dashboard
- **Features**:
  - View list of supervised staff (or all staff for admins)
  - See overall performance percentage for each staff member
  - Quick access to individual staff reports
- **Routes**:
  - `/supervisor/dashboard` - Main supervisor view
  - `/supervisor/staff?id=X` - Individual staff report

#### 5. Comments System
- Supervisors and admins can add comments on:
  - Objectives (general feedback)
  - Activities (specific task feedback)
- Comments display:
  - Commenter name and role
  - Timestamp
  - Comment text
- **Routes**:
  - `/comments/new` - Add comment
  - `/comments/delete?id=X` - Delete comment

#### 6. Task Assignment System
- **Task Types**:
  - Personal - User's own tasks
  - ServiceRequest - Request help from others
  - StaffAssignment - Task assigned to staff member
  - Response - Response to a service request
- **Task Assignment Fields**:
  - `assigned_to_id` - Who the task is assigned to
  - `task_type` - Type of task
  - `requested_by` - Who requested the task

#### 7. Self-Registration
- New staff can register themselves
- Default role: Staff
- **Route**: `/register`
- Auto-login after successful registration

### Database Schema Updates

#### Users Table
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    full_name TEXT,
    email TEXT,
    role TEXT NOT NULL DEFAULT 'Staff',
    supervisor_id INTEGER,
    department TEXT,
    position TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (supervisor_id) REFERENCES users(id)
)
```

#### Tasks Table (Enhanced)
```sql
CREATE TABLE tasks (
    -- ... existing fields ...
    assigned_to_id INTEGER,
    task_type TEXT,
    requested_by TEXT,
    FOREIGN KEY (assigned_to_id) REFERENCES users(id)
)
```

#### Comments Table (New)
```sql
CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    objective_id INTEGER,
    activity_id INTEGER,
    user_id INTEGER NOT NULL,
    comment_text TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (objective_id) REFERENCES objectives(id),
    FOREIGN KEY (activity_id) REFERENCES activities(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
)
```

### New Files Created

1. **staff_handlers.go** - Staff management handlers
2. **supervisor_handlers.go** - Supervisor dashboard and comment handlers
3. **templates/staff_list.html** - Staff listing page
4. **templates/staff_form.html** - Staff add/edit form
5. **templates/registration.html** - Self-registration page
6. **templates/supervisor_dashboard.html** - Supervisor team view
7. **templates/staff_report.html** - Individual staff performance report
8. **templates/comment_form.html** - Add comment form

### Updated Files

1. **models.go**:
   - Added `UserRole` enum (Admin, Supervisor, Staff)
   - Enhanced `User` struct with role, supervisor, department, position
   - Added `TaskType` enum
   - Added `Comment` struct
   - Added view models: `StaffListData`, `StaffFormData`, `SupervisorDashboardData`, `ObjectiveWithComments`, `CommentFormData`

2. **database.go**:
   - Updated `CreateUser`, `UpdateUser`, `DeleteUser`
   - Added `GetAllUsers`, `GetUsersByRole`, `GetStaffBySupervisor`
   - Updated task CRUD for assignment fields
   - Added `GetTasksAssignedToUser`, `GetAllUserTasks`
   - Added comment CRUD: `CreateComment`, `GetCommentsByObjective`, `GetCommentsByActivity`, `DeleteComment`

3. **session.go**:
   - Added `GetSessionValues` function for map-based session access

4. **main.go**:
   - Added routes for staff management
   - Added routes for supervisor features
   - Added comment routes
   - Added registration route

5. **templates/dashboard.html**:
   - Added role-based menu items
   - "My Team" menu for supervisors and admins
   - "Staff" menu for admins only

6. **templates/login.html**:
   - Added registration link

7. **static/css/style.css**:
   - Added styles for staff tables
   - Added comment section styles
   - Added role badges
   - Added performance indicators

### Usage Instructions

#### For Admins:
1. Login with default credentials (admin/admin)
2. Access "Staff" menu to manage users
3. Create supervisors and assign staff to them
4. Access "My Team" to view all staff performance

#### For Supervisors:
1. Login with supervisor credentials
2. Access "My Team" to view assigned staff
3. Click "View Report" to see detailed staff performance
4. Add comments on objectives or activities
5. Monitor staff progress

#### For Staff:
1. Register via "/register" or receive credentials from admin
2. Login and manage own objectives, tasks, and activities
3. View assigned tasks from supervisors
4. See supervisor comments on your work

### Default User
- **Username**: admin
- **Password**: admin
- **Role**: Admin

### Server
Running on: http://localhost:8080

### Testing Checklist
- [x] Multi-user database schema
- [x] Role-based access control
- [x] Staff management CRUD
- [x] Supervisor hierarchy
- [x] Comment system
- [x] Task assignment
- [x] Self-registration
- [x] Role-based menu navigation
- [x] Server running successfully

## Next Steps (Future Enhancements)

1. **Email Notifications**:
   - Notify staff when comments are added
   - Notify about task assignments

2. **Advanced Reporting**:
   - Department-wise performance
   - Trend analysis over time
   - Comparison charts

3. **File Attachments**:
   - Attach documents to objectives
   - Evidence for activity completion

4. **Task Workflows**:
   - Approval workflows for tasks
   - Task dependencies

5. **Advanced Search**:
   - Search objectives across staff
   - Filter by performance range
   - Advanced filters on tasks

6. **Audit Trail**:
   - Track who made changes
   - History of edits

7. **Performance Reviews**:
   - Scheduled review cycles
   - Rating system
   - Review templates
