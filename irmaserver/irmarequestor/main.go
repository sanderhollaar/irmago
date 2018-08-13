package irmarequestor

import (
	"io/ioutil"
	"net/http"

	"github.com/privacybydesign/irmago"
	"github.com/privacybydesign/irmago/irmaserver"
	"github.com/privacybydesign/irmago/irmaserver/backend"
)

type SessionHandler func(*irmaserver.SessionResult)

var handlers = make(map[string]SessionHandler)

func Initialize(configuration *irmaserver.Configuration) error {
	return backend.Initialize(configuration)
}

func StartSession(request irma.SessionRequest, handler SessionHandler) (*irma.Qr, string, error) {
	qr, token, err := backend.StartSession(request)
	if err != nil {
		return nil, "", err
	}
	if handler != nil {
		handlers[token] = handler
	}
	return qr, token, nil
}

func GetSessionResult(token string) *irmaserver.SessionResult {
	return backend.GetSessionResult(token)
}

func HttpHandlerFunc(prefix string) http.HandlerFunc {
	if len(prefix) != 0 && prefix[0] != '/' {
		prefix = "/" + prefix
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var message []byte
		message, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		path := r.URL.Path[len(prefix):]
		status, response, result := backend.HandleProtocolMessage(path, r.Method, r.Header, message)
		w.WriteHeader(status)
		w.Write(response)
		if result != nil {
			if handler, ok := handlers[result.Token]; ok {
				go handler(result)
			}
		}
	}
}
