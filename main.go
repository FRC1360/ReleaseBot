package main

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"strings"
	"fmt"
	"os"
	"net/http"
	"io"
	"github.com/bluele/slack"

)

func main() {

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_TOKEN},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	githubClient := github.NewClient(tc)
	release, err := githubClient.Repositories.GetLatestRelease("FRC1360","Stronghold2016")



	if(err != nil) {
		panic(err)
	}

	if downloadFromUrl(release.ZipballURL) {

		message := "Code Release Downloaded: " + release.AssetsURL + "\nSaved to Backup Server - Running copy cron job now."
		api := slack.New(SLACK_TOKEN)
		channel, err := api.FindChannelByName(SLACK_CHANNEL)
		if err != nil {
			panic(err)
		}
		err = api.ChatPostMessage(channel.Id, message, nil)
		if err != nil {
			panic(err)
		}
	}

}


func downloadFromUrl(url string)(downloaded bool){
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)
	if os.IsExist(fileName) {
		fmt.Println("File already downloaded, ignoring.")
		return false
	}

	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return false
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return false
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return false
	}

	fmt.Println(n, "bytes downloaded.")
	return true
}