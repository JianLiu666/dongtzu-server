package githubSDK

import (
	"bytes"
	"dongtzu/config"
	"dongtzu/pkg/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GithubIssue struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Label string `json:"label"`
}

func CreateIssueForProvider(title string, profile model.Provider) error {
	// create issue to private repo
	apiURL := config.GetGlobalConfig().Github.RepoURL

	issueData := GithubIssue{
		Title: title,
		Body: fmt.Sprintf(`
## DB編號%s註冊申請 

1. 中文姓名(真實): %s
2. 申請Line官方帳號名稱: %s
3. 聯絡的個人Line ID: %s
4. 手機號碼: %s
5. Gamil: %s
6. 邀請碼: %s

以上是申請資料，麻煩客服協助處理喔～
		`, profile.ID, profile.RealName, profile.LineAtName, profile.LineID,
			profile.PhoneNum, profile.GmailAddr, profile.InviteCode),
		Label: "help wanted",
	}
	jsonData, _ := json.Marshal(issueData)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", apiURL, bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "token "+config.GetGlobalConfig().Github.APIToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[GithubSDK] http client execution err : \n %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("[GithubSDK] response http code from github %d\n", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("[GithubSDK] response body from github %d\n", string(body))
		return err
	}
	return nil
}
