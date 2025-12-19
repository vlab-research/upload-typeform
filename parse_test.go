package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/vlab-research/trans"
	"testing"
)

func TestParseCsv(t *testing.T) {
	assert.Equal(t, true, true)
}

func TestBuildField_GetsRefAndDescription(t *testing.T) {
	i, _ := BuildField([]string{"ref", "multiple_choice", "foo", "yes\nno", "description"})
	f := i.(*trans.Field)
	assert.Equal(t, "ref", f.Ref)
	assert.Equal(t, "description", f.Properties.Description)
}

func TestBuildField_ErrorsWhenMultipleChoiceHasNoAnswers(t *testing.T) {
	_, e := BuildField([]string{"ref", "multiple_choice", "foo", "", "description"})
	assert.NotNil(t, e)
}

func TestBuildField_GetsTitleFromOpenQuestions(t *testing.T) {
	i, _ := BuildField([]string{"ref", "short_text", "foo", "", ""})
	f := i.(*trans.Field)
	assert.Equal(t, "foo", f.Title)
}

func TestBuildField_GetsThankyouScreen(t *testing.T) {
	i, _ := BuildField([]string{"ref", "thankyou_screen", "foo", "", ""})
	ty := i.(*ThankyouScreen)
	assert.Equal(t, "foo", ty.Title)
}

func TestBuildField_GetsHiddenVariable(t *testing.T) {
	i, _ := BuildField([]string{"ref", "hidden", "foo", "", ""})
	v := i.(HiddenVariable)
	assert.Equal(t, HiddenVariable("ref"), v)
}

func TestBuildField_GetsTitleFromMultipleChoiceQuestion(t *testing.T) {
	i, _ := BuildField([]string{"ref", "multiple_choice", "foo", "A. yes\nB. no", ""})
	f := i.(*trans.Field)
	assert.Equal(t, "foo\n\nA. yes\nB. no", f.Title)

	i, _ = BuildField([]string{"ref", "multiple_choice", "foo\n", "A. yes\nB. no", ""})
	f = i.(*trans.Field)
	assert.Equal(t, "foo\n\nA. yes\nB. no", f.Title)

	i, _ = BuildField([]string{"ref", "multiple_choice", "foo", "yes\nno", ""})
	f = i.(*trans.Field)
	assert.Equal(t, "foo", f.Title)
}

func TestBuildField_GetsChoicesFromMultipleChoiceQuestionWithLabels(t *testing.T) {
	i, _ := BuildField([]string{"ref", "multiple_choice", "foo", "A. yes\nB. no", ""})
	f := i.(*trans.Field)
	assert.Equal(t, f.Properties.Choices, []*trans.FieldChoice{{"", "A", ""}, {"", "B", ""}})
}

func TestBuildField_GetsChoicesFromMultipleChoiceQuestionWithoutLabels(t *testing.T) {
	i, _ := BuildField([]string{"ref", "multiple_choice", "foo", "yes\nno", ""})
	f := i.(*trans.Field)
	assert.Equal(t, f.Properties.Choices, []*trans.FieldChoice{{"", "yes", ""}, {"", "no", ""}})

	i, _ = BuildField([]string{"ref", "multiple_choice", "foo", "\nyes\nno", ""})
	f = i.(*trans.Field)
	assert.Equal(t, f.Properties.Choices, []*trans.FieldChoice{{"", "yes", ""}, {"", "no", ""}})
}

func TestBuildField_TrimsSpacesFromChoicesWithoutLabels(t *testing.T) {
	// Leading spaces
	i, _ := BuildField([]string{"ref", "multiple_choice", "foo", "  yes\n  no", ""})
	f := i.(*trans.Field)
	assert.Equal(t, []*trans.FieldChoice{{"", "yes", ""}, {"", "no", ""}}, f.Properties.Choices)

	// Trailing spaces
	i, _ = BuildField([]string{"ref", "multiple_choice", "foo", "yes  \nno  ", ""})
	f = i.(*trans.Field)
	assert.Equal(t, []*trans.FieldChoice{{"", "yes", ""}, {"", "no", ""}}, f.Properties.Choices)

	// Both leading and trailing spaces
	i, _ = BuildField([]string{"ref", "multiple_choice", "foo", "  yes  \n  no  ", ""})
	f = i.(*trans.Field)
	assert.Equal(t, []*trans.FieldChoice{{"", "yes", ""}, {"", "no", ""}}, f.Properties.Choices)

	// Preserves internal spaces
	i, _ = BuildField([]string{"ref", "multiple_choice", "foo", "  hello world  \n  foo bar  ", ""})
	f = i.(*trans.Field)
	assert.Equal(t, []*trans.FieldChoice{{"", "hello world", ""}, {"", "foo bar", ""}}, f.Properties.Choices)

	// Skips whitespace-only lines
	i, _ = BuildField([]string{"ref", "multiple_choice", "foo", "yes\n   \nno", ""})
	f = i.(*trans.Field)
	assert.Equal(t, []*trans.FieldChoice{{"", "yes", ""}, {"", "no", ""}}, f.Properties.Choices)
}

func TestBuildField_GetsChoicesFromMultipleChoiceQuestionSkippingLetters(t *testing.T) {
	i, _ := BuildField([]string{"ref", "multiple_choice", "foo", "A. yes\nC. no", ""})
	f := i.(*trans.Field)
	assert.Equal(t, f.Properties.Choices, []*trans.FieldChoice{{"", "A", ""}, {"", "C", ""}})
}

func TestBuildForm_IgnoresBlankLines(t *testing.T) {
	records := [][]string{
		{"", "", "", ""},
		{"", "", "", ""},
		{"ref", "foo", "A. yes\nC. no", ""},
	}

	form, err := BuildForm("foo", records)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(form.Fields))
}

func TestParseMessages_SkipsEmpty(t *testing.T) {
	records := [][]string{
		{"variable", "message"},
		{"", ""},
		{"foo.bar", "baz"},
	}

	m := ParseMessages(records)
	assert.Equal(t, 1, len(m))
	assert.Equal(t, "baz", m["foo.bar"])
}
