# Staff Performance System - User Guide

## Overview

The Staff Performance System now features a **modular dashboard** with three main sections accessible via a visual menu:

1. **Tasks** - Standalone task management
2. **Objectives** - Performance objectives with expected outcomes and activities
3. **Reports** - Comprehensive performance analytics

## Main Dashboard

After logging in, you'll see the home dashboard with:

- **Statistics Cards**: Quick overview of your objectives, tasks, and performance
- **Main Menu**: Four clickable tiles for navigation
  - 🏠 **Home** - Dashboard overview
  - ✓ **Tasks** - Task management
  - 🎯 **Objectives** - Performance objectives
  - 📊 **Reports** - Performance reports
- **Quick Actions**: Shortcuts to create new items

## Tasks Section

### Features
- Create, edit, and delete standalone tasks
- Filter tasks by status (All, Pending, In Progress, Completed, On Hold)
- Track task priority levels
- Monitor due dates

### Task Properties
- **Title** - Task name
- **Description** - Detailed information
- **Priority** - Low, Medium, High, Urgent
- **Status** - Pending, In Progress, Completed, On Hold
- **Due Date** - Deadline for the task

### How to Use Tasks

1. **Create a Task**
   - Click "Tasks" from the main menu
   - Click "+ New Task" button
   - Fill in task details
   - Select priority and status
   - Set due date
   - Click "Create Task"

2. **View Tasks**
   - Tasks are displayed as cards with color-coded badges
   - Filter by status using the filter buttons
   - Each card shows title, description, priority, status, and dates

3. **Edit a Task**
   - Click "Edit" button on any task card
   - Modify the details
   - Click "Update Task"

4. **Complete a Task**
   - Edit the task
   - Change status to "Completed"
   - The completion date is automatically recorded

5. **Delete a Task**
   - Click "Delete" button on the task card
   - Confirm deletion

### Task Viewer Features
- **Visual Cards**: Each task displayed in an easy-to-read card format
- **Color Coding**: 
  - Priority badges (blue, orange, red, purple)
  - Status badges (yellow, blue, green, pink)
- **Quick Filtering**: One-click filtering by status
- **Meta Information**: Created date, due date, completion date

## Objectives Section

### Structure
```
Objective
├── Expected Outcome 1
│   ├── Activity 1
│   ├── Activity 2
│   └── Activity 3
└── Expected Outcome 2
    ├── Activity 1
    └── Activity 2
```

### How to Use Objectives

1. **Create an Objective**
   - Click "Objectives" from the main menu
   - Click "+ New Objective"
   - Enter title, description, start and end dates
   - Click "Create Objective"

2. **Add Expected Outcomes**
   - Find your objective on the objectives page
   - Click "+ Add Outcome" under the objective
   - Enter outcome title and description
   - Click "Add Expected Outcome"

3. **Add Activities to Outcomes**
   - Under an expected outcome, click "+ Add Activity"
   - Fill in:
     - Activity title and description
     - Category (Daily, Weekly, Monthly, Quarterly, Biannually, Annually)
     - Progress percentage (0-100%)
     - Implementation level (current status description)
   - Click "Add Activity"

4. **Track Progress**
   - Edit activities to update progress percentage
   - The objective performance automatically calculates as the mean of all activities
   - Visual progress bars show completion status

## Reports Section

### What's Included

1. **Overview Summary**
   - Total objectives count
   - Completed and pending tasks
   - Average performance across all objectives

2. **Objectives Performance**
   - Detailed breakdown of each objective
   - Performance percentage
   - All expected outcomes and activities
   - Progress visualization

3. **Tasks Summary**
   - Complete list of all tasks
   - Priority and status for each
   - Due dates

4. **Actions**
   - Print report (browser print function)
   - Export data

### How to Use Reports

1. Click "Reports" from the main menu
2. Review the summary statistics
3. Scroll through objectives and tasks
4. Click "🖨️ Print Report" to generate a printable version
5. Use your browser's print dialog to save as PDF or print

## Navigation

### Main Menu
- Always visible at the top of each section
- Click any menu item to switch sections
- Active section is highlighted

### Breadcrumbs
- Forms show breadcrumb navigation
- Example: "Objective: Improve Sales" → "Expected Outcome: Increase Leads"

### Back Navigation
- "Back to Dashboard" links throughout the system
- Browser back button also works

## Performance Calculation

### Objective Performance Formula
```
Objective Performance = (Sum of all activity progress %) / (Number of activities)
```

### Example
If an objective has 4 activities with progress: 75%, 90%, 80%, 60%
- Total: 75 + 90 + 80 + 60 = 305
- Count: 4 activities
- Performance: 305 / 4 = **76.25%**

### Average Performance
The dashboard shows average performance across all objectives:
```
Average Performance = (Sum of all objective performances) / (Number of objectives)
```

## Tips & Best Practices

### For Tasks
- Set realistic due dates
- Use "Urgent" priority sparingly for truly critical items
- Update status regularly to keep the system current
- Use "On Hold" for blocked tasks
- Add detailed descriptions for complex tasks

### For Objectives
- Create SMART objectives (Specific, Measurable, Achievable, Relevant, Time-bound)
- Break down large objectives into multiple expected outcomes
- Create activities that are trackable and measurable
- Update progress percentages weekly or bi-weekly
- Use implementation level to document what's been completed

### For Reports
- Generate reports regularly (monthly/quarterly)
- Use print/PDF function to save historical snapshots
- Share reports with supervisors or team members
- Use performance data to identify areas needing attention

## Keyboard Shortcuts & Quick Tips

- **Tab Navigation**: Use Tab key to navigate through forms
- **Filter Toggle**: Click filter buttons to show/hide task categories
- **Quick Add**: Use quick actions on home dashboard for faster creation
- **Auto-Save**: All changes are immediately saved to the database

## Workflow Examples

### Weekly Review Workflow
1. Open **Tasks** section
2. Filter by "Pending" and "In Progress"
3. Update task statuses as needed
4. Switch to **Objectives** section
5. Update activity progress percentages
6. Check objective performance indicators
7. Go to **Reports** to see overall progress

### Monthly Planning Workflow
1. Go to **Objectives** section
2. Create new monthly objective
3. Add expected outcomes (2-4 per objective)
4. Add activities with monthly/weekly categories
5. Go to **Tasks** section
6. Create supporting tasks with due dates throughout the month
7. Set priorities appropriately

### Performance Review Workflow
1. Open **Reports** section
2. Review average performance metric
3. Identify low-performing objectives
4. Check which activities need attention
5. Review completed vs pending tasks
6. Print/save report for records
7. Create action items in **Tasks** section

## Troubleshooting

### Tasks not showing
- Check the filter buttons - you might have a filter active
- Refresh the page
- Verify you created tasks (not activities)

### Performance shows 0%
- Ensure you've created activities under expected outcomes
- Check that activities have progress percentages set
- Activities must be linked to expected outcomes, which link to objectives

### Can't delete objective
- You can only delete your own objectives
- Deleting an objective will delete all associated outcomes and activities
- This action cannot be undone

## System Requirements

- Modern web browser (Chrome, Firefox, Edge, Safari)
- JavaScript enabled
- Internet connection not required (runs locally)
- Recommended screen resolution: 1024x768 or higher

## Data Management

- All data stored locally in SQLite database
- Database file: `staffperformance.db`
- Automatic timestamps for all records
- Cascading deletes (deleting parent removes children)

---

**Need Help?** Contact your system administrator or refer to the README.md file for technical details.
