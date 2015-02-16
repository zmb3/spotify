// Copyright 2014, 2015 Zac Bergquist
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spotify

import (
	"fmt"
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
	user, err := client.GetUsersPublicProfile("wizzler")
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
	follows, err := client.CurrentUserFollows("artist", ID("74ASZWbe4lXaubB36ztrGX"), ID("08td7MxkoHQkXnWAYD8d6Q"))
	if err != nil {
		t.Error(err)
		return
	}
	if len(follows) != 2 || follows[0] || !follows[1] {
		t.Error("Incorrect result", follows)
	}
}

func TestCurrentUsersTracks(t *testing.T) {
	client := testClientFile(http.StatusOK, "test_data/current_users_tracks.txt")
	addDummyAuth(client)
	tracks, err := client.CurrentUsersTracks()
	if err != nil {
		t.Error(err)
		return
	}
	if tracks.Limit != 20 {
		t.Errorf("Expected limit 20, got %d\n", tracks.Limit)
	}
	if tracks.Endpoint != "https://api.spotify.com/v1/me/tracks?offset=0&limit=20" {
		t.Error("Endpoint incorrect")
	}
	if tracks.Total != 3 {
		t.Errorf("Expect 3 results, got %d\n", tracks.Total)
		return
	}
	if len(tracks.Tracks) != tracks.Total {
		t.Error("Didn't get expected number of results")
		return
	}
	expected := "You & I (Nobody In The World)"
	if tracks.Tracks[0].Name != expected {
		t.Errorf("Expected '%s', got '%s'\n", expected, tracks.Tracks[0].Name)
		fmt.Printf("\n%#v\n", tracks.Tracks[0])
	}
}
