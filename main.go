package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Styerr-Network/internal-bot/commands"
	"github.com/Styerr-Network/internal-bot/commands/music"
	"github.com/Styerr-Network/internal-bot/models"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session   *discordgo.Session
	Config    models.ConfigFile
	Navigator *commands.ImageNavigator
}

var imageNav *commands.ImageNavigator

func main() {

	cfg, err := os.ReadFile("config.json")
	if err != nil {
		panic("Fail to read the token")
	}

	var config models.ConfigFile

	err = json.Unmarshal(cfg, &config)
	if err != nil {
		panic("Fail to read the token")
	}

	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	bot := Bot{
		Session:   dg,
		Config:    config,
		Navigator: commands.NewImageNavigator(config),
	}
	imageNav = commands.NewImageNavigator(config)

	// Handlers
	dg.AddHandler(bot.pushCommands)
	dg.AddHandler(bot.handleInteraction)

	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer dg.Close()

	fmt.Println("El bot está corriendo. Presiona Control + C para salir.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func (bot *Bot) pushCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, ".") {
		return
	}

	parts := strings.Fields(m.Content)

	command := parts[0]
	args := parts[1:]
	img := "https://images-ext-2.discordapp.net/external/tMSvRMWE8HTxKFYc9wTaJnl0-H2NiAdW0IvO7PVHHxM/https/t4.ftcdn.net/jpg/04/23/14/25/360_F_423142507_FsZUpYT6eamfNgyPxRlezjyx8eV1tlXz.jpg?format=webp"

	switch command {
	case ".testapi":
		commands.TestApi(s, m)
	case ".newpass":
		commands.HandleGenerator(s, m, args)
	case ".img":
		bot.Navigator.GetImages(s, m.ChannelID, convertToString(args), 0)
	// Music things
	case ".play":
		music.Play(s, m, "")
	case ".test":
		commands.ToTest(s, m, img)
	}
}

func convertToString(arg []string) string {
	return fmt.Sprintf("%s", arg)
}

func (bot *Bot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	switch i.MessageComponentData().CustomID {
	case "previous":
		if bot.Navigator.CurrentIndex > 0 {
			bot.Navigator.CurrentIndex--
		}
	case "next":
		if bot.Navigator.CurrentIndex < len(bot.Navigator.SearchResults)-1 {
			bot.Navigator.CurrentIndex++
			bot.UpdateMessage(s, i)
		}
	default:
		return
	}
}

func (bot *Bot) UpdateMessage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	imageURL := bot.Navigator.SearchResults[bot.Navigator.CurrentIndex]
	components := bot.makeMessageComponents() // Construye los componentes (botones) para la navegación
	embed := bot.makeImageEmbed(imageURL)     // Construye el embed con la imagen actual

	// Actualiza el mensaje de interacción con la nueva imagen y componentes
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds:     &[]*discordgo.MessageEmbed{embed},
		Components: &components,
	})
	if err != nil {
		fmt.Printf("Error al actualizar el mensaje de interacción: %v\n", err)
	}
}

func (bot *Bot) makeImageEmbed(imageURL string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Image: &discordgo.MessageEmbedImage{
			URL: imageURL,
		},
	}
}

func (bot *Bot) makeMessageComponents() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
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
}
