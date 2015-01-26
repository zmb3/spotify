package spotify

import (
	"net/http"
	"testing"
)

const userResponse = `
{
  "display_name" : "Ronald Pompa",
  "external_urls" : {
    "spotify" : "https://open.spotify.com/user/wizzler"
    },
    "followers" : {
      "href" : null,
      "total" : 3829
    },
    "href" : "https://api.spotify.com/v1/users/wizzler",
    "id" : "wizzler",
    "images" : [ {
      "height" : null,
      "url" : "http://profile-images.scdn.co/images/userprofile/default/9d51820e73667ea5f1e97ea601cf0593b558050e",
      "width" : null
    } ],
    "type" : "user",
    "uri" : "spotify:user:wizzler"
}`

func TestUserProfile(t *testing.T) {
	client := testClientString(http.StatusOK, userResponse)
	user, err := client.UserPublicProfile("wizzler")
	if err != nil {
		t.Error(err)
		return
	}
	if user.ID != "wizzler" {
		t.Error("Expected user wizzler, got ", user.ID)
	}
	if f := user.Followers.Count; f != 3829 {
		t.Errorf("Expected 3829 followers, got %d\n", f)
	}
}

func TestCurrentUser(t *testing.T) {
	json := `{
		"country" : "US",
		"display_name" : null,
		"email" : "username@domain.com",
		"external_urls" : {
			"spotify" : "https://open.spotify.com/user/username"
		},
		"followers" : {
			"href" : null,
			"total" : 0
		},
		"href" : "https://api.spotify.com/v1/users/userame",
		"id" : "username",
		"images" : [ ],
		"product" : "premium",
		"type" : "user",
		"uri" : "spotify:user:username"
	}`
	client := testClientString(http.StatusOK, json)
	addDummyAuth(client)

	me, err := client.CurrentUser()
	if err != nil {
		t.Error(err)
		return
	}
	if me.Country != CountryUSA ||
		me.Email != "username@domain.com" ||
		me.Product != "premium" {
		t.Error("Received incorrect response")
	}
}
