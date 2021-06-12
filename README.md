# isucon-suburi-portal
ISUCONの練習のためのポータルサイト

# sample curl
- add score_log
```
curl -v localhost:8080/report -H "ReportToken:report_token_example" -H "Content-Type:application/json" -d '{"team_name": "won_the_first_prize", "score": 100000, "message": "best score!"}'
```
