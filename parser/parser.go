package parser

import "d_timewheel/job"

type IParser interface {
	Parse() (job.Job, error)
}
