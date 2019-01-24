package app

import (
	"encoding/json"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func following(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	_, ok := ctx.Value(sk{}).(string)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users := []*user{
		&user{
			ID:          1,
			Name:        "Leandro Neko",
			Username:    "leandroneko",
			ProfileURL:  "https://www.instagram.com/leandroneko/",
			ProfilePic:  "https://instagram.fjoi2-1.fna.fbcdn.net/vp/1227b756b3511c008b4653d41c8273bd/5CC03042/t51.2885-19/s150x150/47690721_762539200793315_2621319472480256000_n.jpg?_nc_ht=instagram.fjoi2-1.fna.fbcdn.net",
			IsFollowing: true,
		},
		&user{
			ID:          2,
			Name:        "⚡️PITTY⚡️",
			Username:    "pitty",
			ProfileURL:  "https://www.instagram.com/pitty/",
			ProfilePic:  "https://instagram.fjoi2-1.fna.fbcdn.net/vp/b35b26b4dc3c5dcd68d674a8a0f699d2/5CCFBF39/t51.2885-19/s150x150/36149132_190177528500343_3582746510120452096_n.jpg?_nc_ht=instagram.fjoi2-1.fna.fbcdn.net",
			IsFollowing: true,
		},
		&user{
			ID:          3,
			Name:        "Ana Abelha",
			Username:    "euanaabelha",
			ProfileURL:  "https://www.instagram.com/euanaabelha/",
			ProfilePic:  "https://instagram.fjoi2-1.fna.fbcdn.net/vp/470ef8ec63ad46d5dbe12f9c5adbdb19/5CBFEEF7/t51.2885-19/s150x150/42890822_188202832072110_1885326991805120512_n.jpg?_nc_ht=instagram.fjoi2-1.fna.fbcdn.net",
			IsFollowing: true,
		},
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Errorf(ctx, "unable to marshal response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type user struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Username    string `json:"username,omitempty"`
	ProfileURL  string `json:"profileUrl,omitempty"`
	ProfilePic  string `json:"profilePic,omitempty"`
	IsFollowing bool   `json:"isFollowing,omitempty"`
}
