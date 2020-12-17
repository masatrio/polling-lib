package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gitlab.com/ruangguru/polling/poller"
)

type agent struct {
	httpClient        *http.Client
	httpReq           *http.Request
	timeout, interval time.Duration
	ctx               context.Context
}

type Response struct {
	Resp *http.Response
	Err  error
}

func Agent(
	httpReq *http.Request,
	interval, timeout int64,
) poller.Poller {

	if timeout == 0 {
		timeout = 30
	}

	if interval == 0 {
		interval = 5
	}

	return &agent{
		httpClient: &http.Client{},
		httpReq:    httpReq,
		interval:   interval,
		timeout:    timeout,
	}
}

func (s *agent) Run(ctx context.Context, respCh chan interface{}) error {
	if s.ctx != nil {
		return fmt.Errorf("Job already running")
	}

	if ctx == nil {
		return fmt.Errorf("Cannot use nil context")
	}

	s.ctx = ctx

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		s.httpClient.Timeout = time.Duration(s.timeout) * time.Second

		res, err := s.httpClient.Do(s.httpReq)

		respCh <- &Response{
			Resp: res,
			Err:  err,
		}

		time.Sleep(time.Duration(s.interval) * time.Second)
	}
}

func (s *agent) Stop() {

	s.ctx.Done()

	s.ctx = nil

}

func (s *agent) SetTimeout(timeout int64) {
	s.timeout = timeout
}
