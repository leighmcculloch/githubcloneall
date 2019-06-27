package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/fatih/color"
)

var username string = "";
var directory string = "";
var token string = "";
var githubType string = "";

func main() {
	flag.StringVar(&username, "u", "", "GitHub username")
	printHelp := flag.Bool("h", false, "Print help")
	flag.StringVar(&directory, "d", "", "Output directory");
	flag.StringVar(&token, "token", "", "Github personal access token or oauth token");
	flag.StringVar(&githubType, "type", "users", "github type (users, orgs)")

	flag.Usage = func() {
		fmt.Println("Usage: githubcloneall -u username -d dir -token TOKEN -type orgs")
		fmt.Println("")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *printHelp || username == "" {
		flag.Usage()
		return
	}

	checkForMoreRepos(1);	
}

func checkForMoreRepos(pageNumber int) {
	client := &http.Client{};
	req, err := http.NewRequest("GET", "https://api.github.com/" + githubType + "/" + username + "/repos?per_page=200&page=" + strconv.Itoa(pageNumber), nil)
	req.Header.Add("Authorization", "token " + token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	dec := json.NewDecoder(bytes.NewReader(body))
	repos := []Repo{}
	err = dec.Decode(&repos)
	if err != nil {
		fmt.Printf("Error: %s\n%s\n", err, string(body))
		return
	}

	if len(repos) != 0 {
		downloadRepos(repos);
		checkForMoreRepos(pageNumber + 1);
	}
}

func downloadRepos(repos []Repo) {
	for i, r := range repos {
		if exists(r.Name) {
			color.Yellow("%d/%d Skipping already cloned repo %s.\n", i, len(repos), r.SSHURL)
			continue
		}
		if r.Archived {
			color.Yellow("%d/%d Skipping archived repo %s.\n", i, len(repos), r.SSHURL)
			continue
		}
		color.Green("%d/%d Cloning repo %s:\n", i, len(repos), r.SSHURL)
		cmd := exec.Command("git", "clone", "--depth", "1", r.SSHURL, directory + "/" + r.Name);
		color.Yellow("Running command: %s", cmd);
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	}
}

type Repo struct {
	Name     string `json:"name"`
	SSHURL   string `json:"ssh_url"`
	Archived bool   `json:"archived"`
}

func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
