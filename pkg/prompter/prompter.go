package prompter

import (
	"github.com/AlecAivazis/survey/v2"
)

type Prompter interface {
	Select(msg string, defaultValue string, options []string) (result int, err error)
}

func New() Prompter {
	return &surveyPrompter{}
}

type surveyPrompter struct{}

// Select return index of selected item, and -1 if no item selected and error occurred
func (p *surveyPrompter) Select(message, defaultValue string, options []string) (int, error) {
	const prompterPageSize = 15
	q := &survey.Select{
		Message:  message,
		Options:  options,
		PageSize: prompterPageSize,
	}

	if defaultValue != "" {
		q.Default = defaultValue
	}

	var result int

	err := survey.AskOne(q, &result)
	if err != nil {
		return -1, err
	}

	return result, nil
}
