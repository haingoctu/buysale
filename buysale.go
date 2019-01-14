package main

import (
        "fmt"
		"flag"
		"os"
		"io/ioutil"
		"strings"
		
        "github.com/turnage/graw"
        "github.com/turnage/graw/reddit"
		"buysale/parse"
		"buysale/sheets"
		"buysale/mylog"
)

const VERSION = "20181114" //removed **

var (
 automod = flag.Bool("automod", false, "is the post created by an automod?")
 title = flag.String("title", "Buy/Sell/Trade", "posts are identified by this substring")
 subreddit = flag.String("subreddit", "hai_test", "our subreddit scope") 
 sheetID = flag.String("sheetid", "", "current sheet id")
 currPost = flag.String("currentpost", "", "current buy/sell posting in sub")
 permalink = flag.String("permalink", "", "permanent link of the current post")
 file = flag.String("file", "buysale_history", "file that logs the most recent posting")
 
 logger = mylog.GetInstance()
 
 item = map[int]*parse.Shoe{}
)

type currentPost struct {
	postDate	uint64
	postID 		string
	sheet		string
	permaLink	string
}

type buysaleBot struct {
	bot 		reddit.Bot
	cur 		currentPost
}

//Post identifies if the post is of Buy/Sell type. If it is,
// we check if it's newer than a previous one. In addition
// create a new sheet for a new Buy/Sell thread
func (r *buysaleBot) Post(p *reddit.Post) error {
		if parse.IdentifyPost(p.Title, *title, *automod) {
			err := findLatestThread(r, p.Permalink)
			if err != nil {
				logger.Print("find last thread broke")
			}
			logger.Print("current poste date ",r.cur.postDate," new post date: ", p.CreatedUTC)
			if r.cur.postDate < p.CreatedUTC {
				r.cur.sheet = sheets.Create(p.Title)
				sheets.ShareLink(r.cur.sheet)
				r.cur.postDate = p.CreatedUTC
				logger.Print("Setting new postDate: ",r.cur.postDate )
				r.cur.postID = p.ID
				logger.Print("Setting new postID: ",r.cur.postID )
			}
		}
	return nil
}

// Looks at all comments in our subreddit
func (r *buysaleBot) Comment(p *reddit.Comment) error{
	//check to see if this comment belongs to the latest Buy/Sell thread
	if strings.Contains(p.ParentID, r.cur.postID) {
		//retrieve the Buy/Sell thread's postDate if it doesn't exist
		if r.cur.postDate == 0 {
			pl, err := r.bot.Thread(*permalink)
			if err != nil {
				logger.Print("Failed to fetch permalink: ", err)
			}
			r.cur.postDate = pl.CreatedUTC
		}
		if r.cur.sheet == "" {
			logger.Print("creating new sheet")
			r.cur.sheet = sheets.Create(p.LinkTitle)
			sheets.ShareLink(r.cur.sheet)
		}
		
		item = parse.GetSaleItems(p)
		logger.Printf("There are %v in items\n", len(item))
		if len(item) > 0 {
			sheets.PostToSheets(item, r.cur.sheet)
		}
	}
	return nil
}


func main() {
		logger.Print("Running version: ", VERSION)
        if bot, err := reddit.NewBotFromAgentFile("buysale.agent", 0); err != nil {
                fmt.Println("Failed to create bot handle: ", err)
        } else {
				flag.Parse()
                cfg := graw.Config{Subreddits: []string{*subreddit}, SubredditComments: []string{*subreddit}}
				current := &currentPost{
					postID: *currPost,
					sheet:	*sheetID,
					permaLink: *permalink,
				}
				logger.Printf("current post, %s", current.postID)
				logger.Printf("main sheetid is, %s", current.sheet)
                handler := &buysaleBot{
					bot: bot,
					cur: *current,
				}
                if _, wait, err := graw.Run(handler, bot, cfg); err != nil {
                        fmt.Println("Failed to start graw run: ", err)
                } else {
                        fmt.Println("graw run failed: ", wait())
                }
        }
}

func findLatestThread(r *buysaleBot, newpost_permalink string) error {
	var file_data string
	linkToWrite := r.cur.permaLink
	if _, err := os.Stat(*file); !os.IsNotExist(err) {
		temp_data, err := ioutil.ReadFile(*file)
		if err != nil {
			logger.Print("Read file error: ", err)
		}
		file_data = string(temp_data)
	}
	
    new_pl, err := r.bot.Thread(newpost_permalink)
	old_pl_UTC := uint64(0)
	if err != nil {
		logger.Print("Failed to fetch permalink: ", err)
	}
	
	//find the date of the permalink in file if there's something there
	if file_data != "" {
		old_thread, err := r.bot.Thread(file_data)
		if err != nil {
			logger.Print("Failed to fetch permalink, file_data: ", err)
		}
		old_pl_UTC = old_thread.CreatedUTC
	}
	
	//find the date of current permalink stored
	if r.cur.permaLink != "" {
		old_thread, err := r.bot.Thread(r.cur.permaLink)
		if err != nil {
			logger.Print("Failed to fetch permalink, r.cur: ", err)
		}
		old_pl_UTC = old_thread.CreatedUTC
	}
	
	//compare old vs new dates
	if new_pl.CreatedUTC > old_pl_UTC {
		linkToWrite = new_pl.Permalink
	} else {
		linkToWrite = r.cur.permaLink
	}
	
	if linkToWrite != "" {
		f, err := os.Create(*file)
			if err != nil {
				logger.Print("Create file error: ", err)
			}
		defer f.Close()
		
		f.WriteString(linkToWrite)
	}
	
	return nil
}