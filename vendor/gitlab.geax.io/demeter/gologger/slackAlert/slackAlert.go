package slackAlert

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/rs/xid"
)

var startTmpl, shutdownTmpl, warnTmpl, errorTmpl *template.Template
var client *http.Client
var sendChan chan *SlackWarn
var flushChan chan chan struct{}
var guid, gitcommitnum, slackUrl, backendmodulesGitcommitnum, serverName, slackChannel, env string
var slackEnable bool

type SlackAlertsInfo struct {
	SlackChannel string `json:"slack_channel"`
	Env          string `json:"env"`
	SlackEnable  bool   `json:"slack_enable"`
	SlackUrl     string `json:"slack_url"` 
}

type SlackAlert struct {
	ServerName          string
	SlackChannel        string
	Env                 string
	Time                string
	Guid                string
	Message             string
	Gitcommitnum        string
	BackendGitcommitnum string
}

type SlackWarnAlert struct {
	ServerName          string
	SlackChannel        string
	Env                 string
	Time                string
	Guid                string
	Warns               []SlackWarn
	Gitcommitnum        string
	BackendGitcommitnum string
}

type SlackWarn struct {
	Timestamp string
	Warn      string
}

func localDial(network, addr string) (net.Conn, error) {
	dial := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 600 * time.Second,
	}

	conn, err := dial.Dial(network, addr)
	if err != nil && gitcommitnum != "" && slackUrl != "" {
		log.Printf("slackAlert/LocalDial error: %v", err)
		return conn, err
	}

	return conn, err
}

/*
第二個 string 參數之後依序是
1.serverName
2.xid
3.gitCommitNum
3.backendmodulesGitCommitnum
*/
func Run(sUrl string, alertsInfo *SlackAlertsInfo, sName, gitCommitNum, bmCommitnum string) {
	slackEnable = alertsInfo.SlackEnable
	slackChannel = alertsInfo.SlackChannel
	env = alertsInfo.Env
	if slackEnable && gitCommitNum != "" && sUrl != "" {
		slackUrl = sUrl
		serverName = sName
		guid = xid.New().String()
		gitcommitnum = gitCommitNum
		backendmodulesGitcommitnum = bmCommitnum
		client = &http.Client{
			Transport: &http.Transport{
				Dial:         localDial,
				MaxIdleConns: 50,
			},
		}
		startTmpl, _ = template.New("tmplStart").Parse(`
	{
        "channel": "{{.SlackChannel}}",
        "username": "{{.ServerName}} Start",
        "icon_emoji": ":large_blue_circle:",
        "attachments": [
        {
            "color": "#36a64f",
            "pretext": "{{.ServerName}} Start",
            "title": "{{.ServerName}} Start",
			"mrkdwn_in":["fields"],
            "fields": [
                {
                    "title": "環境",
                    "value": "{{.Env}}",
                    "short": true
                },{
                    "title": "Gitcommitnum",
                    "value": "{{.Gitcommitnum}}",
                    "short": true
                }{{if .BackendGitcommitnum}},{
                    "title": "BackendGitcommitnum",
                    "value": "{{.BackendGitcommitnum}}",
                    "short": true
                }{{end}},{
                    "title": "啟動辨識碼",
                    "value": "{{.Guid}}",
                    "short": false
                },{
                    "title": "Message",
                    "value": "{{.Message}}",
                    "short": false
                }
            ],
            "footer": "發生在",
            "ts": {{.Time}}
            }
        ]
    }`)

		shutdownTmpl, _ = template.New("tmplShutdown").Parse(`
	{
        "channel": "{{.SlackChannel}}",
        "username": "{{.ServerName}} Shutdown",
        "icon_emoji": ":red_circle:",
        "attachments": [
        {
            "color": "#4f4b48",
            "pretext": "{{.ServerName}} Shutdown",
            "title": "{{.ServerName}} Shutdown",
			"mrkdwn_in":["fields"],
            "fields": [
                {
                    "title": "環境",
                    "value": "{{.Env}}",
                    "short": true
                },{
                    "title": "Gitcommitnum",
                    "value": "{{.Gitcommitnum}}",
                    "short": true
                }{{if .BackendGitcommitnum}},{
                    "title": "BackendGitcommitnum",
                    "value": "{{.BackendGitcommitnum}}",
                    "short": true
                }{{end}},{
                    "title": "啟動辨識碼",
                    "value": "{{.Guid}}",
                    "short": false
                },{
                    "title": "Message",
                    "value": "{{.Message}}",
                    "short": false
                }
            ],
            "footer": "發生在",
            "ts": {{.Time}}
            }
        ]
	}`)

		errorTmpl, _ = template.New("tmplError").Parse(`
	{
        "channel": "{{.SlackChannel}}",
        "username": "{{.ServerName}} Error",
        "icon_emoji": ":red_circle:",
        "attachments": [
        {
            "color": "#f22307",
            "pretext": "{{.ServerName}} Error",
            "title": "{{.ServerName}} Error",
			"mrkdwn_in":["fields"],
            "fields": [
                {
                    "title": "環境",
                    "value": "{{.Env}}",
                    "short": true
                },{
                    "title": "Gitcommitnum",
                    "value": "{{.Gitcommitnum}}",
                    "short": true
                }{{if .BackendGitcommitnum}},{
                    "title": "BackendGitcommitnum",
                    "value": "{{.BackendGitcommitnum}}",
                    "short": true
                }{{end}},{
                    "title": "啟動辨識碼",
                    "value": "{{.Guid}}",
                    "short": false
                },{
                    "title": "Message",
                    "value": "{{.Message}}",
                    "short": false
                }
            ],
            "footer": "發生在",
            "ts": {{.Time}}
            }
        ]
    }`)

		warnTmpl, _ = template.New("tmplWarn").Parse(`
	{
        "channel": "{{.SlackChannel}}",
        "username": "{{.ServerName}} Warns",
        "icon_emoji": ":warning:",
        "attachments": [
        {
            "color": "#fbff1c",
            "pretext": "{{.ServerName}} Warns",
            "title": "{{.ServerName}} Warns",
			"mrkdwn_in":["fields"],
            "fields": [
                {
                    "title": "環境",
                    "value": "{{.Env}}",
                    "short": true
                },{
                    "title": "Gitcommitnum",
                    "value": "{{.Gitcommitnum}}",
                    "short": true
                }{{if .BackendGitcommitnum}},{
                    "title": "BackendGitcommitnum",
                    "value": "{{.BackendGitcommitnum}}",
                    "short": true
                }{{end}},{
                    "title": "啟動辨識碼",
                    "value": "{{.Guid}}",
                    "short": false
                }{{ range .Warns }},{
                    "title": "{{.Timestamp}}",
                    "value": "{{.Warn}}",
                    "short": false
                }{{end}}
            ],
            "footer": "發生在",
            "ts": {{.Time}}
            }
        ]
	}`)

		sendChan = make(chan *SlackWarn, 100)
		flushChan = make(chan chan struct{})
		listWarns := make([]*SlackWarn, 100)

		go func() {
			loopCheck := time.NewTicker(300 * time.Second)
			count := 0
			for {
				select {
				case v := <-sendChan:
					listWarns[count] = v
					count++
					if count >= 99 {
						errs := make([]SlackWarn, count)
						for i := 0; i < count; i++ {
							errs[i] = SlackWarn{
								Timestamp: listWarns[i].Timestamp,
								Warn:      listWarns[i].Warn,
							}
							listWarns[i] = nil
						}
						sendWarns(errs)
						count = 0
					}
				case <-loopCheck.C:
					if count != 0 {
						errs := make([]SlackWarn, count)
						for i := 0; i < count; i++ {
							errs[i] = SlackWarn{
								Timestamp: listWarns[i].Timestamp,
								Warn:      listWarns[i].Warn,
							}
							listWarns[i] = nil
						}
						sendWarns(errs)
						count = 0
					}
				case cb := <-flushChan:
					if count != 0 {
						errs := make([]SlackWarn, count)
						for i := 0; i < count; i++ {
							errs[i] = SlackWarn{
								Timestamp: listWarns[i].Timestamp,
								Warn:      listWarns[i].Warn,
							}
							listWarns[i] = nil
						}
						sendWarns(errs)
						count = 0
					}
					cb <- struct{}{}
				}
			}
		}()
	}
}

func SendStart(message string) {
	if slackEnable && gitcommitnum != "" && slackUrl != "" {
		var tpl bytes.Buffer
		sa := SlackAlert{BackendGitcommitnum: backendmodulesGitcommitnum, ServerName: serverName, Env: env, SlackChannel: slackChannel, Time: fmt.Sprint(time.Now().Unix()), Guid: guid, Message: message, Gitcommitnum: gitcommitnum}
		_ = startTmpl.Execute(&tpl, sa)
		req, err := http.NewRequest("POST", slackUrl, &tpl)
		if err != nil {
			log.Printf("slackAlert/SendStart error: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("slackAlert/SendStart error: %v", err)
			return
		}
		defer resp.Body.Close()
		// 拋棄 resp.Body 的資料才可以重複利用 Conns
		_, _ = io.Copy(ioutil.Discard, resp.Body)
	}
}

func SendShutdown(message string) {
	if slackEnable && gitcommitnum != "" && slackUrl != "" {
		var tpl bytes.Buffer
		sa := SlackAlert{Gitcommitnum: gitcommitnum, BackendGitcommitnum: backendmodulesGitcommitnum, ServerName: serverName, Env: env, SlackChannel: slackChannel, Time: fmt.Sprint(time.Now().Unix()), Guid: guid, Message: message}
		_ = shutdownTmpl.Execute(&tpl, sa)
		req, err := http.NewRequest("POST", slackUrl, &tpl)
		if err != nil {
			log.Printf("slackAlert/SendShutdown error: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("slackAlert/SendShutdown error: %v", err)
			return
		}
		defer resp.Body.Close()
		// 拋棄 resp.Body 的資料才可以重複利用 Conns
		_, _ = io.Copy(ioutil.Discard, resp.Body)
	}
}

func SendError(message string) {
	if slackEnable && gitcommitnum != "" && slackUrl != "" {
		FlushSync()
		var tpl bytes.Buffer
		sa := SlackAlert{Gitcommitnum: gitcommitnum, BackendGitcommitnum: backendmodulesGitcommitnum, ServerName: serverName, Env: env, SlackChannel: slackChannel, Time: fmt.Sprint(time.Now().Unix()), Guid: guid, Message: message}
		_ = errorTmpl.Execute(&tpl, sa)
		req, err := http.NewRequest("POST", slackUrl, &tpl)
		if err != nil {
			log.Printf("slackAlert/SendStart error: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("slackAlert/SendStart error: %v", err)
			return
		}
		defer resp.Body.Close()
		// 拋棄 resp.Body 的資料才可以重複利用 Conns
		_, _ = io.Copy(ioutil.Discard, resp.Body)
	}
}

func SendWarns(message string) {
	if slackEnable && gitcommitnum != "" && slackUrl != "" {
		go func(m string) {
			currentTime := time.Now()
			sendChan <- &SlackWarn{
				Warn:      m,
				Timestamp: currentTime.Format("2006-01-02 15:04:05"),
			}
		}(message)
	}
}

func sendWarns(warns []SlackWarn) {
	if slackEnable && gitcommitnum != "" && slackUrl != "" {
		var tpl bytes.Buffer
		sa := SlackWarnAlert{Gitcommitnum: gitcommitnum, BackendGitcommitnum: backendmodulesGitcommitnum, ServerName: serverName, Env: env, SlackChannel: slackChannel, Time: fmt.Sprint(time.Now().Unix()), Guid: guid, Warns: warns}
		_ = warnTmpl.Execute(&tpl, sa)
		req, err := http.NewRequest("POST", slackUrl, &tpl)
		if err != nil {
			log.Printf("slackAlert/SendErrors error: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("slackAlert/SendErrors error: %v", err)
			return
		}
		defer resp.Body.Close()
		// 拋棄 resp.Body 的資料才可以重複利用 Conns
		_, _ = io.Copy(ioutil.Discard, resp.Body)
	}
}

// 同步等待送出 slack
func FlushSync() {
	if slackEnable && gitcommitnum != "" && slackUrl != "" {
		cb := make(chan struct{})
		flushChan <- cb
		<-cb
	}
}

func SendErrorOnError(message string, e error) {
	if slackEnable && gitcommitnum != "" && slackUrl != "" {
		if e != nil {
			SendError(fmt.Sprint(message, e.Error()))
		}
	}
}

func IsSlackEnable() bool {
	if slackEnable && gitcommitnum != "" && slackUrl != "" {
		return true
	}
	return false
}
