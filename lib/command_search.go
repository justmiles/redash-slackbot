package redashslackbot

import (
	"fmt"
	"regexp"

	"github.com/shomali11/slacker"
	"github.com/slack-go/slack"
)

// Search for queries
func Search(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {

	var match = regexp.MustCompile(`search (.*)`).FindStringSubmatch(botCtx.Event().Msg.Text)
	var searchQuery string
	if len(match) == 0 {
		searchQuery = "*"
	} else {
		searchQuery = match[1]
	}

	queries, err := client.SearchQueries(searchQuery, false, 5)
	if err != nil {
		response.Reply(err.Error())
	}

	if len(queries) == 0 {
		response.Reply("No matching queries")
		return
	}

	var attachments []slack.Attachment
	for _, query := range queries {

		query, err := client.GetQueryByID(query.ID)
		if err != nil {
			response.Reply(err.Error())
		}

		attachment := slack.Attachment{
			Title:      query.Name,
			Text:       query.Description,
			TitleLink:  fmt.Sprintf("%s/queries/%d", client.BaseURL, query.ID),
			AuthorName: query.User.Name,
		}

		for _, vis := range query.Visualizations {
			attachment.Text = fmt.Sprintf("%s\n - <%s/embed/query/%d/visualization/%d?api_key=%s|%s>", attachment.Text, client.BaseURL, query.ID, vis.ID, query.APIKey, vis.Name)
			// url := "https://redash.prd.i-edo.net/embed/query/752/visualization/904?api_key=YsONjjfGS1E8G9NwGAzPMyF1TCUUZkA0bfGRgpum"

		}

		attachments = append(attachments, attachment)

	}

	response.Reply("", slacker.WithAttachments(attachments))

}
