package gerrit

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"testing"
	"time"
)

func TestHellos(t *testing.T) {
	labels := make(map[string]string)
	labels["code-review"] = "-2"
	review := &ReviewInput{}
	review.Message = "Hello Message"
	review.Labels = labels
	reviews, _ := json.MarshalIndent(review, "", "    ")

	log.Println(string(reviews))
}

func TestAuthPwd(t *testing.T) {
	auth := "admin" + ":" + "admin"
	log.Println(base64.StdEncoding.EncodeToString([]byte(auth)))
}

func TestSS(t *testing.T) {
	ct := "2020-04-03 04:58:35.727000000"
	CTF := "2006-01-02 15:04:05.000000000"
	RFC3339 := "2006-01-02T15:04:05+0000" //The time format used by fetch sonar issues.
	createDate, _ := time.Parse(CTF, ct)
	fmt.Println(createDate.Add(-5 * time.Minute).Format(RFC3339))
}

func TestReviewInput(t *testing.T) {
	ci1 := &CommentInput{
		Path: "test.java",
		Line: 0,

		Message: "Var name is wrong!",
	}

	ci1.Range.Start_line = 9
	ci1s, err := json.MarshalIndent(ci1, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(string(ci1s))

	ci2 := &CommentInput{
		Path: "test2.java",
		Line: 0,

		Message: "Value name is wrong!",
	}

	reviewInput := &ReviewInput{
		Message:  "Sonar Issue found. fix it!",
		Labels:   map[string]string{"Code-Review": "+2"},
		Comments: map[string][]*CommentInput{"test.java": []*CommentInput{ci1, ci2}},
	}

	ris, err := json.MarshalIndent(reviewInput, "", "    ")
	log.Println(ris)

}

func TestReviewInputJson(t *testing.T) {
	data := `
{
   "tag": "jenkins",
   "message": "Some nits need to be fixed.",
   "labels": {
     "Code-Review": -1
   },
   "comments": {
     "gerrit-server/src/main/java/com/google/gerrit/server/project/RefControl.java": [
       {
         "line": 23,
         "message": "[nit] trailing whitespace"
       },
       {
         "line": 49,
         "message": "[nit] s/conrtol/control"
       },
       {
         "range": {
           "start_line": 50,
           "start_character": 0,
           "end_line": 55,
           "end_character": 20
         },
         "message": "Incorrect indentation"
       }
     ]
   }
 }
`

	ri := &ReviewInput{}
	json.Unmarshal([]byte(data), ri)
	log.Println(ri)
}

func TestJso(t *testing.T) {
	abc := map[string]int{"gerrit.CodeReviewLabel": +1}
	res, _ := json.MarshalIndent(abc, "", "    ")
	fmt.Println(string(res))
}

func TestYaml(t *testing.T) {
	data := `
Url: "http://localhost:8080"
Name: wxdlong
`
	var gerrit map[string]string
	err := yaml.Unmarshal([]byte(data), &gerrit)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Println(gerrit)
}

func TestGetChange(t *testing.T) {
	changes := `)]}'
  {
    "id": "myProject~master~I8473b95934b5732ac55d26311a706c9c2bde9940",
    "project": "myProject",
    "branch": "master",
    "change_id": "I8473b95934b5732ac55d26311a706c9c2bde9940",
    "subject": "Implementing Feature X",
    "status": "NEW",
    "created": "2013-02-01 09:59:32.126000000",
    "updated": "2013-02-21 11:16:36.775000000",
    "mergeable": true,
    "insertions": 34,
    "deletions": 101,
    "_number": 3965,
    "owner": {
      "name": "John Doe"
    }
  }
`
	change := &ChangeInfo{}
	err := json.Unmarshal([]byte(changes[4:]), change)
	if err != nil {
		fmt.Println(err)
	}
	log.Println(change)
}
