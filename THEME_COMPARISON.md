# Azion Theme - Before vs After

## Color Changes

### Before (Default Theme)
- **Labels**: Color "99" (generic purple/blue)
- **Answers**: Color "86" (generic cyan/green)
- **Titles**: Color "212" (generic pink/orange)
- **Success**: Color "42" (generic green)

### After (Azion Theme)
- **Labels**: `#b5b1f4` (Azion brand purple) ✨
- **Answers**: `#f3652b` (Azion brand orange) 🔥
- **Titles**: `#b5b1f4` (Azion brand purple) ✨
- **Success**: `#f3652b` (Azion brand orange) 🔥

## Visual Example

### Before
```
Category: JavaScript                    [generic purple: generic cyan]
Template: Hello World                   [generic purple: generic cyan]
Project Name: brave-villain             [generic purple: generic cyan]

Creating your project...                [generic pink]
  ✓ Template downloaded                 [generic green]
  ✓ Files extracted                     [generic green]
  ✓ Configuration generated             [generic green]

Template successfully configured        [generic green]
```

### After (Azion Branded)
```
Category: JavaScript                    [#b5b1f4: #f3652b]
Template: Hello World                   [#b5b1f4: #f3652b]
Project Name: brave-villain             [#b5b1f4: #f3652b]

Creating your project...                [#b5b1f4]
  ✓ Template downloaded                 [#f3652b]
  ✓ Files extracted                     [#f3652b]
  ✓ Configuration generated             [#f3652b]

Template successfully configured        [#f3652b]
```

## Interactive Elements Theme

### Select Menus
- **Selector (">")**: Orange `#f3652b`
- **Selected option**: Orange `#f3652b`
- **Unselected option**: Normal foreground
- **Navigation arrows**: Orange `#f3652b`

### Input Fields
- **Prompt**: Orange `#f3652b`
- **Cursor**: Orange `#f3652b`
- **Placeholder**: Subtle gray
- **User text**: Normal foreground

### Buttons
- **Focused**: Orange background `#f3652b` with cream text
- **Blurred**: Gray background with normal text

### Confirm Dialogs
- **Title**: Purple `#b5b1f4`, bold
- **Yes/No buttons**: Orange theme when focused

## Theme Structure

The theme is based on the Huh library's `ThemeBase()` and customized with:

1. **Azion Purple** (`#b5b1f4`) for:
   - All titles and headings
   - Directory names
   - Labels in output

2. **Azion Orange** (`#f3652b`) for:
   - All interactive selectors
   - Success indicators (✓)
   - User input prompts
   - Focused buttons
   - Selected options
   - Navigation indicators

3. **Preserved from Base**:
   - Error indicators: Red
   - Subtle text: Gray
   - Normal text: Default foreground

## Implementation Details

### Theme Application Points

1. **Category Selection** - Line 221 in init.go
2. **Framework Selection** - Line 235 in init.go
3. **Template Selection** - Line 283 in init.go
4. **Project Name Input** - Line 297 in init.go
5. **Confirmation Prompt** - Line 324 in init.go
6. **Dependency Confirmations** - Line 572 in init.go
7. **Utility Forms** - Lines 26, 226 in utils.go

### Style Helper Functions

All located in `theme.go`:
- `ThemeAzion()` - Main theme constructor
- `GetAzionLabelStyle()` - For static labels
- `GetAzionAnswerStyle()` - For user answers
- `GetAzionSuccessStyle()` - For success messages
- `GetAzionTitleStyle()` - For section titles

## Benefits

✅ **Consistent Branding**: All colors match Azion's brand identity
✅ **Better Visual Hierarchy**: Purple for structure, Orange for action/success
✅ **Improved UX**: Clear distinction between labels and values
✅ **Maintainable**: Single source of truth for all theme colors
✅ **Extensible**: Easy to add more themed components
