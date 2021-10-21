package cli

import (
	"context"
	"flag"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"log"
)

var auth = spotifyauth.New(spotifyauth.WithRedirectURL("http://localhost:3000/login_check"))

func main() {
	code := flag.String("code", "", "authorization code to negotiate by token")
	flag.Parse()

	if *code == "" {
		log.Fatal("code required")
	}
	if err := Authorize(*code); err != nil {
		log.Fatal("error while negotiating the token: ", err)
	}
}

func Authorize(code string) error {
	ctx := context.Background()
	token, err := auth.Exchange(ctx, code)
	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	client := spotify.New(httpClient)

	user, err := client.CurrentUser(ctx)
	if err != nil {
		return err
	}
	log.Printf("Logged in as %s\n", user.DisplayName)

	return nil
}
