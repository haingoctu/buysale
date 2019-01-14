# buysale
A program that watches r/goodyearwelt for their bi-weekly buy/sell threads, and aggregates the seller's posting on to a Google Sheets.

Dependencies:
  github.com/turnage/graw


TODO:
If the program crashes, and you rerun it, it does not pick up the latest buy/sell thread automactically. We want to read from the saved file, the current thread and Google sheet ID so we can continue from where we left of.
