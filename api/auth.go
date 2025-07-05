package api

import (
	"context"
	"encoding/json"
	"go-drive/ui"
	"net/http"
	"os"

	"github.com/charmbracelet/log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func GetGoogleDriveService(ctx context.Context, credentials_path, token_path, scope string, allow_web_auth bool) (*drive.Service, APIError) {
	if client, err := getGoogleClient(ctx, credentials_path, token_path, scope, allow_web_auth); err != nil {
		return nil, err
	} else {
		if srv, err := drive.NewService(ctx, option.WithHTTPClient(client)); err != nil {
			return nil, &GoogleDriveError{err}
		} else {
			log.Debug("Google Drive service initialized successfully")
			return srv, nil
		}
	}
}

func getGoogleClient(ctx context.Context, credentials_path, token_path, scope string, allow_web_auth bool) (*http.Client, APIError) {
	if credentials, err := os.ReadFile(credentials_path); err != nil {
		return nil, &ReadFileFailed{err, credentials_path}
	} else {
		if config, err := google.ConfigFromJSON(credentials, scope); err != nil {
			return nil, &ParseCredentialsFailed{err}
		} else {
			if client, err := getClient(ctx, token_path, config, allow_web_auth); err != nil {
				return nil, err
			} else {
				return client, nil
			}
		}
	}
}

func getClient(ctx context.Context, token_path string, config *oauth2.Config, allow_web_auth bool) (*http.Client, APIError) {
	tok, err := tokenFromFile(token_path)
	if err != nil {
		if allow_web_auth {
			if tok, err = getTokenFromWeb(ctx, config); err != nil {
				return nil, err
			}
			if err = saveToken(token_path, tok); err != nil {
				return nil, err
			}
		} else {
			log.Debugf("unable to authenticate with token file %s", token_path)
			log.Print("run `go-drive a` to authenticate via web")
			return nil, err
		}
	}
	return config.Client(ctx, tok), nil
}

func tokenFromFile(file string) (*oauth2.Token, APIError) {
	if f, err := os.Open(file); err != nil {
		return nil, &OpenFileFailed{err, file}
	} else {
		defer f.Close()
		tok := &oauth2.Token{}
		if err := json.NewDecoder(f).Decode(tok); err != nil {
			return nil, &TokenDecodeFailed{err}
		} else {
			return tok, nil
		}
	}
}

func getTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, APIError) {
	auth_url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("go to the following link in your browser then type the authorization code: \n\n%v\n", auth_url)
	if auth_code, err := ui.Prompt("auth code", ""); err != nil {
		return nil, &ParseTokenFailed{err}
	} else if tok, err := config.Exchange(ctx, auth_code); err != nil {
		return nil, &AuthTokenFailed{err, auth_code}
	} else {
		return tok, nil
	}
}

func saveToken(path string, token *oauth2.Token) APIError {
	log.Printf("Saving credential file to: %s\n", path)
	if f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return &WriteTokenFailed{err, path}
	} else {
		defer f.Close()
		json.NewEncoder(f).Encode(token)
		return nil
	}
}
