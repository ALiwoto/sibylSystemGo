package tests

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	sibylSystemGo "github.com/ALiwoto/sibylSystemGo/sibylSystem"
)

func TestBanUser01(t *testing.T) {
	token := getToken()
	if len(token) == 0 {
		log.Println("token is empty; exiting")
	}

	client := sibylSystemGo.NewClient(token, sibylSystemGo.GetDefaultConfig())
	const reason01 = "Spam adding +99 members to an anime group"
	const msg = "https://t.me/telegram/832471"
	const src = "https://t.me/AnimeKaizoku/6176165"
	b, err := client.BanUser(1478, reason01, &sibylSystemGo.BanConfig{
		Message: msg,
		SrcUrl:  src,
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(b)
}

func TestBanUser02(t *testing.T) {
	token := getToken()
	if len(token) == 0 {
		log.Println("token is empty; exiting")
	}

	config := &sibylSystemGo.SibylConfig{
		HostUrl:    "https://psychopass.animekaizoku.com",
		HttpClient: http.DefaultClient,
	}

	client := sibylSystemGo.NewClient(token, config)
	//print(client)
	fmt.Print(client)

	const reason01 = "Spam adding +99 members to an anime group"
	const msg = "https://t.me/telegram/832471"
	const src = "https://t.me/AnimeKaizoku/6176165"
	b, err := client.BanUser(1478, reason01, &sibylSystemGo.BanConfig{
		Message: msg,
		SrcUrl:  src,
	})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(b)
}

func TestUnBanUser01(t *testing.T) {
	token := getToken()
	if len(token) == 0 {
		log.Println("token is empty; exiting")
	}

	client := sibylSystemGo.NewClient(token, sibylSystemGo.GetDefaultConfig())
	b, err := client.RemoveBan(1478, "", nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(b)
}
