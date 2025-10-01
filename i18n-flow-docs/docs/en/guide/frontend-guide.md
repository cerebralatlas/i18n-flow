# Frontend Usage Guide

The i18n Flow admin interface is a modern frontend application developed with React, providing a user-friendly interface to manage your internationalization content. This guide will help you understand how to use the various features of the admin interface.

## Login

When you first access the i18n Flow admin interface, you need to log in:

1. Visit the admin interface (default is <http://localhost> or your production deployment address)
2. Enter your username and password
   - Default admin account: `admin`
   - Default password: `admin123`
3. Click the "Login" button

::: warning Note
After your first login, immediately change the default password to ensure system security!
:::

## Dashboard Overview

After successful login, you will see the main dashboard, which provides overall statistics about the system:

- Total number of projects
- Total number of languages
- Total number of translation keys
- Translation completion rate
- Recent activity

## Project Management

### Creating a New Project

1. Click "Projects" in the sidebar menu
2. Click the "Create Project" button
3. Fill in the project information:
   - Project name (required)
   - Project description (optional)
   - Project identifier (slug, used for API and CLI calls)
4. Click the "Save" button

### Viewing Project List

The project list page displays all projects and their basic information:

- Project name
- Creation date
- Number of languages
- Number of translation keys
- Completion rate

### Editing a Project

1. Find the project you want to edit in the project list
2. Click the "Edit" button
3. Modify the project information
4. Click the "Save" button

### Deleting a Project

1. Find the project you want to delete in the project list
2. Click the "Delete" button
3. Click "Confirm" in the confirmation dialog

::: danger Warning
Deleting a project will permanently remove all translation data for that project. This action cannot be undone!
:::

## Project Member Management

### Understanding Roles and Permissions

i18n Flow implements a comprehensive role-based access control system:

**Project Roles:**
- **Viewer**: Can view projects, translations, and export data
- **Editor**: Viewer permissions + can create, update, and delete translations
- **Owner**: Editor permissions + can manage project settings and members

**System Roles:**
- **Admin**: Full system access, can manage users, languages, and all projects
- **User**: Standard user with project-specific permissions

### Adding Project Members

1. Navigate to your project in the project list
2. Click on the project name to open project details
3. Switch to the "Members" tab
4. Click the "Add Member" button
5. Enter the user's email or username
6. Select the appropriate role from the dropdown
7. Click "Save" to send the invitation

### Managing Existing Members

- **Update Role**: Click the "Edit" button next to a member to change their role
- **Remove Member**: Click the "Remove" button to remove a member from the project
- **View Member Details**: Click on a member's name to see their profile and permissions

### Viewing Accessible Projects

Users can easily see all projects they have access to:

1. Click "Projects" in the sidebar menu
2. The project list will show both projects you own and projects you've been invited to
3. Use the filter options to view projects by role (All, Owner, Editor, Viewer)

### Best Practices for Team Collaboration

1. **Assign roles based on responsibilities**: 
   - Viewers: Stakeholders who only need to review translations
   - Editors: Translators and developers who work with translations
   - Owners: Project managers who need full control

2. **Regular access reviews**: Periodically review team member access
3. **Use descriptive project names**: Help team members quickly identify relevant projects
4. **Communicate changes**: Notify team members when their roles change

## Language Management

### Adding a Language

1. Click "Languages" in the sidebar menu
2. Click the "Add Language" button
3. Fill in the language information:
   - Language name (e.g., English)
   - Language code (e.g., en)
   - Language locale (e.g., en-US)
   - RTL support (right-to-left text direction)
4. Click the "Save" button

### Editing a Language

1. Find the language you want to edit in the language list
2. Click the "Edit" button
3. Modify the language information
4. Click the "Save" button

## Translation Management

### Translation Interface

1. Click "Projects" in the sidebar menu
2. Click on the name of the project whose translations you want to manage
3. Switch to the "Translations" tab on the project details page

The translation interface displays the translation content in a table format, with each row representing a translation key and each column representing a language.

### Adding a Translation Key

1. Click the "Add Key" button in the translation interface
2. Fill in the following information:
   - Translation key (e.g., common.buttons.save)
   - Description (optional, helps translators understand the context)
   - Translation text for the default language
3. Click the "Save" button

### Editing Translations

1. Find the translation key you want to edit in the translation table
2. Click on the cell for the corresponding language column
3. Enter the translation text
4. Click outside the cell or press Enter to save

### Bulk Import of Translations

1. Click the "Import" button in the translation interface
2. Select the import format (JSON, Excel, CSV)
3. Upload a file or paste content
4. Select import options (merge or overwrite)
5. Click the "Import" button

### Exporting Translations

1. Click the "Export" button in the translation interface
2. Select the export format (JSON, Excel, CSV)
3. Select the languages to export
4. Click the "Export" button

## User Management

### Understanding User Roles

i18n Flow supports different user roles with varying levels of system access:

**Admin Users:**
- Full system access
- Can manage all projects, languages, and users
- Can access system settings and configuration

**Regular Users:**
- Can access projects they've been invited to
- Can create and manage their own projects
- Can be assigned specific roles within projects (viewer, editor, owner)

### Creating a User

1. Click "Users" in the sidebar menu
2. Click the "Create User" button
3. Fill in the user information:
   - Username (required)
   - Email (required)
   - Password (required)
   - Role (admin or regular user)
   - Status (active/inactive)
4. Click the "Save" button

### Editing a User

1. Find the user you want to edit in the user list
2. Click the "Edit" button
3. Modify the user information:
   - Update basic information (username, email)
   - Change user role if needed
   - Reset password if required
   - Update user status
4. Click the "Save" button

### Managing User Access

- **Reset Password**: Click "Reset Password" to generate a new temporary password
- **Deactivate User**: Change user status to inactive to temporarily disable access
- **View User Projects**: See all projects a user has access to and their roles
- **Delete User**: Remove a user from the system (requires admin confirmation)

### User Profile Management

Users can manage their own profiles:

1. Click on the user avatar in the top right corner
2. Select "Profile" from the dropdown
3. Update personal information:
   - Change password
   - Update email address
   - Modify display preferences
4. Click "Save" to apply changes

## Settings

### Profile Settings

1. Click on your user avatar in the top right corner
2. Select "Profile"
3. Modify your personal information
4. Click the "Save" button

### Changing Password

1. Click on your user avatar in the top right corner
2. Select "Change Password"
3. Enter your current password and new password
4. Click the "Save" button

## Common Issues

### Unable to Login

- Confirm that your username and password are correct
- Check if the backend service is running properly
- Clear your browser cache and try again

### Unable to Save Translations

- Confirm that you have permission to edit the project
- Check your network connection
- Look for error messages in the browser console
