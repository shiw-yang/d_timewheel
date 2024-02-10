package parser

import "d_timewheel/job"

type QuartzJob struct{}

func (q *QuartzJob) Parse() job.Job {
	return job.Job{}
}

