package redashslackbot

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

// Show a visualization
func Show(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {

	var match = regexp.MustCompile(`show (.*)`).FindStringSubmatch(botCtx.Event().Msg.Text)
	var searchQuery string
	if len(match) == 0 {
		response.Reply("Show what? Try `show <query>`")
		return
	}
	searchQuery = match[1]

	queries, err := client.SearchQueries(searchQuery, false, 1)
	if err != nil {
		response.Reply(err.Error())
	}

	if len(queries) == 0 {
		response.Reply("No matching queries")
		return
	}

	query, err := client.GetQueryByID(queries[0].ID)
	if err != nil {
		response.Reply(err.Error())
	}
	var attachments []slack.Attachment

	for _, vis := range query.Visualizations {
		if vis.Type == "TABLE" {
			continue
		}
		embeddedURL := fmt.Sprintf("%s/embed/query/%d/visualization/%d?api_key=%s", client.BaseURL, query.ID, vis.ID, query.APIKey)
		fmt.Println(embeddedURL)

		file := GetVisualization(embeddedURL, fmt.Sprintf("redash_%d_%d.png", query.ID, vis.ID))
		defer os.Remove(file.Name())
		client := botCtx.Client()

		slackFile, err := client.UploadFile(slack.FileUploadParameters{
			File:     file.Name(),
			Filename: fmt.Sprintf("%s (%s)", query.Name, vis.Name),
			Channels: []string{"redash-images"},
		})
		if err != nil {
			response.Reply(err.Error())
		}

		attachment := slack.Attachment{
			Title:     fmt.Sprintf("%s (%s)", query.Name, vis.Name),
			TitleLink: embeddedURL,
			ImageURL:  slackFile.Permalink,
		}

		attachments = append(attachments, attachment)

	}

	if len(attachments) == 0 {
		response.Reply("This query has no visualisations")
	}

	// TODO: consider using blocks!
	// blocks := []slack.Block{}
	// blocks = append(blocks, slack.NewContextBlock("1",
	// 	slack.NewTextBlockObject("mrkdwn", "Hi!", false, false)),
	// )
	// response.Reply("", slacker.WithAttachments(attachments), slacker.WithBlocks(blocks))
	response.Reply("", slacker.WithAttachments(attachments))

}

// GetVisualization uses chromedp to screenshot the visulization page
func GetVisualization(url, filename string) *os.File {

	file, err := os.Create(path.Join(os.TempDir(), filename))
	if err != nil {
		log.Fatal(err)
	}

	// Start Chrome
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(url, &buf)); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(file.Name(), buf, 0644); err != nil {
		log.Fatal(err)
	}
	return file
}

func elementScreenshot(url string, res *[]byte) chromedp.Tasks {
	width, height := 1024, 1080
	return chromedp.Tasks{
		emulation.SetDeviceMetricsOverride(int64(width), int64(height), 1.0, false),
		chromedp.Navigate(url),
		chromedp.WaitVisible(`document.querySelector("#app-content > visualization-embed > div")`, chromedp.ByJSPath),
		chromedp.SetAttributeValue(`document.querySelector("#app-content > visualization-embed > div > div.tile__bottom-control > span.hidden-print")`, `class`, `hidden`, chromedp.ByJSPath),
		chromedp.Screenshot(`document.querySelector("#app-content > visualization-embed > div")`, res, chromedp.NodeVisible, chromedp.ByJSPath),
	}
}
