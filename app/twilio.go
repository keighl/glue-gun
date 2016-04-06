package api

import (
	"net/http"
	"encoding/xml"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"golang.org/x/net/context"
	"github.com/mholt/binding"
	r "github.com/unrolled/render"
)

func NotFound(w http.ResponseWriter, req *http.Request) {
	render(w, 404, errorsJSON("Not found"))
}

type TwilioSMSRequestData struct {
	MessageSid          string
	SmsSid              string
	AccountSid          string
	MessagingServiceSid string
	From                string
	To                  string
	Body                string
	NumMedia            int
}

// Then provide a field mapping (pointer receiver is vital)
func (d *TwilioSMSRequestData) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&d.MessageSid:          "MessageSid",
		&d.SmsSid:              "SmsSid",
		&d.AccountSid:          "AccountSid",
		&d.MessagingServiceSid: "MessagingServiceSid",
		&d.From:                "From",
		&d.To:                  "To",
		&d.Body:                "Body",
		&d.NumMedia:            "NumMedia",
	}
}

type TwiML struct {
	XMLName xml.Name `xml:"Response"`
	Message string   `xml:",omitempty"`
}

// https://www.twilio.com/docs/api/twiml/sms/twilio_request#twilio-data-passing
func TwilioInbound(c context.Context, w http.ResponseWriter, req *http.Request) {
	ctx := appengine.NewContext(req)

	data := new(TwilioSMSRequestData)
	errs := binding.Bind(req, data)
	if errs.Handle(w) {
		log.Errorf(ctx, errs.Error())
		r.New().XML(w, 500, TwiML{})
		return
	}

	log.Infof(ctx, "MessageSid: %s", data.MessageSid)
	log.Infof(ctx, "SmsSid: %s", data.SmsSid)
	log.Infof(ctx, "AccountSid: %s", data.AccountSid)
	log.Infof(ctx, "MessagingServiceSid: %s", data.MessagingServiceSid)
	log.Infof(ctx, "From: %s", data.From)
	log.Infof(ctx, "To: %s", data.To)
	log.Infof(ctx, "Body: %s", data.Body)
	log.Infof(ctx, "NumMedia: %d", data.NumMedia)

	r.New().XML(w, http.StatusOK, TwiML{})
}

