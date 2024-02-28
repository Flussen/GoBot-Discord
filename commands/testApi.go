package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type userData struct {
	UserID    int
	ID        int
	Title     string
	completed bool
}

func (user *userData) toString() string {
	return fmt.Sprintf("userID: %d, ID: %d, Title: %s, Completed: %t", user.UserID, user.ID, user.Title, user.completed)
}

func TestApi(s *discordgo.Session, m *discordgo.MessageCreate) {
	url := "https://jsonplaceholder.typicode.com/todos/1"

	response, err := http.Get(url)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "fallo al regresar los datos del api")
		return
	}
	defer response.Body.Close()

	var user userData

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error al leer el cuerpo de la respuesta")
		return
	}

	err = json.Unmarshal([]byte(responseBody), &user)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error al leer el cuerpo de la respuesta")
		return
	}

	s.ChannelMessageSend(m.ChannelID, user.toString())
}
