package redash

import (
	"fmt"
	"time"

	prettyTime "github.com/andanhm/go-prettytime"
)

// JobResponse represents responses from the API for Job objects
type JobResponse struct {
	Job `json:"job"`
}

// Job is a Redash Job
type Job struct {
	Status        int    `json:"status"`
	Error         string `json:"error"`
	ID            string `json:"id"`
	QueryResultID int    `json:"query_result_id"`
	UpdatedAt     int    `json:"updated_at"`
}

// Poll a job until it is complete
func (j *Job) Poll(client *Client, i int) error {
	now := time.Now()
	for {
		j, err := client.GetJobByID(j.ID)
		if err != nil {
			return err
		}
		if j.Error != "" {
			return fmt.Errorf("%s", j.Error)
		}

		if j.Status == 2 && client.DebugEnabled {
			fmt.Printf("Job is running. Started %s\n", prettyTime.Format(now))
		}

		if j.Status == 3 {
			fmt.Println("Job is has finished")
			return nil
		}
		time.Sleep(time.Second * time.Duration(i))
	}
}

// GetJobByID returns a job given a job id
func (c *Client) GetJobByID(id string) (j *Job, err error) {
	var jr JobResponse
	req, err := c.newRequest("GET", fmt.Sprintf("/api/jobs/%s", id), nil, nil)
	if err != nil {
		return j, err
	}
	_, err = c.do(req, &jr)
	return &jr.Job, err
}
