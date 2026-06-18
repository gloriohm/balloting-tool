package sdimport

import (
	"ballot-tool/internal/utils/logging"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
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
	Logger    *slog.Logger
}

func NewClient(dev bool, params Parameters) *Client {
	url, key := setEnv(dev)
	file, _ := logging.NewFile("errors.log")
	return &Client{
		HTTP:      &http.Client{Timeout: 30 * time.Second},
		BaseURL:   url,
		apiKey:    key,
		keyHeader: keyHeader,
		Params:    params,
		Logger:    slog.New(slog.NewTextHandler(file, nil)),
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
		return Publication{}, fmt.Errorf("no publications in response")
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

	if totalRecords == 0 {
		return nil, nil
	}

	lastPage := (totalRecords - 1) / c.Params.PageSize

	for page := 0; page <= lastPage; page++ {
		var resp Response

		if page == 0 {
			resp, err = dumpResponse(first)
			if err != nil {
				return nil, fmt.Errorf("error dumping response for page %d: %w", page, err)
			}
		} else {
			path := c.Params.buildRequestString(pubType, page)
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

func (c *Client) GetPublicationByProject(urn string, pubStatus string) (Publication, error) {
	status := make(map[string]struct{})
	switch pubStatus {
	case "published", "withdrawn":
		status["PUBLISHED"] = struct{}{}
		status["WITHDRAWN"] = struct{}{}
	case "draft":
		status["DRAFT"] = struct{}{}
	}
	proj, err := c.GetProject(urn)
	if err != nil {
		return Publication{}, fmt.Errorf("could not get project by %s", urn)
	}

	pubs := proj.getPublicationURNs()
	if len(pubs) == 0 {
		return Publication{}, fmt.Errorf("project %s has no publications", urn)
	}

	for _, p := range pubs {
		pub, err := c.GetPublication(p)
		if err != nil {
			log.Printf("failed to get publication with urn %s", p)
			c.Logger.Info(fmt.Sprintf("%s: %s", p, err))
			continue
		}

		if _, ok := status[pub.Status]; ok {
			return pub, nil
		}
	}

	return Publication{}, fmt.Errorf("could not get publication with status %s for project %s", status, urn)
}

func (c *Client) GetFile(ref ContentRef, dir string) error {
	url := c.BaseURL + ref.URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	log.Println("requesting ", url)

	req.Header.Set("User-Agent", "SN-Utils")
	req.Header.Set(c.keyHeader, c.apiKey)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	path := filepath.Join(dir, ref.FileName)

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}

	return nil
}
