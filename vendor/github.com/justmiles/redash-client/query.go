package redash

import (
	"fmt"
	"strconv"
	"time"
)

// Query is returned for all /api/queries
type Query struct {
	User           `json:"user,omitempty"`
	Visualizations `json:"visualizations,omitempty"`

	Schedule struct {
		Interval  float64     `json:"interval,omitempty"`
		Until     string      `json:"until,omitempty"`
		DayOfWeek string      `json:"name,omitempty"`
		Value     interface{} `json:"day_of_week,omitempty"`
		Time      time.Time   `json:"time,omitempty"`
	} `json:"parameters,omitempty"`

	Options struct {
		Parameters []struct {
			Global bool        `json:"global,omitempty"`
			Type   string      `json:"type,omitempty"`
			Name   string      `json:"name,omitempty"`
			Value  interface{} `json:"value,omitempty"`
			Title  string      `json:"title,omitempty"`
		} `json:"parameters,omitempty"`
	} `json:"options,omitempty"`

	Message           string    `json:"message,omitempty"`
	DataSourceID      int       `json:"data_source_id,omitempty"`
	LastModifiedByID  int       `json:"last_modified_by_id,omitempty"`
	LatestQueryDataID int       `json:"latest_query_data_id,omitempty"`
	IsArchived        bool      `json:"is_archived,omitempty"`
	RetrievedAt       time.Time `json:"retrieved_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
	Query             string    `json:"query,omitempty"`
	IsDraft           bool      `json:"is_draft,omitempty"`
	ID                int       `json:"id,omitempty"`
	Description       string    `json:"description,omitempty"`
	Runtime           float64   `json:"runtime,omitempty"`
	Name              string    `json:"name,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	Version           int       `json:"version,omitempty"`
	QueryHash         string    `json:"query_hash,omitempty"`
	APIKey            string    `json:"api_key,omitempty"`
	IsFavorite        bool      `json:"is_favorite,omitempty"`
	Tags              []string  `json:"tags,omitempty"`
	IsSafe            bool      `json:"is_safe,omitempty"`
}

// ListQueriesResponse represents the API response to a ListQueries request
type ListQueriesResponse struct {
	Count    int     `json:"count"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
	Results  []Query `json:"results"`
}

// ListQueries returns all Redash queries in a single slice
// see ListQueriesWithPagination for a paginated option
func (c *Client) ListQueries() (queries []Query, err error) {
	var total = 1
	var page = 1

	for len(queries) < total {
		p := ListQueriesResponse{}

		req, err := c.newRequest("GET", "/api/queries", nil, &map[string]string{
			"page":      fmt.Sprintf("%d", page),
			"page_size": fmt.Sprintf("%d", 25),
		})

		if err != nil {
			return queries, err
		}
		_, err = c.do(req, &p)
		if err != nil {
			return queries, err
		}
		total = p.Count
		page++
		queries = append(queries, p.Results...)
	}
	return queries, nil
}

// CreateQuery creates a query
func (c *Client) CreateQuery(q Query) (Query, error) {
	req, err := c.newRequest("POST", fmt.Sprintf("/api/queries"), q, nil)
	if err != nil {
		return q, err
	}
	_, err = c.do(req, &q)
	if q.Message != "" {
		return q, fmt.Errorf("Error Creating Query: %s", q.Message)
	}

	return q, err
}

func (c *Client) DeleteQuery(q Query) (err error) {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/api/queries/%d", q.ID), nil, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, nil)
	return err
}

// GetQueryByID returns Query
func (c *Client) GetQueryByID(i int) (q *Query, err error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/api/queries/%d", i), nil, nil)
	if err != nil {
		return q, err
	}
	_, err = c.do(req, &q)
	return q, err
}

// ListQueriesWithPagination returns the paginated ListQueriesResponse
func (c *Client) ListQueriesWithPagination(options *map[string]string) (p ListQueriesResponse, err error) {

	req, err := c.newRequest("GET", "/api/queries", nil, options)
	if err != nil {
		return p, err
	}
	_, err = c.do(req, &p)
	if err != nil {
		return p, err
	}
	return p, nil
}

// SearchQueries searches Redash for queries
// q - query to search
// includeDrafts - boolean whether or not to include drafts
func (c *Client) SearchQueries(q string, includeDrafts bool, maxResults int) (queries []Query, err error) {
	var (
		page  = 1
		total = 1
	)

	for len(queries) < total && len(queries) < maxResults {
		p := ListQueriesResponse{}

		req, err := c.newRequest("GET", "/api/queries/search", nil, &map[string]string{
			"page":           fmt.Sprintf("%d", page),
			"page_size":      fmt.Sprintf("%d", 25),
			"q":              q,
			"include_drafts": strconv.FormatBool(includeDrafts),
		})

		if err != nil {
			return queries, err
		}
		_, err = c.do(req, &p)
		if err != nil {
			return queries, err
		}
		total = p.Count
		page++
		queries = append(queries, p.Results...)
	}

	if len(queries) > maxResults {
		queries = queries[0:maxResults]
	}
	return queries, nil
}

// DownloadResults will write the latest query results to the filesystem
func (c *Client) DownloadResults(q Query, filelocation, filetype string) (err error) {
	if filetype == "" {
		filetype = "xlsx"
	}
	if filetype != "xlsx" && filetype != "csv" {
		return fmt.Errorf(`unable to download file of type "%s". Please specify "xlsx" or "csv"`, filetype)
	}
	req, err := c.newRequest("GET", fmt.Sprintf("/api/queries/%d/results/%d.%s", q.ID, q.LatestQueryDataID, filetype), nil, nil)
	if err != nil {
		return err
	}
	_, err = c.download(req, filelocation)
	if err != nil {
		return err
	}
	return nil
}

// RefreshQuery refreshes a query
func (c *Client) RefreshQuery(q *Query) (j Job, err error) {
	var jr JobResponse
	req, err := c.newRequest("POST", fmt.Sprintf("/api/queries/%d/refresh", q.ID), nil, nil)
	if err != nil {
		return j, err
	}

	_, err = c.do(req, &jr)
	if err != nil {
		return j, err
	}

	return jr.Job, err
}

// RefreshQueryWait execute's a query and waits polls for the response before returning
func (c *Client) RefreshQueryWait(q *Query, interval int) (*Query, error) {
	job, err := c.RefreshQuery(q)
	if err != nil {
		return q, err
	}

	err = job.Poll(c, interval)
	if err != nil {
		return q, err
	}

	updatedQuery, err := c.GetQueryByID(q.ID)
	return updatedQuery, err

}

func (q *Query) String() string {
	return fmt.Sprintf("%d\t%s\t%s\t%s", q.ID, q.Name, q.User.Name, q.RetrievedAt)
}
