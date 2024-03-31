package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func ToTest(s *discordgo.Session, m *discordgo.MessageCreate, imageURL string) {

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				&discordgo.Button{
					Label:    "",
					CustomID: "previous",
					Style:    discordgo.PrimaryButton,
					Emoji: discordgo.ComponentEmoji{
						Name: "⏮️",
					},
				},
				&discordgo.Button{
					Label:    "",
					CustomID: "next",
					Style:    discordgo.PrimaryButton,
					Emoji: discordgo.ComponentEmoji{
						Name: "⏭️",
					},
				},
			},
		},
	}

	embed := &discordgo.MessageEmbed{
		Image: &discordgo.MessageEmbedImage{
			URL: imageURL,
		},
	}

	msg, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})
	if err != nil {
		fmt.Printf("Error al enviar mensaje complejo: %v\n", err)
		return
	}
	fmt.Printf("Mensaje enviado: %v\n", msg.ID)
}
