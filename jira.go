package main

import (
  "encoding/json"
  "io/ioutil"
  "log"
  "net/http"
  "fmt"
  "bytes"
)

type LineRequest struct {
	Expand     string `json:"expand"`
	StartAt    int    `json:"startAt"`
	MaxResults int    `json:"maxResults"`
	Total      int    `json:"total"`
	Issues     []struct {
		Expand string `json:"expand"`
		ID     string `json:"id"`
		Self   string `json:"self"`
		Key    string `json:"key"`
		Fields struct {
			Assignee struct {
				Self       string `json:"self"`
				Name       string `json:"name"`
				Key        string `json:"key"`
				AccountID  string `json:"accountId"`
				AvatarUrls struct {
					Four8X48  string `json:"48x48"`
					Two4X24   string `json:"24x24"`
					One6X16   string `json:"16x16"`
					Three2X32 string `json:"32x32"`
				} `json:"avatarUrls"`
				DisplayName string `json:"displayName"`
				Active      bool   `json:"active"`
				TimeZone    string `json:"timeZone"`
				AccountType string `json:"accountType"`
			} `json:"assignee"`
		} `json:"fields"`
	} `json:"issues"`
}

var lineReq = LineRequest{} // 構造体つくっておく

type RequestBody struct {
    Channel string `json:"channel"`
    Name    string `json:"username"`
    Text    string `json:"text"`
}

func slack() {
  var namelist1 =jira("10100")
  var namelist2 =jira("10101")

  var data = RequestBody{
        Channel: "jira_test",
        Name:    "jiraリマインド",
        Text:    "メールでも送付されていると思いますが、直近期限のJIRAチケット状況を共有します\n ※期限を変更する際は、変更して問題ない理由を記載した上で担当マネージャに確認してもらってください。\n\n 期限切れチケット\n" + namelist1 + "\n https://daconsortium.atlassian.net/issues/?filter=10100 \n\n 本日期限チケット\n" + namelist2 + "\n https://daconsortium.atlassian.net/issues/?filter=10101 \n\n 未完了タスク一覧です\n https://daconsortium.atlassian.net/secure/Dashboard.jspa?selectPageId=10103 \n お忙しい中恐れ入りますが、こちらご対応のほどよろしくお願いいたします。",
    }

    dataJSON, err := json.Marshal(&data)
    if err != nil {
        fmt.Print(err)
    }

    req2, _ :=http.NewRequest("POST", "https://hooks.slack.com/services/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", bytes.NewBuffer(dataJSON))
    req2.Header.Set("Content-Type", "application/json")
    client2 := &http.Client{}
    resp2, err := client2.Do(req2)
    if err != nil {
        fmt.Print(err)
    }
    //fmt.Print(resp)
    defer resp2.Body.Close()
}

func jira(filter string)string{//期限切れ
  req, _ := http.NewRequest("GET", "https://daconsortium.atlassian.net/rest/api/2/search?maxResults=50&fields=assignee&jql=filter="+ filter,nil)
  req.Header.Set("Authorization","Basic c3VtaXJlLXVubm9AZGFjLmNvLmpwOnQ3UTQxMWxjR1NqUU5UVlBYRGplMUQ0Mw==")
  client := new(http.Client)
  resp, _ := client.Do(req)
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  //fmt.Printf("%d",resp.StatusCode) //ステータス確認
  if err != nil {
    log.Fatal(err)
  }

  err = json.Unmarshal(body, &lineReq)
  if err != nil {
    log.Fatal(err)
  }

  m := make(map[string]bool)
  for i, _ := range lineReq.Issues {
    m[lineReq.Issues[i].Fields.Assignee.DisplayName] = true
  }

  var namelist string
  for i, _ := range m{
    //fmt.Println(i)
    namelist = namelist + i + " "
  }
  return namelist
}

func main(){
  slack()
}
