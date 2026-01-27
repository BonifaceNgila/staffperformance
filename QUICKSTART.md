# Quick Start Guide

## Login
1. Go to http://localhost:8080
2. Login with:
   - Username: `admin`
   - Password: `admin123`

## Create Your First Objective

1. Click **"+ New Objective"** button
2. Fill in:
   - Objective Title (e.g., "Improve Customer Satisfaction")
   - Description (optional)
   - Start Date
   - End Date
3. Click **"Create Objective"**

## Add Expected Outcomes

1. From the dashboard, find your objective
2. Click **"+ Add Outcome"** under Expected Outcomes section
3. Fill in:
   - Expected Outcome Title (e.g., "Reduce response time")
   - Description (optional)
4. Click **"Add Expected Outcome"**

## Add Activities/Tasks

1. From the dashboard, find your expected outcome
2. Click **"+ Add Activity"**
3. Fill in:
   - Activity Title (e.g., "Implement automated email responses")
   - Description (optional)
   - Category (Daily/Weekly/Monthly/Quarterly/Biannually/Annually)
   - Progress Percentage (0-100%)
   - Implementation Level (describe current status)
4. Click **"Add Activity"**

## Track Progress

1. View your objectives on the dashboard
2. See the **Performance badge** showing overall progress
3. Edit activities to update progress percentages
4. The objective performance updates automatically

## Example Workflow

```
Objective: "Improve Team Productivity"
├── Expected Outcome: "Reduce project delivery time"
│   ├── Activity: "Implement daily standups" (Weekly, 75%, "Meeting daily at 9am")
│   └── Activity: "Use project management tool" (Daily, 90%, "All tasks tracked in Jira")
└── Expected Outcome: "Improve code quality"
    ├── Activity: "Conduct code reviews" (Daily, 80%, "All PRs reviewed within 24h")
    └── Activity: "Increase test coverage" (Monthly, 60%, "Currently at 65% coverage")

Overall Objective Performance = (75 + 90 + 80 + 60) / 4 = 76.25%
```

## Tips

- **Update regularly**: Keep activity progress percentages current for accurate performance tracking
- **Use categories wisely**: Choose frequency categories that match how often you'll review each activity
- **Be specific**: Detailed implementation levels help track what's actually been done
- **Start small**: Begin with a few objectives and expand as you get comfortable with the system

## Common Actions

### Edit an Item
Click the **"Edit"** button next to any objective, outcome, or activity

### Delete an Item
Click the **"Delete"** button (note: deleting an objective deletes all associated outcomes and activities)

### View Performance
Check the green performance badge on each objective card

### Logout
Click **"Logout"** in the top navigation bar
