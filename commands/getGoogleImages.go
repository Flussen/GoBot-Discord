package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Styerr-Network/internal-bot/models"
	"github.com/bwmarrin/discordgo"
)

type GoogleSearchResponse struct {
	Items []struct {
		Link string `json:"link"`
	} `json:"items"`
}

type ImageNavigator struct {
	SearchQuery   string
	CurrentIndex  int
	SearchResults []string
	Config        models.ConfigFile
}

func NewImageNavigator(config models.ConfigFile) *ImageNavigator {
	return &ImageNavigator{
		Config:        config,
		SearchResults: []string{},
		CurrentIndex:  0,
	}
}

func (nav *ImageNavigator) GetImages(s *discordgo.Session, channelID string, arg string, startIndex int) {
	// Verificar si la búsqueda ya está en curso y si es una solicitud para la siguiente imagen
	if arg == "" && len(nav.SearchResults) > 0 {
		// Incrementa el índice para obtener la siguiente imagen
		nav.CurrentIndex++
		if nav.CurrentIndex < len(nav.SearchResults) {
			nav.SendImage(s, channelID, nav.SearchResults[nav.CurrentIndex])
			return
		}
		// Si no hay más imágenes en los resultados actuales, podrías cargar más aquí
	}

	if len(arg) < 1 {
		s.ChannelMessageSend(channelID, "Por favor utiliza un argumento, .img [busqueda]")
		return
	}

	nav.SearchQuery = arg // Guarda la búsqueda actual en caso de que el usuario quiera "navegar"
	nav.CurrentIndex = 0  // Restablece el índice para una nueva búsqueda

	encodedArg := url.QueryEscape(arg)
	URL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&searchType=image&start=%d",
		nav.Config.GoogleAPIKey, nav.Config.GoogleSearchEngineID, encodedArg, nav.CurrentIndex+1)

	resp, err := http.Get(URL)
	if err != nil {
		s.ChannelMessageSend(channelID, "Error al buscar la imagen")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.ChannelMessageSend(channelID, "Error al leer la respuesta")
		return
	}

	var searchResponse GoogleSearchResponse
	if err := json.Unmarshal(body, &searchResponse); err != nil {
		s.ChannelMessageSend(channelID, "Error al procesar la respuesta")
		return
	}

	// Almacenar todos los enlaces de imágenes válidos
	nav.SearchResults = []string{}
	for _, item := range searchResponse.Items {
		if isImageURL(item.Link) {
			nav.SearchResults = append(nav.SearchResults, item.Link)
		}
	}

	if len(nav.SearchResults) > 0 {
		nav.SendImage(s, channelID, nav.SearchResults[nav.CurrentIndex])
	} else {
		s.ChannelMessageSend(channelID, "No se encontraron imágenes válidas.")
	}
}

func isImageURL(url string) bool {
	lowerURL := strings.ToLower(url)
	return strings.HasSuffix(lowerURL, ".jpg") || strings.HasSuffix(lowerURL, ".jpeg") ||
		strings.HasSuffix(lowerURL, ".png") || strings.HasSuffix(lowerURL, ".gif") ||
		strings.HasSuffix(lowerURL, ".bmp") || strings.HasSuffix(lowerURL, ".webp")
}

func (nav *ImageNavigator) SendImage(s *discordgo.Session, channelID string, imageURL string) {
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

	_, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})
	if err != nil {
		fmt.Printf("Error al enviar mensaje complejo: %v\n", err)
		return
	}
}
