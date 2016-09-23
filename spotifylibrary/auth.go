package spotifylibrary

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/unixpickle/ezserver"
)

const (
	ClientIDEnvVar     = "SPOTIFY_CLIENT_ID"
	CallbackPathEnvVar = "SPOTIFY_CALLBACK_PATH"
	PortsEnvVar        = "SPOTIFY_PORTS"
)

// These are default parameters used for spotify's
// "Implicit Grant Flow" authorization.
const (
	DefaultAuthClientID     = "dac4cf3243db4bd1b055efc7af84be72"
	DefaultAuthCallbackPath = "/spotifydone"
	DefaultAuthPorts        = "14505,19535,19548"
)

type authResult struct {
	Token string
	Error error
}

// Auth authenticates with Spotify and returns an access
// token for the spotify API.
//
// This will look at three environment variables that may
// override default parameters:
//
// The SPOTIFY_CLIENT_ID variable may specify the ID for
// a spotify application to use for the request.
// By default, this is DefaultAuthClientID.
//
// The SPOTIFY_CALLBACK_PATH variable may specify the path
// of the callback URL (after the http://localhost:port
// part of the URL).
// By default, this is DefaultAuthCallbackPath.
//
// The SPOTIFY_PORTS variable may specify a comma-separated
// list of ports on localhost which are whitelisted as
// redirect URIs.
// By default, this is DefaultAuthPorts.
func Auth() (string, error) {
	serverRes, port, err := serveCallbackURL()
	if err != nil {
		return "", err
	}

	query := url.Values{}
	query.Add("client_id", clientID())
	query.Add("response_type", "token")
	query.Add("redirect_uri", "http://localhost:"+strconv.Itoa(port)+callbackPath())
	query.Add("scope", "user-library-read")

	u := url.URL{
		Scheme:   "https",
		Host:     "accounts.spotify.com",
		Path:     "/authorize",
		RawQuery: query.Encode(),
	}

	fmt.Println("Please navigate to:")
	fmt.Println(u.String())

	result := <-serverRes
	return result.Token, result.Error
}

func serveCallbackURL() (<-chan authResult, int, error) {
	ports, err := authPorts()
	if err != nil {
		return nil, 0, err
	}

	resChan := make(chan authResult, 1)
	var server *ezserver.HTTP
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, callbackPath()) {
			return
		}
		if r.URL.RawQuery == "" {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<!doctype html><html><head><script>" +
				"window.location.search = window.location.hash.substr(1);" +
				"</script></head><body>Authorizing...</body></html>"))
			return
		}

		w.Header().Set("Content-Type", "text/plain")

		res := authResult{Token: r.FormValue("access_token")}
		if errStr := r.FormValue("error"); errStr != "" {
			res.Error = errors.New("authorization error: " + errStr)
			w.Write([]byte("Authorization error: " + errStr + "."))
		} else {
			w.Write([]byte("Authorization successful."))
		}

		select {
		case resChan <- res:
		default:
		}

		// If we stop the server right away, Safari complains
		// that it cannot connect to the server, even though it
		// obviously did if we got a token.
		go func() {
			time.Sleep(time.Second)
			server.Stop()
		}()
	})

	for _, port := range ports {
		server = ezserver.NewHTTP(handler)
		if err := server.Start(port); err == nil {
			return resChan, port, nil
		}
	}

	return nil, 0, errors.New("could not listen on any callback ports")
}

func clientID() string {
	if s := os.Getenv(ClientIDEnvVar); s != "" {
		return s
	}
	return DefaultAuthClientID
}

func callbackPath() string {
	if s := os.Getenv(CallbackPathEnvVar); s != "" {
		return s
	}
	return DefaultAuthCallbackPath
}

func authPorts() ([]int, error) {
	unsplit := DefaultAuthPorts
	if s := os.Getenv(PortsEnvVar); s != "" {
		unsplit = s
	}
	var res []int
	for _, portStr := range strings.Split(unsplit, ",") {
		num, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, errors.New("invalid port number: " + portStr)
		}
		res = append(res, num)
	}
	return res, nil
}
