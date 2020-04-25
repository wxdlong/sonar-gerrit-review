package sonar

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Sonar struct {
	Url           string `yaml:"url"`
	Name          string `yaml:"name"`
	Pwd           string `yaml:"pwd"`
	Timeout       int    `yaml:timeout`
	Author        string `yaml:author`
	Task_url      string `yaml:task_url`
	ComponentKeys string
	CreatedAfter  string
	Tasks         *Tasks
}

//api/issues/search
//Search for issues.
//At most one of the following parameters can be provided at the same time: componentKeys and componentUuids.
//Requires the 'Browse' permission on the specified project(s).
//http://localhost:9000/web_api/api/issues
type NewIssues struct {
	Total  int
	Server string
	Paging struct {
		PageIndex int
		PageSize  int
		Total     int
	}
	Issues []struct {
		Key       string
		Rule      string
		Severity  string //INFO
		Component string
		Project   string
		Line      int
		TextRange struct {
			StartLine   int
			EndLine     int
			StartOffset int
			EndOffset   int
		}
		Status  string //OPEN
		Message string
		Type    string
		Effort  string
	}
}

var sonars = &Sonar{}
var client = &http.Client{}
var checkCount = 0

type Task struct {
	Id              string
	Type            string
	ComponentId     string
	ComponentKey    string
	ComponentName   string
	AnalysisId      string
	Status          string
	SubmittedAt     string
	StartedAt       string
	ExecutionTimeMs int
	ErrorMessage    string
}

type Tasks struct {
	Task Task
}

//waitSonarResult wait sonar scan result from server by the taskURL.
//loop this method when task status is PENDING.
//exit 1 when task status.
func (s *Sonar) waitSonarResult() {
	log.Println("Check Sonar Task Status: ", s.Task_url)
	req, _ := http.NewRequest("GET", s.Task_url, nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth(s.Name, s.Pwd)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	s.Tasks = &Tasks{}
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response statuscode", resp.StatusCode)

	json.Unmarshal(body, s.Tasks)
	log.Println("task status", s.Tasks.Task.Status)

}

func (s *Sonar) WaitSonarResult() bool {
	for i := 0; i < s.Timeout; i++ {
		sonars.waitSonarResult()
		switch sonars.Tasks.Task.Status {
		case "PENDING":
			log.Println("Sonar scanner pending......")
			time.Sleep(30 * time.Second)
			continue
		case "SUCCESS":
			log.Println("Sonar scanner SUCCESS......")
			return true
		default:
			log.Println("Sonar scanner error:", sonars.Tasks.Task.ErrorMessage)
			os.Exit(1)
		}
	}
	return false
}

func (s *Sonar) FetchNewIssues() (newIssues *NewIssues) {
	baseUrl, err := url.Parse(s.Url + "/api/issues/search")

	if err != nil {
		log.Fatal("Malformed URL: ", err.Error())
		os.Exit(1)
	}

	// Add a Path Segment (Path segment is automatically escaped)

	// Prepare Query Parameters
	params := url.Values{}
	params.Add("authors", s.Author)
	params.Add("componentKeys", s.ComponentKeys)
	params.Add("createdAfter", s.CreatedAfter)
	baseUrl.RawQuery = params.Encode()

	log.Println("Fetch new issues: ", baseUrl.String())
	req, _ := http.NewRequest("GET", baseUrl.String(), nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth(s.Name, s.Pwd)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(resp.StatusCode)
	log.Println(string(body))
	newIssues = &NewIssues{}
	json.Unmarshal(body, &newIssues)
	newIssues.Server = s.Url
	return
}

func (s Sonar) String() string {
	return fmt.Sprintf("Sonar Url:%s,  task_url:%s,  Author: %s, CreateAfter:%s", s.Url, s.Task_url, s.Author, s.CreatedAfter)
}
