package command

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/loghinalexandru/resonator/pkg/rest"
)

type animeData struct {
	Id    string `json:"id"`
	Stats struct {
		Title    string `json:"canonicalTitle"`
		Type     string `json:"showType"`
		Rated    string `json:"ageRating"`
		Status   string `json:"status"`
		Start    string `json:"startDate"`
		End      string `json:"endDate"`
		Episodes int32  `json:"episodeCount"`
		Length   int32  `json:"episodeLength"`
		Total    int32  `json:"totalLength"`
	} `json:"attributes"`
}

type animeWrapper struct {
	Content []animeData `json:"data"`
}

func NewAnime() *rest.REST[animeWrapper] {
	url := "https://kitsu.io/api/edge/anime?filter[text]=%v&page[limit]=10"
	def := &discordgo.ApplicationCommand{
		Name:        "anime",
		Description: "This command is used find anime via Kitsu API!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "keyword",
				Description: "Keyword to search for: ",
				Required:    true,
			},
		},
	}

	return rest.New(def, url, http.DefaultClient, animeFormatter)
}

func animeFormatter(resp animeWrapper) string {
	var sb strings.Builder

	if len(resp.Content) > 0 {
		sb.WriteString("Here's a list: \n")
		for _, r := range resp.Content {
			sb.WriteString(fmt.Sprintf("**%s - %s**\n", r.Stats.Title, strings.ToUpper(r.Stats.Status)))
			sb.WriteString(fmt.Sprintf("> Type: %v, Episodes: %v, Length:  %vm, Total Time: %vm\n", r.Stats.Type, r.Stats.Episodes, r.Stats.Length, r.Stats.Total))
		}
	} else {
		sb.WriteString("No match found :(")
	}

	return sb.String()
}
