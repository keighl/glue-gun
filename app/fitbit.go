package api

import (
	"net/http"
	"encoding/json"
	"io/ioutil"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func FitbitAuth(c context.Context, w http.ResponseWriter, req *http.Request) {
    url := fitbitOauthConf.AuthCodeURL("state")
    http.Redirect(w, req, url, 302)
}

func FitbitAuthCallback(c context.Context, w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)
	code := req.FormValue("code")
    token, err := fitbitOauthConf.Exchange(ctx, code)
    if err != nil {
       	log.Errorf(ctx, err.Error())
       	render(w, 200, errorsJSON(err.Error()))
       	return
    }

    render(w, 200, token)
}

func FitbitPlayground(c context.Context, w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	token := &oauth2.Token{}
	err := json.Unmarshal([]byte(fitbitToken), token)
	if err != nil {
        log.Errorf(ctx, err.Error())
       	render(w, 200, errorsJSON(err.Error()))
       	return
    }

    client := fitbitOauthConf.Client(ctx, token)
    res, err := client.Get("https://api.fitbit.com/1/user/-/activities/steps/date/today/1m.json")
    if err != nil {
        log.Errorf(ctx, err.Error())
       	render(w, 200, errorsJSON(err.Error()))
       	return
    }

    defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")
	w.WriteHeader(200)
	w.Write(body)
}

