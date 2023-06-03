package slacker

// JobDefinition structure contains definition of the job
type JobDefinition struct {
	Spec        string
	Description string
	Handler     JobHandler

	// HideHelp will hide this job definition from appearing in the `help` results.
	HideHelp bool
}

// newJob creates a new job object
func newJob(spec string, definition *JobDefinition) Job {
	definition.Spec = spec
	return &job{
		definition: definition,
	}
}

// Job interface
type Job interface {
	Definition() *JobDefinition
}

// job structure contains the job's spec and handler
type job struct {
	definition *JobDefinition
}

// Definition returns the job's definition
func (c *job) Definition() *JobDefinition {
	return c.definition
}
