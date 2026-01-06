package init

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// ThemeAzion returns a custom Azion theme with brand colors
// Primary color: #b5b1f4 (light purple) for titles and labels
// Secondary color: #f3652b (orange) for success messages and interactive elements
func ThemeAzion() *huh.Theme {
	t := huh.ThemeBase()

	var (
		normalFg      = lipgloss.AdaptiveColor{Light: "235", Dark: "252"}
		azionPurple   = lipgloss.Color("#b5b1f4") // Light purple for titles
		azionOrange   = lipgloss.Color("#f3652b") // Orange for success/interactive
		cream         = lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}
		red           = lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#ED567A"}
		subtleGray    = lipgloss.AdaptiveColor{Light: "", Dark: "243"}
	)

	// Focused state styles
	t.Focused.Base = t.Focused.Base.BorderForeground(lipgloss.Color("238"))
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = t.Focused.Title.Foreground(azionPurple).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(azionPurple).Bold(true).MarginBottom(1)
	t.Focused.Directory = t.Focused.Directory.Foreground(azionPurple)
	t.Focused.Description = t.Focused.Description.Foreground(subtleGray)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(red)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(red)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(azionOrange)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(azionOrange)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(azionOrange)
	t.Focused.Option = t.Focused.Option.Foreground(normalFg)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(azionOrange)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(azionOrange)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(azionOrange).SetString("✓ ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(subtleGray).SetString("• ")
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(normalFg)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(cream).Background(azionOrange)
	t.Focused.Next = t.Focused.FocusedButton
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(normalFg).Background(lipgloss.AdaptiveColor{Light: "252", Dark: "237"})

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(azionOrange)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(lipgloss.AdaptiveColor{Light: "248", Dark: "238"})
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(azionOrange)

	// Blurred state styles
	t.Blurred = t.Focused
	t.Blurred.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	// Group styles
	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description

	// Help styles
	t.Help = help.New().Styles

	return t
}

// GetAzionLabelStyle returns the lipgloss style for labels (purple)
func GetAzionLabelStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#b5b1f4")).Bold(true)
}

// GetAzionAnswerStyle returns the lipgloss style for answers (orange)
func GetAzionAnswerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#f3652b"))
}

// GetAzionSuccessStyle returns the lipgloss style for success messages (orange)
func GetAzionSuccessStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#f3652b")).Bold(true)
}

// GetAzionTitleStyle returns the lipgloss style for titles (purple)
func GetAzionTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#b5b1f4"))
}
