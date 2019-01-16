# buysale
A program that watches r/goodyearwelt for their bi-weekly buy/sell threads, and aggregates the sellers' posting on to a Google Sheets.

The parsing is based upon the 'recommended' formatting of an item being sold.

This bot no longer runs due to a fallout between me and the r/goodyearwelt moderators.
Sample of what the Google Sheets looked like, https://docs.google.com/spreadsheets/d/1qggDIIWCtAjweY3n2G_SbhOaTGmNTan__y6ZyX2Vcnc/edit?usp=sharing

Dependencies:
  github.com/turnage/graw (you'll need to create your own config file)

#####TODO:
- If the program crashes, and you rerun it, it does not pick up the latest buy/sell thread automactically. We want to read from the saved file, the current thread and Google sheet ID so we can continue from where we left of.
- Refactor parse.go so that the parsing for many of the items is in one func instead of seperate cases.
