package gerrit

import (
	"encoding/json"
	"fmt"
	"github.com/wxdlong/sonar-gerrit-review/sonar"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	REVIEW              = "/a/changes/%s/revisions/%s/review" //Review PATH
	ApplicationJsonUtf8 = "application/json; charset=UTF-8"
)

var client = &http.Client{}

//Set Review
//'POST /changes/{change-id}/revisions/{revision-id}/review'
//
//Sets a review on a revision, optionally also publishing draft comments,
//setting labels, adding reviewers or CCs, and modifying the work in progress property.
//
//The review must be provided in the request body as a ReviewInput entity.
//
//A review cannot be set on a change edit. Trying to post a review for a change edit fails with 409 Conflict.
//
//Here is an example of using this method to set labels:
//POST /changes/myProject~master~I8473b95934b5732ac55d26311a706c9c2bde9940/revisions/674ac754f91e64a0efb8087e59a176484bd534d1/review HTTP/1.0
//Content-Type: application/json; charset=UTF-8
//{
//    "tag": "jenkins",
//    "message": "Some nits need to be fixed.",
//    "labels": {
//      "Code-Review": -1
//    },
//    "comments": {
//      "gerrit-server/src/main/java/com/google/gerrit/server/project/RefControl.java": [
//        {
//          "line": 23,
//          "message": "[nit] trailing whitespace"
//        },
//        {
//          "line": 49,
//          "message": "[nit] s/conrtol/control"
//        },
//        {
//          "range": {
//            "start_line": 50,
//            "start_character": 0,
//            "end_line": 55,
//            "end_character": 20
//          },
//          "message": "Incorrect indentation"
//        }
//      ]
//    }
//  }
type ReviewInput struct {
	Message string `json:"message"`
	//Tag      string
	Labels   map[string]string          `json:"labels"`
	Comments map[string][]*CommentInput `json:"comments"`
}

type CommentInput struct {
	// Id      string
	Path  string `json:"path"`
	Line  int    `json:"line,omitempty"`
	Range struct {
		Start_line      int `json:"start_line"`
		Start_character int `json:"start_character"`
		End_line        int `json:"end_line"`
		End_character   int `json:"end_character"`
	} `json:"range"`
	Message string `json:"message"`
}

//Get Change
//'GET /changes/{change-id}'
//
//Retrieves a change.
//
//Additional fields can be obtained by adding o parameters, each option requires more database
//lookups and slows down the query response time to the client so they are generally disabled
//by default. Fields are described in Query Changes.
//
//Request
//GET /changes/myProject~master~I8473b95934b5732ac55d26311a706c9c2bde9940 HTTP/1.0
//As response a ChangeInfo entity is returned that describes the change.
//
//Response
//HTTP/1.1 200 OK
//Content-Disposition: attachment
//Content-Type: application/json; charset=UTF-8
type ChangeInfo struct {
	Id      string `json:"Id"`
	Project string `json:"Project"`
	Branch  string
	Subject string
	Status  string
	Created string
	Updated string
}

type Gerrit struct {
	Url             string `yaml:"url"`
	Name            string `yaml:"name"`
	Pwd             string `yaml:"pwd"`
	ChangeId        string `yaml:"changeId"`
	ReviewId        string `yaml:"reviewId"`
	ReviewInput     *ReviewInput
	ChangeInfo      *ChangeInfo
	CodeReviewLabel string `yaml:"codeReviewLabel"`
}

const msg = `
%s SonarQube violation:
%s

Read more: 
%s
`

func (gerrit *Gerrit) toJson() string {
	reviews, _ := json.MarshalIndent(gerrit.ReviewInput, "", "    ")
	return string(reviews)
}

func (gerrit *Gerrit) PostComment() {
	review_url := fmt.Sprintf(REVIEW, gerrit.ChangeId, gerrit.ReviewId)
	commentJson := gerrit.toJson()
	comments := strings.NewReader(commentJson)
	log.Println("Post review:", review_url)
	log.Println("   Comment:", commentJson)
	req, _ := http.NewRequest("POST", gerrit.Url+review_url, comments)
	req.Header.Set("Content-Type", ApplicationJsonUtf8)
	req.SetBasicAuth(gerrit.Name, gerrit.Pwd)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("responseCode", resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
	if resp.StatusCode != 200 {
		os.Exit(1)
	}
}

func (gerrit *Gerrit) GetChange() {
	change_url := gerrit.Url + "/a/changes/" + gerrit.ChangeId
	log.Println("GetChange: ", change_url)
	req, _ := http.NewRequest("GET", change_url, nil)
	req.Header.Set("Content-Type", ApplicationJsonUtf8)
	req.SetBasicAuth(gerrit.Name, gerrit.Pwd)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("responseCode", resp.StatusCode)

	if resp.StatusCode != 200 {
		os.Exit(1)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
	gerrit.ChangeInfo = &ChangeInfo{}
	json.Unmarshal(body[5:], gerrit.ChangeInfo)
	gerrit.ChangeInfo.Created = timeConvert(gerrit.ChangeInfo.Created)
}

//timeConvert convert time format
func timeConvert(src string) string {
	CTF := "2006-01-02 15:04:05.000000000"
	RFC3339 := "2006-01-02T15:04:05+0000" //The time format used by fetch sonar issues.
	createDate, _ := time.Parse(CTF, src)
	return createDate.Add(-1 * time.Minute).Format(RFC3339)
}

//https://localhost:9001/coding_rules#rule_key=squid%3AClassVariableVisibilityCheck
func (gerrit *Gerrit) Issue2Comment(issues *sonar.NewIssues, sonar_url string) {
	comments := make(map[string][]*CommentInput)

	for _, issue := range issues.Issues {
		comment := &CommentInput{}
		comment.Line = issue.Line
		msgs := fmt.Sprintf(msg, issue.Severity, issue.Message, sonar_url+"/coding_rules#rule_key="+url.PathEscape(issue.Rule))
		comment.Message = msgs
		comment.Path = strings.Replace(issue.Component, issue.Project+":", "", 1)
		comment.Range.Start_line = issue.TextRange.StartLine
		comment.Range.Start_character = issue.TextRange.StartOffset
		comment.Range.End_line = issue.TextRange.EndLine
		comment.Range.End_character = issue.TextRange.EndOffset
		if _, ok := comments[comment.Path]; ok {
			comments[comment.Path] = append(comments[comment.Path], comment)
		} else {
			comments[comment.Path] = []*CommentInput{comment}
		}
	}
	gerrit.ReviewInput = &ReviewInput{}
	if len(comments) == 0 {
		gerrit.ReviewInput.Message = "Nice Code!"
		gerrit.ReviewInput.Labels = map[string]string{gerrit.CodeReviewLabel: "+1"}
	} else {
		gerrit.ReviewInput.Message = "Sonar Issue Found,Please fix and recommit."
		gerrit.ReviewInput.Labels = map[string]string{gerrit.CodeReviewLabel: "-2"}
		gerrit.ReviewInput.Comments = comments
	}
}

func (gerrit Gerrit) String() string {
	return fmt.Sprintf("Gerrit Url:%s, changeId:%s, reviewIed:%s", gerrit.Url, gerrit.ChangeId, gerrit.ReviewId)
}
