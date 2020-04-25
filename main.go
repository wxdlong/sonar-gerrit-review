package main

import (
	"flag"
	"github.com/wxdlong/sonar-gerrit-review/gerrit"
	"github.com/wxdlong/sonar-gerrit-review/sonar"
	"log"
	"os"
)

var (
	gerrit_url          string
	gerrit_change_id    string
	gerrit_review_id    string
	gerrit_name         string
	gerrit_pwd          string
	gerrit_review_label string

	sonar_url      string
	sonar_name     string
	sonar_pwd      string
	sonar_task_url string
	sonar_author   string
	sonar_timeout  int
)

func init() {
	flag.StringVar(&gerrit_url, "gurl", "https://localhost:8080", "gerrit url")
	flag.StringVar(&gerrit_change_id, "gcid", "", "gerrit change id")
	flag.StringVar(&gerrit_review_id, "grid", "", "gerrit review id")
	flag.StringVar(&gerrit_name, "gname", "admin", "gerrit name")
	flag.StringVar(&gerrit_pwd, "gpwd", "admin", "gerrit pwd")
	flag.StringVar(&gerrit_review_label, "glabel", "Code-Review", "gerrit review label")
	flag.StringVar(&sonar_url, "surl", "https://localhost:9001", "sonar url")
	flag.StringVar(&sonar_name, "sname", "admin", "sonar name")
	flag.StringVar(&sonar_pwd, "spwd", "admin", "sonar pwd")
	flag.StringVar(&sonar_task_url, "stask", "", "sonar task url")
	flag.StringVar(&sonar_author, "sauthor", "", "sonar author")
	flag.IntVar(&sonar_timeout, "stimeout", 40, "sonar timeout for checking task status")
}

func main() {
	flag.Parse()

	sonars := sonar.Sonar{Url: sonar_url,
		Name:     sonar_name,
		Pwd:      sonar_pwd,
		Author:   sonar_author,
		Task_url: sonar_task_url,
		Timeout:  sonar_timeout,
	}

	gerrits := gerrit.Gerrit{
		Url:      gerrit_url,
		Name:     gerrit_name,
		Pwd:      gerrit_pwd,
		ChangeId: gerrit_change_id,
		ReviewId: gerrit_review_id,
		//CodeReviewLabel: "Sonar-Verified",
		CodeReviewLabel: gerrit_review_label,
	}

	log.Println(sonars)
	log.Println(gerrits)
	if !sonars.WaitSonarResult() {
		log.Println("WaitSonarResult Failed!")
		os.Exit(1)
	}

	gerrits.GetChange()
	sonars.CreatedAfter = gerrits.ChangeInfo.Created
	newIssues := sonars.FetchNewIssues()

	gerrits.Issue2Comment(newIssues, sonars.Url)
	gerrits.PostComment()
}
