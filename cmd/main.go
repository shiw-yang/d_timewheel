package main

import (
	"d_timewheel/config"
	"d_timewheel/job"
	"d_timewheel/parser"
)

type DTWServer struct {
}

func (s *DTWServer) Start() error {
	return nil
}

func (s *DTWServer) Stop() {
}

func (s *DTWServer) registerJob(jobInfo job.Job) (job.JobDto, error) {
	return job.JobDto{}, nil
}

func (s *DTWServer) RegisterParserJob(parser parser.IParser) (job.JobDto, error) {
	jobInfo, err := parser.Parse()
	if err != nil {
		return job.JobDto{}, err
	}
	return s.registerJob(jobInfo)
}

func (s *DTWServer) StopJob(job int64) error {
	return nil
}

type IDTWServer interface {
	Start() error
	Stop()
	registerJob(job job.Job) (job.JobDto, error)
	RegisterParserJob(parser parser.IParser) (job.JobDto, error)
	StopJob(job int64) error
	QueryJob(job int64) (job.JobDto, error)
}

func InitDTWServer(conf config.Config) (*DTWServer, error) {
	return &DTWServer{}, nil
}

func main() {
	conf, err := config.InitConfig("./conf/config.toml")
	if err != nil {
		panic(err)
	}
	server, err := InitDTWServer(conf)
	if err != nil {
		panic(err)
	}
	server.Start()
	defer server.Stop()
}
