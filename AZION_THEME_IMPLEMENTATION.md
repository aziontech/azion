# Azion Custom Theme Implementation

## Overview
Created a custom theme for the Azion CLI `init` command using the Huh library's theme system. The theme uses Azion's brand colors to provide a consistent and branded user experience.

## Color Scheme
- **Primary (Purple)**: `#b5b1f4` - Used for titles, labels, and category displays
- **Secondary (Orange)**: `#f3652b` - Used for success messages, checkmarks, and interactive elements

## Files Modified/Created

### 1. `/pkg/cmd/init/theme.go` (NEW)
Created a new theme file containing:

- **`ThemeAzion()`**: Main theme function that returns a customized `huh.Theme`
  - Based on `ThemeBase()` from the Huh library
  - Applies Azion brand colors to all interactive elements
  - Configures focused and blurred states
  - Sets up form, group, and field styles

- **Helper Functions**:
  - `GetAzionLabelStyle()`: Returns purple style for labels
  - `GetAzionAnswerStyle()`: Returns orange style for user answers
  - `GetAzionSuccessStyle()`: Returns orange bold style for success messages
  - `GetAzionTitleStyle()`: Returns purple bold style for titles

### 2. `/pkg/cmd/init/init.go` (MODIFIED)
Updated to use the new Azion theme:

- **Lines 213-214**: Replaced hardcoded color styles with theme helper functions
- **Line 221**: Added `.WithTheme(ThemeAzion())` to category selection
- **Line 235**: Added `.WithTheme(ThemeAzion())` to framework selection
- **Line 283**: Added `.WithTheme(ThemeAzion())` to template selection
- **Line 297**: Added `.WithTheme(ThemeAzion())` to project name input
- **Line 324**: Added `.WithTheme(ThemeAzion())` to confirmation prompt
- **Line 448**: Updated success style to use `GetAzionSuccessStyle()`
- **Line 572**: Added `.WithTheme(ThemeAzion())` to confirmWithHuh form
- **Line 588**: Updated printSuccess to use `GetAzionSuccessStyle()`
- **Line 598**: Updated printTitle to use `GetAzionTitleStyle()`
- **Line 17**: Removed unused `lipgloss` import (now handled in theme.go)

### 3. `/pkg/cmd/init/utils.go` (MODIFIED)
Applied theme to utility forms:

- **Line 26**: Added `.WithTheme(ThemeAzion())` to askForInput form
- **Line 226**: Added `.WithTheme(ThemeAzion())` to confirmProceed form

## Visual Output Example

```
Category: JavaScript                    (purple label, orange value)
Template: Hello World                   (purple label, orange value)
Project Name: brave-villain             (purple label, orange value)

Creating your project...                (purple title)
  ✓ Template downloaded                 (orange checkmark)
  ✓ Files extracted                     (orange checkmark)
  ✓ Configuration generated             (orange checkmark)

Template successfully configured        (orange bold text)
```

## Theme Features

### Interactive Elements
- **Select menus**: Orange selector (">") and indicators
- **Input fields**: Orange cursor and prompt
- **Buttons**: Orange background when focused
- **Multi-select**: Orange checkmarks for selected items

### Text Elements
- **Titles**: Purple, bold
- **Descriptions**: Subtle gray
- **Options**: Normal foreground color
- **Error messages**: Red (preserved from base theme)

### States
- **Focused**: Full color scheme with visible borders
- **Blurred**: Hidden borders, maintains readability

## Benefits

1. **Brand Consistency**: Uses official Azion colors throughout the CLI
2. **Better UX**: Clear visual hierarchy with distinct colors for different element types
3. **Maintainability**: Centralized theme in one file, easy to update
4. **Reusability**: Helper functions can be used anywhere in the init package
5. **Accessibility**: Maintains good contrast ratios for readability

## Testing

To test the theme:
```bash
cd /Users/vitor.eltz/git/azioncli/azion
go build
./azion init
```

The interactive prompts will now display with the new Azion color scheme.
