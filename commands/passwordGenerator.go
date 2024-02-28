package commands

import (
	"math/rand"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+{}:?><"
)

func PasswordGenerator(leng int) string {

	password := make([]byte, leng)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}

	return string(password)
}

func HandleGenerator(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {

	if len(args) < 1 {
		s.ChannelMessageSend(m.ChannelID, "USO: .newpass [tamaño]")
		return
	}

	toNum, err := strconv.Atoi(args[0])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error del parametro recibido, posiblemente no has ingresado un número!")
		return
	}
	password := PasswordGenerator(toNum)
	s.ChannelMessageSend(m.ChannelID, password)
}
