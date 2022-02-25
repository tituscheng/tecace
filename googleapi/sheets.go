package googleapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	SPREADSHEET_ID   = "SPREADSHEET_ID"
	READ_WRITE_RANGE = "Sheet1!A1:B"
	RAW              = "RAW"
)

var srv = createSheet()

func Get() (map[string]string, error) {
	data := make(map[string]string)
	resp, err := srv.Spreadsheets.Values.Get(os.Getenv(SPREADSHEET_ID), READ_WRITE_RANGE).Do()
	if err != nil {
		return data, err
	}
	for _, row := range resp.Values {
		if len(row) > 0 {
			data[row[0].(string)] = row[1].(string)
		}
	}
	return data, nil
}

func Post(key, val string) error {
	var vr sheets.ValueRange
	vr.Values = append(vr.Values, []interface{}{key, val})

	valueInputOption := RAW
	_, err := srv.Spreadsheets.Values.Append(os.Getenv(SPREADSHEET_ID), READ_WRITE_RANGE, &vr).ValueInputOption(valueInputOption).Do()
	if err != nil {
		return err
	}
	return nil
}

func Update(key, val string) error {

	var vr sheets.ValueRange
	vr.Values = append(vr.Values, []interface{}{key, val})

	// How the input data should be interpreted
	valueInputOption := RAW

	ctx := context.Background()

	_, err := srv.Spreadsheets.Values.Update(os.Getenv(SPREADSHEET_ID), READ_WRITE_RANGE, &vr).ValueInputOption(valueInputOption).Context(ctx).Do()
	if err != nil {
		return err
	}
	return nil
}

func Delete(key string) error {
	ctx := context.Background()

	resp, err := srv.Spreadsheets.Values.Get(os.Getenv(SPREADSHEET_ID), READ_WRITE_RANGE).Do()
	if err != nil {
		return err
	}

	index := 0
	for idx, row := range resp.Values {
		if len(row) > 0 && row[0] == key {
			index = idx + 1
		}
	}

	deleteRange := fmt.Sprintf("Sheet1!A%d:B%d", index, index)

	rb := &sheets.ClearValuesRequest{
		// TODO: Add desired fields of the request body.
	}
	_, err = srv.Spreadsheets.Values.Clear(os.Getenv(SPREADSHEET_ID), deleteRange, rb).Context(ctx).Do()
	if err != nil {
		return err
	}
	return nil
}

func createSheet() *sheets.Service {
	var err error
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		// return c.JSON(response.Response{Result: 500, Description: "problem accessing google api"})
		return nil
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		// return c.JSON(response.Response{Result: 500, Description: "problem accessing google api"})
		return nil
	}

	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		// return c.JSON(response.Response{Result: 500, Description: "unable to read data"})
		return nil
	}
	return srv
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
