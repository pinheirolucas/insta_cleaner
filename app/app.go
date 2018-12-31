package app

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/ahmdrz/goinsta"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/pinheirolucas/insta_cleaner/cleaner"
	"github.com/pinheirolucas/insta_cleaner/logger"
)

func init() {
	http.HandleFunc("/unfollow/morning", unfollow)
	http.HandleFunc("/unfollow/night", unfollow)
}

func unfollow(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	username := os.Getenv("INSTAGRAM_USERNAME")
	password := os.Getenv("INSTAGRAM_PASSWORD")
	project := os.Getenv("GCLOUD_PROJECT_ID")

	ustr := fmt.Sprintf("%s@%s", username, password)
	usha := sha1.New()
	usha.Write([]byte(ustr))
	session := fmt.Sprintf("%x", usha.Sum(nil))

	object, err := getSession(ctx, project, session)
	if err != nil {
		log.Errorf(ctx, "getSession: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var insta *goinsta.Instagram

	if rc, err := object.NewReader(ctx); err == nil {
		insta, err = goinsta.ImportReader(rc)
		if err != nil {
			log.Errorf(ctx, "goinsta.ImportReader: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if err == storage.ErrObjectNotExist {
		insta = goinsta.New(username, password)

		if err := insta.Login(); err != nil {
			log.Errorf(ctx, "(*goinsta.Instagram).Login: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		wc := object.NewWriter(ctx)
		defer wc.Close()

		if err := goinsta.Export(insta, wc); err != nil {
			log.Errorf(ctx, "goinsta.Export: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		log.Errorf(ctx, "(*storage.ObjectHandle).NewReader: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fcli, err := firestore.NewClient(ctx, project)
	if err != nil {
		log.Errorf(ctx, "firestore.NewClient: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	instagramService := cleaner.NewGoinstaInstagramService(insta)
	whitelistService := cleaner.NewFirebaseWhitelistService(fcli)
	l := logger.NewAppengine(ctx)
	service := cleaner.NewService(instagramService, whitelistService, l)

	if err := service.Clean(); err != nil {
		log.Errorf(ctx, "(*cleaner.Service).Clean: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getSession(ctx context.Context, bucket, store string) (*storage.ObjectHandle, error) {
	scli, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "storage.NewClient")
	}

	return scli.Bucket(bucket).Object(store), nil
}