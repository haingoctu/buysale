package sheets

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
		"buysale/parse"

        "golang.org/x/net/context"
        "golang.org/x/oauth2"
        "golang.org/x/oauth2/google"
        "google.golang.org/api/sheets/v4"
		"google.golang.org/api/drive/v3"
		"buysale/mylog"
)

const COL_IN_SHEETS int = 10
var Fields = [COL_IN_SHEETS]string {
	"User",
	"Make",		
	"Size",
	"Leather",
	"Sole",
	"Price",
	"Condition",
	"ImageLink",
	"Notes",
	"PermLink",
}

var logger = mylog.GetInstance()

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
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
        logger.Print("Go to the following link in your browser then type the "+
                "authorization code: \n%v\n", authURL)

        var authCode string
        if _, err := fmt.Scan(&authCode); err != nil {
                logger.Fatalf("Unable to read authorization code: %v", err)
        }

        tok, err := config.Exchange(oauth2.NoContext, authCode)
        if err != nil {
                logger.Fatalf("Unable to retrieve token from web: %v", err)
        }
        return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
        f, err := os.Open(file)
        defer f.Close()
        if err != nil {
                return nil, err
        }
        tok := &oauth2.Token{}
        err = json.NewDecoder(f).Decode(tok)
        return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
        logger.Print("Saving credential file to: %s\n", path)
        f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
        defer f.Close()
        if err != nil {
                logger.Fatalf("Unable to cache oauth token: %v", err)
        }
        json.NewEncoder(f).Encode(token)
}

func Create(t string) string {
        b, err := ioutil.ReadFile("credentials.json")
        if err != nil {
                logger.Fatalf("Unable to read client secret file: %v", err)
        }

        // If modifying these scopes, delete your previously saved token.json.
        config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets", "https://www.googleapis.com/auth/drive")
        if err != nil {
                logger.Fatalf("Unable to parse client secret file to config: %v", err)
        }
        client := getClient(config)

        srv, err := sheets.New(client)
        if err != nil {
                logger.Fatalf("Unable to retrieve Sheets client: %v", err)
        }
	
		spreadsheetId := "1pDNm-5Du6uG0hU3_qdcjTkrU8bZYhm3f2eaaLcG2GME"
		destinationSpreadsheetId := ""
		
		rb := &sheets.Spreadsheet{
			//EMPTY
        }
		//Create new sheet
		newSheet, err := srv.Spreadsheets.Create(rb).Do()
        if err != nil {
                logger.Fatal(err)
        }
		
		destinationSpreadsheetId = fmt.Sprint(newSheet.SpreadsheetId)

		rb2 := &sheets.CopySheetToAnotherSpreadsheetRequest{
                DestinationSpreadsheetId: destinationSpreadsheetId,
        }
		
		//copy our template sheet to new sheet
        _, err = srv.Spreadsheets.Sheets.CopyTo(spreadsheetId, 0, rb2).Do()
        if err != nil {
                logger.Fatalf("Unable to retrieve data from sheet: %v", err)
        }

		sheetProp, err := srv.Spreadsheets.Get(destinationSpreadsheetId).Do()
		if err != nil {
				logger.Fatalf("Unable to retrieve data from sheet: %v", err)
		}
		
		//Req Body: delete sheet request start
		val := sheetProp.Sheets[0].Properties.SheetId
		
		deleteSheetRequest := &sheets.DeleteSheetRequest{
			SheetId: val,
		}
		//delete sheet request end
		
		//Req Body: update the title of the spreadsheet start
		spreadSheetProperties := &sheets.SpreadsheetProperties{
			Title: t,
		}
		updateSSheet := &sheets.UpdateSpreadsheetPropertiesRequest{
			Properties: spreadSheetProperties,
			Fields: "title",
		}
		//Req Body: update the title of the spreadsheet end
		
		/*/Req Body: update sheet title start
		sheetProperties := &sheets.SheetProperties{
			Title: "Main - Contact haingoctu@gmail.com",
			SheetId: 0,
		}
		
		updateSheet := &sheets.UpdateSheetPropertiesRequest{
			Properties: sheetProperties,
			Fields: "title",
		}
		//Req Body: update sheet title end*/
		
		requests := []*sheets.Request{
			{DeleteSheet: deleteSheetRequest},
			{UpdateSpreadsheetProperties: updateSSheet},
			
		}
		rb3 := &sheets.BatchUpdateSpreadsheetRequest{
                Requests: requests,
        }

		resp2, err := srv.Spreadsheets.BatchUpdate(destinationSpreadsheetId, rb3).Do()
		if err != nil {
			logger.Fatalf("Unable to retrieve data from sheet: %v", err)
		}
		
		logger.Print("Created: ", resp2, destinationSpreadsheetId, "\n")
		return destinationSpreadsheetId
}


func PostToSheets(items map[int]*parse.Shoe, sheetID string) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
			logger.Fatalf("Unable to read client secret file: %v", err)
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
			logger.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.New(client)
	if err != nil {
			logger.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	// The A1 notation of a range to search for a logical table of data.
	// Values will be appended after the last row of the table.
	range2 := "A2" // TODO: Update placeholder value.

	// How the input data should be interpreted.
	valueInputOption := "RAW" // TODO: Update placeholder value.

	// How the input data should be inserted.
	insertDataOption := "INSERT_ROWS" // TODO: Update placeholder value.
	
	rb := &sheets.ValueRange{
	}
	var shoeVals [][]interface{}
	for _, shoes := range items {
		logger.Print(shoes.Make)
		for _, field := range Fields {
			valz := shoes.Reflectz(field)
			//logger.Print(valz)
			shoeVals = append(shoeVals, []interface{}{valz})
		}
		rb = &sheets.ValueRange{
				//Range: "A2",
				MajorDimension: "COLUMNS",
				Values: shoeVals,
		}
		
		resp, err := srv.Spreadsheets.Values.Append(sheetID, range2, rb).ValueInputOption(valueInputOption).InsertDataOption(insertDataOption).Do()
		if err != nil {
			logger.Fatalf("uh oh %v %v", resp, err)
		}
		shoeVals = nil
	}
	
}

func resetFilter() *sheets.BatchUpdateSpreadsheetRequest{
	
	gridRange := &sheets.GridRange{
		StartColumnIndex: 0,
		EndColumnIndex: 10,
	}
	
	basicFilter := &sheets.BasicFilter{
		Range: gridRange,
	}
	
	setBasicFilterRequest := &sheets.SetBasicFilterRequest{
		Filter: basicFilter,
	}
	
	requests := []*sheets.Request{
		{SetBasicFilter: setBasicFilterRequest},
	}
	
	rb := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
        // TODO: Add desired fields of the request body.
	}
	
	return rb
}

func ShareLink(ss string) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
			logger.Fatalf("Unable to read client secret file: %v", err)
	}
	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
			logger.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := drive.New(client)
	if err != nil {
			logger.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	
	rb := &drive.Permission{
		Role: "reader",
		Type: "anyone",
	}
	_, err = srv.Permissions.Create(ss, rb).Do()
	if err != nil{
		logger.Fatal("Unable: %v", err)
	}
	
}