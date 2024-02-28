package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Styerr-Network/internal-bot/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// Definir un tipo para las funciones de comando
type CommandFunc func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Manejador de mensajes
	dg.AddHandler(pushCommands)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer dg.Close()

	fmt.Println("El bot está corriendo. Presiona Control + C para salir.")

	// Esperar señales de terminación
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func pushCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, ".") {
		return
	}

	parts := strings.Fields(m.Content)

	command := parts[0]
	args := parts[1:]

	switch command {
	case ".testapi":
		commands.TestApi(s, m)
	case ".newpass":
		commands.HandleGenerator(s, m, args)
	}
}
