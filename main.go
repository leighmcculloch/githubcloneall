package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

func main() {
	username := flag.String("u", "", "GitHub username")
	printHelp := flag.Bool("h", false, "Print help")

	flag.Usage = func() {
		fmt.Println("Usage: githubcloneall -u username")
		fmt.Println("")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *printHelp || *username == "" {
		flag.Usage()
		return
	}

	resp, err := http.Get("https://api.github.com/user/" + *username + "/repos?per_page=200")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	dec := json.NewDecoder(resp.Body)
	repos := []Repo{}
	err = dec.Decode(&repos)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

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
		cmd := exec.Command("git", "clone", r.SSHURL)
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
