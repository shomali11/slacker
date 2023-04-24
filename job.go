package slacker

// JobDefinition structure contains definition of the job
type JobDefinition struct {
	Description string
	Handler     JobHandler

	// HideHelp will hide this job definition from appearing in the `help` results.
	HideHelp bool
}

// newJob creates a new job object
func newJob(spec string, definition *JobDefinition) Job {
	return &job{
		spec:       spec,
		definition: definition,
	}
}

// Job interface
type Job interface {
	Spec() string
	Definition() *JobDefinition
	Callback(JobContext) func()
}

// job structure contains the job's spec and handler
type job struct {
	spec       string
	definition *JobDefinition
}

// Spec returns the job's spec
func (c *job) Spec() string {
	return c.spec
}

// Definition returns the job's definition
func (c *job) Definition() *JobDefinition {
	return c.definition
}

// Callback returns cron job callback
func (c *job) Callback(jobCtx JobContext) func() {
	return func() {
		c.Definition().Handler(jobCtx)
	}
}
