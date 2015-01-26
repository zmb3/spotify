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

func TestFollowUsersMissingScope(t *testing.T) {
	json := `{
		"error": {
			"status": 403,
			"message": "Insufficient client scope"
		}
	}`
	client := testClientString(http.StatusForbidden, json)
	addDummyAuth(client)

	err := client.Follow(ID("exampleuser01"))
	if serr, ok := err.(Error); !ok {
		t.Error("Expected insufficient client scope error")
	} else {
		if serr.Status != http.StatusForbidden {
			t.Error("Expected HTTP 403")
		}
	}
}

func TestFollowUsersInvalidToken(t *testing.T) {
	json := `{
		"error": {
			"status": 401,
			"message": "Invalid access token"
		}
	}`
	client := testClientString(http.StatusUnauthorized, json)
	addDummyAuth(client)

	err := client.Follow(ID("dummyID"))
	if serr, ok := err.(Error); !ok {
		t.Error("Expected invalid token error")
	} else {
		if serr.Status != http.StatusUnauthorized {
			t.Error("Expected HTTP 401")
		}
	}
}

func TestUserFollows(t *testing.T) {
	json := "[ false, true ]"
	client := testClientString(http.StatusOK, json)
	addDummyAuth(client)
	follows, err := client.UserFollows("artist", ID("74ASZWbe4lXaubB36ztrGX"), ID("08td7MxkoHQkXnWAYD8d6Q"))
	if err != nil {
		t.Error(err)
		return
	}
	if len(follows) != 2 || follows[0] || !follows[1] {
		t.Error("Incorrect result", follows)
	}
}
