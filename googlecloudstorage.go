package rule110

import (
  "net/http"

  "golang.org/x/net/context"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
  "google.golang.org/cloud"
  "google.golang.org/appengine"
  "google.golang.org/appengine/urlfetch"
  "google.golang.org/cloud/storage"
)

const bucket = "rule110-go.appspot.com"

func cloudAuthContext(r *http.Request) (context.Context, error) {
  c := appengine.NewContext(r)

  hc := &http.Client{
      Transport: &oauth2.Transport{
          Source: google.AppEngineTokenSource(c, storage.ScopeFullControl),
          Base:   &urlfetch.Transport{Context: c},
      },
  }
  return cloud.WithContext(c, "rule110-go", hc), nil

}

