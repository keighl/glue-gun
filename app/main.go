package api

import (
	"net/http"
	"encoding/json"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"golang.org/x/net/context"
	"goji.io"
    "goji.io/pat"
    "github.com/rs/cors"
)

func init() {
	mux := goji.NewMux()

	corz := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders: []string{"*", "x-requested-with", "Content-Type", "If-Modified-Since", "If-None-Match"},
		ExposedHeaders: []string{"Content-Length"},
	})

	mux.Use(corz.Handler)
	mux.UseC(Recoverer)

	mux.HandleFuncC(pat.Get("/fitbit/auth"), FitbitAuth)
	mux.HandleFuncC(pat.Get("/fitbit/auth/callback"), FitbitAuthCallback)
	mux.HandleFuncC(pat.Get("/fitbit/playground"), FitbitPlayground)
	mux.HandleFuncC(pat.Post("/twilio/inbound"), TwilioInbound)

	mux.HandleFunc(pat.Get("/*"), NotFound)

	http.Handle("/", mux)
}

func Recoverer(inner goji.Handler) goji.Handler {
	mw := func(c context.Context, w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				ctx := appengine.NewContext(req)
				log.Errorf(ctx, "%v", err)
				render(w, 500, errorsJSON("There was an unexpected error"))
			}
		}()
		inner.ServeHTTPC(c, w, req)
	}

	return goji.HandlerFunc(mw)
}

func Authenticator(inner goji.Handler) goji.Handler {

	mw := func(c context.Context, w http.ResponseWriter, req *http.Request) {

		ctx := appengine.NewContext(req)
		key := pat.Param(ctx, "key")
		if key != conf.APIKey {
			log.Errorf(ctx, "Wrong key %s", key)
			render(w, 401, errorsJSON("Invalid JWT"))
			return
		}

		inner.ServeHTTPC(c, w, req)
	}
	return goji.HandlerFunc(mw)
}

func render(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")
	w.WriteHeader(code)

	if data != nil {
		d, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(d)
	}
}

type jsonPayload map[string]interface{}
type ErrorResource struct {
	Title string `json:"title"`
}

func errorsJSON(data ...string) jsonPayload {
	errors := make([]ErrorResource, len(data))
	for i, e := range data {
		errors[i] = ErrorResource{Title: e}
	}
	return jsonPayload{"errors": errors}
}

func messagesJSON(messages ...string) jsonPayload {
	return jsonPayload{"messages": messages}
}
