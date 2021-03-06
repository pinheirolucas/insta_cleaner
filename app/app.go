package app

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/ahmdrz/goinsta"
	"github.com/pkg/errors"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"

	"github.com/pinheirolucas/insta_cleaner/cleaner"
	"github.com/pinheirolucas/insta_cleaner/logger"
	"github.com/pinheirolucas/insta_cleaner/whitelist"
)

func init() {
	http.HandleFunc("/tasks/unfollow", unfollow)
}

func unfollow(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	username := os.Getenv("INSTAGRAM_USERNAME")
	password := os.Getenv("INSTAGRAM_PASSWORD")
	bucket := os.Getenv("FIREBASE_STORAGE_BUCKET")
	databaseURL := os.Getenv("FIREBASE_REALTIME_DATABASE_URL")

	maxUnfollows, err := strconv.Atoi(os.Getenv("INSTA_CLEANER_MAX_UNFOLLOWS"))
	if err != nil {
		maxUnfollows = 10
	}

	ustr := fmt.Sprintf("%s@%s", username, password)
	usha := sha1.New()
	usha.Write([]byte(ustr))
	session := fmt.Sprintf("%x", usha.Sum(nil))

	app, err := firebase.NewApp(ctx, &firebase.Config{
		DatabaseURL:   databaseURL,
		StorageBucket: bucket,
	})
	if err != nil {
		log.Errorf(ctx, "firebase.NewApp: %v \n", err)
	}

	object, err := getSession(ctx, app, session)
	if err != nil {
		log.Errorf(ctx, "getSession: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var insta *goinsta.Instagram
	instaOptions := []goinsta.Option{
		goinsta.WithHTTPClient(urlfetch.Client(ctx)),
	}

	if rc, err := object.NewReader(ctx); err == nil {
		log.Infof(ctx, "using existing session: %s", session)

		insta, err = goinsta.ImportReader(rc, instaOptions...)
		if err != nil {
			log.Errorf(ctx, "goinsta.ImportReader: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if err == storage.ErrObjectNotExist {
		log.Infof(ctx, "creating new instagram session: %s", session)

		insta = goinsta.New(username, password, instaOptions...)

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

	whitelistService, err := whitelist.NewService(ctx, os.Getenv("WHITELIST_SERVICE"), app)
	if err != nil {
		log.Errorf(ctx, "whitelist.NewService: %v \n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	instagramService := cleaner.NewGoinstaInstagramService(insta)
	l := logger.NewAppengine(ctx)
	service := cleaner.NewService(instagramService, whitelistService, l, cleaner.WithMaxUnfollows(uint32(maxUnfollows)))

	if err := service.Clean(); err != nil {
		log.Errorf(ctx, "(*cleaner.Service).Clean: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getSession(ctx context.Context, app *firebase.App, store string) (*storage.ObjectHandle, error) {
	scli, err := app.Storage(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "(*firebase.App).Storage")
	}

	bucket, err := scli.DefaultBucket()
	if err != nil {
		return nil, errors.Wrap(err, "(*storage.Client).DefaultBucket")
	}

	return bucket.Object(store), nil
}
