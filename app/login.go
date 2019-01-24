package app

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/ahmdrz/goinsta"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

type sk struct{}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)

		auth := r.Header.Get("Authorization")

		if strings.TrimSpace(auth) == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		const prefix = "Basic "
		if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sessionKey := auth[len(prefix):]

		bucket := os.Getenv("FIREBASE_STORAGE_BUCKET")
		app, err := firebase.NewApp(ctx, &firebase.Config{
			StorageBucket: bucket,
		})
		if err != nil {
			log.Errorf(ctx, "firebase.NewApp: %v \n", err)
			return
		}

		object, err := getSession(ctx, app, sessionKey)
		if err != nil {
			log.Errorf(ctx, "getSession: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := object.NewReader(ctx); err == storage.ErrObjectNotExist {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		c := appengine.WithContext(context.WithValue(ctx, sk{}, sessionKey), r)

		next.ServeHTTP(w, r.WithContext(c))
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	if r.Body == nil {
		log.Warningf(ctx, "nil request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warningf(ctx, "unable to unmarshal body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Username) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bucket := os.Getenv("FIREBASE_STORAGE_BUCKET")
	ustr := fmt.Sprintf("%s@%s", req.Username, req.Password)
	usha := sha1.New()
	usha.Write([]byte(ustr))
	sessionKey := fmt.Sprintf("%x", usha.Sum(nil))

	app, err := firebase.NewApp(ctx, &firebase.Config{
		StorageBucket: bucket,
	}, option.WithCredentialsFile("C:\\Users\\Lucas Pinheiro\\works\\src\\github.com\\pinheirolucas\\insta_cleaner\\keys\\insta-cleaner-227321-firebase-adminsdk-kcv6c-df23818d6d.json"))
	if err != nil {
		log.Errorf(ctx, "firebase.NewApp: %v \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session, err := getSession(ctx, app, sessionKey)
	if err != nil {
		log.Errorf(ctx, "getSession: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = session.NewReader(ctx)
	switch err {
	case nil:
		// continue
	case storage.ErrObjectNotExist:
		insta := goinsta.New(req.Username, req.Password, goinsta.WithHTTPClient(urlfetch.Client(ctx)))

		if err := insta.Login(); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		wc := session.NewWriter(ctx)
		defer wc.Close()

		if err := goinsta.Export(insta, wc); err != nil {
			log.Errorf(ctx, "goinsta.Export: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		log.Errorf(ctx, "(*storage.ObjectHandle).NewReader: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(loginResponse{SessionKey: sessionKey}); err != nil {
		log.Errorf(ctx, "unable to marshal response: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type loginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type loginResponse struct {
	SessionKey string `json:"sessionKey,omitempty"`
}
