package tests

import (
	"log"
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
	const msg = "https://t.me/MusicSingAlong/832471"
	const src = "https://t.me/AnimeKaizoku/6176165"
	b, err := client.BanUser(1478, reason01, msg, src)
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
	b, err := client.RemoveBan(1478)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(b)
}
