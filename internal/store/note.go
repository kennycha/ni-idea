package store

type NoteType string

const (
	TypeProblem   NoteType = "problem"
	TypeDecision  NoteType = "decision"
	TypeKnowledge NoteType = "knowledge"
	TypePractice  NoteType = "practice"
)

var AllTypes = []NoteType{TypeProblem, TypeDecision, TypeKnowledge, TypePractice}

var DefaultSearchTypes = []NoteType{TypeProblem, TypeDecision}

func (t NoteType) String() string {
	return string(t)
}

func (t NoteType) IsValid() bool {
	switch t {
	case TypeProblem, TypeDecision, TypeKnowledge, TypePractice:
		return true
	}
	return false
}

func (t NoteType) Directory() string {
	switch t {
	case TypeProblem:
		return "problems"
	case TypeDecision:
		return "decisions"
	case TypeKnowledge:
		return "knowledge"
	case TypePractice:
		return "practice"
	}
	return ""
}

type NoteMeta struct {
	Title   string   `yaml:"title"`
	Type    NoteType `yaml:"type"`
	Tags    []string `yaml:"tags"`
	Private bool     `yaml:"private"`
	Created string   `yaml:"created"`
	Updated string   `yaml:"updated"`
}

type Note struct {
	Meta    NoteMeta
	Content string
	Path    string
}
