package sdimport

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const keyHeader = "Ocp-Apim-Subscription-Key"

type Client struct {
	HTTP      *http.Client
	BaseURL   string
	apiKey    string
	keyHeader string
	Params    Parameters
}

type Parameters struct {
	Vendor              string
	PageSize            int
	FromDate            string
	ToDate              string
	LastChangeTimestamp string
	Originator          string
}

func NewParameters(from, to string) Parameters {
	return Parameters{
		Vendor:              "sarepta",
		PageSize:            25,
		FromDate:            from,
		ToDate:              to,
		LastChangeTimestamp: "2023-01-01T00:00:00",
		Originator:          "SN",
	}
}

func NewClient(dev bool, params Parameters) *Client {
	url, key := setEnv(dev)
	return &Client{
		HTTP:      &http.Client{Timeout: 15 * time.Second},
		BaseURL:   url,
		apiKey:    key,
		keyHeader: keyHeader,
		Params:    params,
	}
}

func setEnv(dev bool) (string, string) {
	if dev {
		return os.Getenv("IMPORT_TEST_URL"), os.Getenv("IMPORT_TEST_API_KEY")
	} else {
		return os.Getenv("IMPORT_PROD_URL"), os.Getenv("IMPORT_PROD_API_KEY")
	}
}

func (c *Client) GetProject(urn string) (Project, error) {
	endpoint := fmt.Sprintf("/projects/%s", urn)
	resp, err := c.getWrapper(endpoint)
	if err != nil {
		return Project{}, err
	}
	defer resp.Body.Close()
	respDump, err := dumpResponse(resp)
	if err != nil {
		return Project{}, err
	}

	if len(respDump.Project) == 0 {
		return Project{}, fmt.Errorf("no projects in response")
	}
	proj := respDump.Project[0]

	return proj, nil
}

func (c *Client) GetPublication(urn string) (Publication, error) {
	endpoint := fmt.Sprintf("/publications/%s", urn)
	resp, err := c.getWrapper(endpoint)
	if err != nil {
		return Publication{}, err
	}

	defer resp.Body.Close()
	respDump, err := dumpResponse(resp)
	if err != nil {
		return Publication{}, err
	}

	if len(respDump.Publication) == 0 {
		return Publication{}, fmt.Errorf("no projects in response")
	}
	Publication := respDump.Publication[0]

	return Publication, nil
}

func (c *Client) GetMultipleWrapper(pubType string) ([]Response, error) {
	var all []Response
	path := c.Params.buildRequestString(pubType, 0)
	log.Println("fetching from ", path)

	first, err := c.getWrapper(path)
	if err != nil {
		return nil, fmt.Errorf("error getting first page %d: %w", 0, err)
	}
	totalRecords, err := strconv.Atoi(first.Header.Get("totalrecords"))
	if err != nil {
		return nil, fmt.Errorf("total records not a number: %w", err)
	}
	pages := totalRecords / c.Params.PageSize //pages are zero indexed

	for page := 0; page <= pages; page++ {
		var resp Response

		if page == 0 {
			resp, err = dumpResponse(first)
			if err != nil {
				return nil, fmt.Errorf("error dumping response for page %d: %w", page, err)
			}
		} else {
			path := c.Params.buildRequestString("projects", page)
			next, err := c.getWrapper(path)
			if err != nil {
				return nil, fmt.Errorf("error getting next response for page %d: %w", page, err)
			}
			resp, err = dumpResponse(next)
			if err != nil {
				return nil, fmt.Errorf("error dumping response for page %d: %w", page, err)
			}
		}

		all = append(all, resp)
	}

	return all, nil
}

func (c *Client) GetPublications() ([]Publication, error) {
	var publications []Publication
	data, err := c.GetMultipleWrapper("publications")
	if err != nil {
		return nil, err
	}

	for _, d := range data {
		publications = append(publications, d.Publication...)
	}

	return publications, nil
}

func (c *Client) GetProjects() ([]Project, error) {
	var projects []Project
	data, err := c.GetMultipleWrapper("projects")
	if err != nil {
		return nil, err
	}

	for _, d := range data {
		projects = append(projects, d.Project...)
	}

	return projects, nil
}

func (c *Client) getWrapper(endpoint string) (*http.Response, error) {
	url := c.BaseURL + endpoint
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	log.Println("requesting ", url)

	req.Header.Set("User-Agent", "SN-Utils")
	req.Header.Set("Accept", "application/json")
	req.Header.Set(c.keyHeader, c.apiKey)

	return c.HTTP.Do(req)
}

func (p *Parameters) buildRequestString(pubType string, page int) string {
	return fmt.Sprintf("/%s/%s/%s/%d/%d?publicationDateFrom=%s&publicationDateTo=%s&originator=%s", pubType, p.Vendor, p.LastChangeTimestamp, page, p.PageSize, p.FromDate, p.ToDate, p.Originator)
}

func dumpResponse(resp *http.Response) (Response, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return Response{}, err
	}

	return response, nil
}
