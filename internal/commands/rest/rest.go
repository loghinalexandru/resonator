package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type REST struct {
	URL       string
	Type      any
	Formatter func(payload any) string
	Def       *discordgo.ApplicationCommand
}

func (cmd *REST) Definition() *discordgo.ApplicationCommand {
	return cmd.Def
}

func (cmd *REST) Handler(sess *discordgo.Session, inter *discordgo.InteractionCreate) error {
	var args []any
	for _, v := range inter.ApplicationCommandData().Options {
		args = append(args, v.Value.(string))
	}

	customURL := fmt.Sprintf(cmd.URL, args...)
	response, err := http.Get(customURL)

	if err != nil {
		return nil
	}

	defer response.Body.Close()
	content, _ := io.ReadAll(response.Body)

	err = json.Unmarshal(content, cmd.Type)

	if err != nil {
		return err
	}

	sess.InteractionRespond(inter.Interaction, cmd.createReponse())
	return nil
}

func (cmd *REST) createReponse() *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: cmd.Formatter(cmd.Type),
		},
	}
}