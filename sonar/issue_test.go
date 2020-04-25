package sonar

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	fmt.Println(time.Now().Format("2006-01-02T15:04:05Z0700"))
}

func TestUnDYJSON(t *testing.T) {
	jsonStr := `{"abc":true,"intv":0,"listv":[{"name":"wxdlong"},{"age":44}]}`
	dynamic := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &dynamic)
	fmt.Println(dynamic)
}

func TestHello45(t *testing.T) {
	sonars.FetchNewIssues()
}

//api/ce/task
func TestCETask(t *testing.T) {
	ceTaskStr := `
  {
    "task": {
      "organization": "my-org-1",
      "id": "AVAn5RKqYwETbXvgas-I",
      "type": "REPORT",
      "componentId": "AVAn5RJmYwETbXvgas-H",
      "componentKey": "project_1",
      "componentName": "Project One",
      "componentQualifier": "TRK",
      "analysisId": "123456",
      "status": "FAILED",
      "submittedAt": "2015-10-02T11:32:15+0200",
      "startedAt": "2015-10-02T11:32:16+0200",
      "executedAt": "2015-10-02T11:32:22+0200",
      "executionTimeMs": 5286,
      "errorMessage": "Fail to extract report AVaXuGAi_te3Ldc_YItm from database",
      "logs": false,
      "hasErrorStacktrace": true,
      "errorStacktrace": "java.lang.IllegalStateException: Fail to extract report AVaXuGAi_te3Ldc_YItm from database\n\tat org.sonar.server.computation.task.projectanalysis.step.ExtractReportStep.execute(ExtractReportStep.java:50)",
      "scannerContext": "SonarQube plugins:\n\t- Git 1.0 (scmgit)\n\t- Java 3.13.1 (java)",
      "hasScannerContext": true
    }
  }
  `

	ceTask := &Tasks{}
	err := json.Unmarshal([]byte(ceTaskStr), ceTask)
	if err != nil {
		t.Fatal(err)
	}

	if ceTask.Task.Status != "FAILED" {
		t.Fatal("ceTask status should be FAILED but is", ceTask.Task.Status)
	}

	if ceTask.Task.AnalysisId != "123456" {
		t.Fatal("ceTask AnalysisId should be 123456, but ", ceTask.Task.AnalysisId)
	}

	if !strings.Contains(ceTask.Task.ErrorMessage, "Fail to extract report") {
		t.Fatal("ceTask ErrorMessage should contain `Fail to extract report`, but ", ceTask.Task.ErrorMessage)
	}
}

func TestSearchIssue(t *testing.T) {
	issues := `
   {
     "paging": {
       "pageIndex": 1,
       "pageSize": 100,
       "total": 1
     },
     "issues": [
       {
         "key": "01fc972e-2a3c-433e-bcae-0bd7f88f5123",
         "component": "com.github.kevinsawicki:http-request:com.github.kevinsawicki.http.HttpRequest",
         "project": "com.github.kevinsawicki:http-request",
         "rule": "checkstyle:com.puppycrawl.tools.checkstyle.checks.coding.MagicNumberCheck",
         "status": "RESOLVED",
         "resolution": "FALSE-POSITIVE",
         "severity": "MINOR",
         "message": "'3' is a magic number.",
         "line": 81,
         "hash": "a227e508d6646b55a086ee11d63b21e9",
         "author": "Developer 1",
         "effort": "2h1min",
         "creationDate": "2013-05-13T17:55:39+0200",
         "updateDate": "2013-05-13T17:55:39+0200",
         "tags": [
           "bug"
         ],
         "type": "RELIABILITY",
         "comments": [
           {
             "key": "7d7c56f5-7b5a-41b9-87f8-36fa70caa5ba",
             "login": "john.smith",
             "htmlText": "Must be &quot;final&quot;!",
             "markdown": "Must be \"final\"!",
             "updatable": false,
             "createdAt": "2013-05-13T18:08:34+0200"
           }
         ],
         "attr": {
           "jira-issue-key": "SONAR-1234"
         },
         "transitions": [
           "unconfirm",
           "resolve",
           "falsepositive"
         ],
         "actions": [
           "comment"
         ],
         "textRange": {
           "startLine": 2,
           "endLine": 2,
           "startOffset": 0,
           "endOffset": 204
         },
         "flows": [
           {
             "locations": [
               {
                 "textRange": {
                   "startLine": 16,
                   "endLine": 16,
                   "startOffset": 0,
                   "endOffset": 30
                 },
                 "msg": "Expected position: 5"
               }
             ]
           },
           {
             "locations": [
               {
                 "textRange": {
                   "startLine": 15,
                   "endLine": 15,
                   "startOffset": 0,
                   "endOffset": 37
                 },
                 "msg": "Expected position: 6"
               }
             ]
           }
         ]
       }
     ],
     "components": [
       {
         "key": "com.github.kevinsawicki:http-request:src/main/java/com/github/kevinsawicki/http/HttpRequest.java",
         "enabled": true,
         "qualifier": "FIL",
         "name": "HttpRequest.java",
         "longName": "src/main/java/com/github/kevinsawicki/http/HttpRequest.java",
         "path": "src/main/java/com/github/kevinsawicki/http/HttpRequest.java"
       },
       {
         "key": "com.github.kevinsawicki:http-request",
         "enabled": true,
         "qualifier": "TRK",
         "name": "http-request",
         "longName": "http-request"
       }
     ],
     "rules": [
       {
         "key": "checkstyle:com.puppycrawl.tools.checkstyle.checks.coding.MagicNumberCheck",
         "name": "Magic Number",
         "status": "READY",
         "lang": "java",
         "langName": "Java"
       }
     ],
     "users": [
       {
         "login": "admin",
         "name": "Administrator",
         "active": true,
         "avatar": "ab0ec6adc38ad44a15105f207394946f"
       }
     ]
   }
   `

	newIssues := &NewIssues{}
	err := json.Unmarshal([]byte(issues), newIssues)
	if err != nil {
		log.Fatal(err)
	}
	if newIssues.Paging.PageIndex != 1 {
		log.Fatal("pageIndex should be 1, but ", newIssues.Paging.PageIndex)
	}

	if len(newIssues.Issues) != 1 {
		t.Fatal("Issues Num should be 1, but ", len(newIssues.Issues))
	}

	if newIssues.Issues[0].Line != 81 {
		t.Fatal("Issue 1 line should be 81, but ", newIssues.Issues[0].Line)
	}

	if newIssues.Issues[0].Status != "RESOLVED" {
		t.Fatal("Issue 1 status should be RESOLVED, but ", newIssues.Issues[0].Status)
	}

	if newIssues.Issues[0].Rule != "checkstyle:com.puppycrawl.tools.checkstyle.checks.coding.MagicNumberCheck" {
		t.Fatal("Issue 1 Rule should be checkstyle:com.puppycrawl.tools.checkstyle.checks.coding.MagicNumberCheck, but ", newIssues.Issues[0].Rule)
	}

	if newIssues.Issues[0].Effort != "2h1min" {
		t.Fatal("Issue 1 Effort should be 2h1min, but ", newIssues.Issues[0].Effort)
	}

	if newIssues.Issues[0].TextRange.StartLine != 2 {
		t.Fatal("Issue 1 Text Range.StartLine should be 2, but ", newIssues.Issues[0].TextRange.StartLine)
	}
}
