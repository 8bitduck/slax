///
/// Copyright (c) 2015 8bit Duck LLC
/// Author: Bryan Rehbein <bryan@8bitduck.com>
///
/// This code is licensed under the terms of the MIT License, see the LICENSE file for more details
///

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("slax")
	viper.AddConfigPath("/etc/slax")
	viper.AddConfigPath("$HOME/.slax")
	vipererr := viper.ReadInConfig() // Find and read the config file
	if vipererr != nil {             // Handle errors reading the config file
		fmt.Errorf("Fatal error config file: %s \n", vipererr)
	}
}

type SlackText struct {
	Channel   string `json:"channel"`
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	Text      string `json:"text"`
}

func main() {
	var slatCmd = &cobra.Command{
		Use:   "slax",
		Short: "Slax sends personality messages to a slack channel via a webhook",
		Long: `Slax sends messages as different personas to a slack channel using
webhooks. This way you can give voice to your CI bots or any other
character you deem necessary.`,
	}
	var sayCmd = &cobra.Command{
		Use:   "say",
		Short: "Say something",
		Run: func(cmd *cobra.Command, args []string) {
			q := viper.GetBool("quiet")
			if !q {
				fmt.Println("Slat at your command!")
			}

			if len(args) == 0 {
				panic("You need to say something")
			}

			// Build the Slack Text object
			txt := &SlackText{}
			txt.Username = "Jarvis"
			txt.IconEmoji = ":jarvis:"

			// Determine the Persona
			persona := viper.GetString("persona")

			switch persona {
			case "jarvis":
				txt.Username = "Jarvis"
				txt.IconEmoji = ":jarvis:"
			case "butler":
				txt.Username = "Sebastian - One Hell of a Butler"
				txt.IconEmoji = ":butler:"
			case "hodor":
				txt.Username = "Hodor"
				txt.IconEmoji = ":hodor:"
			case "tron":
				txt.Username = "uselesstron9000"
				txt.IconEmoji = ":uselesstron9000:"
			case "redbeard":
				txt.Username = "Capt Redbeard of the Silicon Seas"
				txt.IconEmoji = ":bryan:"
			}

			channel := viper.GetString("channel")
			if strings.HasPrefix(channel, "@") || strings.HasPrefix(channel, "#") {
				txt.Channel = channel
			} else {
				txt.Channel = "#" + channel
			}
			txt.Text = strings.Join(args, " ")

			if !q {
				fmt.Printf("Sending to %s...\n", txt.Channel)
				fmt.Printf("@%s: %s\n", txt.Username, txt.Text)
			}

			url := fmt.Sprintf(
				"https://%s.slack.com/services/hooks/incoming-webhook?token=%s",
				viper.GetString("slack_account_name"),
				viper.GetString("slack_webhook_token"),
			)

			jsonBytes, err := json.Marshal(txt)
			if err != nil {
				panic(err)
			}

			resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))

			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			if bytes.Equal(body, []byte(`ok`)) {
				if !q {
					fmt.Println("Sent!")
				}
			} else {
				fmt.Println("Error sending, sorry.")
			}

		},
	}

	slatCmd.PersistentFlags().StringP("channel", "c", "general", "Slack Channel")
	slatCmd.PersistentFlags().StringP("persona", "p", "jarvis", "Persona to use (jarvis, butler, hodor, tron)")
	slatCmd.PersistentFlags().BoolP("quiet", "q", false, "Quiet operation")
	viper.BindPFlag("channel", slatCmd.PersistentFlags().Lookup("channel"))
	viper.BindPFlag("persona", slatCmd.PersistentFlags().Lookup("persona"))
	viper.BindPFlag("quiet", slatCmd.PersistentFlags().Lookup("quiet"))
	slatCmd.AddCommand(sayCmd)
	slatCmd.Execute()

}
