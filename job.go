package slacker

// JobDefinition structure contains definition of the job
type JobDefinition struct {
	CronExpression string
	Name           string
	Description    string
	Middlewares    []JobMiddlewareHandler
	Handler        JobHandler

	// HideHelp will hide this job definition from appearing in the `help` results.
	HideHelp bool
}

// newJob creates a new job object
func newJob(definition *JobDefinition) *Job {
	return &Job{
		definition: definition,
	}
}

// Job structure contains the job's spec and handler
type Job struct {
	definition *JobDefinition
}

// Definition returns the job's definition
func (c *Job) Definition() *JobDefinition {
	return c.definition
}
