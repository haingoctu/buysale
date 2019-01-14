# buysale
A program that watches r/goodyearwelt for their bi-weekly buy/sell threads, and aggregates the seller's posting on to a Google Sheets.

This bot no longer runs due to a fallout between me and the r/goodyearwelt moderators.
Sample of what the Google Sheets looked like, https://docs.google.com/spreadsheets/d/1qggDIIWCtAjweY3n2G_SbhOaTGmNTan__y6ZyX2Vcnc/edit?usp=sharing

Dependencies:
  github.com/turnage/graw


TODO:
If the program crashes, and you rerun it, it does not pick up the latest buy/sell thread automactically. We want to read from the saved file, the current thread and Google sheet ID so we can continue from where we left of.
