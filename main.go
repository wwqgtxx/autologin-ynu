package main

import (
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

var Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

func getRad(client *http.Client) (map[string]interface{}, error) {
	url1, err := url.Parse("http://202.203.208.5/cgi-bin/rad_user_info")
	if err != nil {
		return nil, err
	}
	jsonp := NewJsonp()
	query := url1.Query()
	query.Set("callback", jsonp.CallbackString)
	query.Set("_", timestampString())
	url1.RawQuery = query.Encode()
	req, err := http.NewRequest("GET", url1.String(), strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	setCommonHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	b = jsonp.RemoveJsonP(b)

	var jsonData map[string]interface{}
	err = json.Unmarshal(b, &jsonData)

	return jsonData, err
}

func getUserIPAndAcID(client *http.Client) (string, string, error) {
	req, err := http.NewRequest("GET", "http://202.203.208.5/", strings.NewReader(""))
	if err != nil {
		return "", "", err
	}
	setCommonHeaders(req)
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	acID := resp.Request.URL.Query().Get("ac_id")
	defer resp.Body.Close()
	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", err
	}
	//document.Find("script").EachWithBreak(func(i int, selection *goquery.Selection) bool {
	//	text := selection.Text()
	//	if len(text) > 0 {
	//		fmt.Println(text)
	//	}
	//	return true
	//})
	if len(acID) == 0 {
		acID, _ = document.Find("#ac_id").First().Attr("value")
	}
	userIP, _ := document.Find("#user_ip").First().Attr("value")
	return userIP, acID, err
}

func getChallenge(client *http.Client, userName string, userIP string) (string, error) {
	url1, err := url.Parse("http://202.203.208.5/cgi-bin/get_challenge")
	if err != nil {
		return "", err
	}
	jsonp := NewJsonp()
	query := url1.Query()
	query.Set("callback", jsonp.CallbackString)
	query.Set("username", userName)
	query.Set("ip", userIP)
	query.Set("_", timestampString())
	url1.RawQuery = query.Encode()
	req, err := http.NewRequest("GET", url1.String(), strings.NewReader(""))
	if err != nil {
		return "", err
	}
	setCommonHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	b = jsonp.RemoveJsonP(b)

	var jsonData map[string]interface{}
	if err = json.Unmarshal(b, &jsonData); err == nil {
		if value, ok := jsonData["challenge"]; ok {
			if result, ok := value.(string); ok {
				return result, nil
			} else {
				return "", errors.New("fail to get challenge")
			}
		} else {
			return "", errors.New("fail to get challenge")
		}
	}
	return "", err
}

func getPortalLogin(client *http.Client, userName string, passWord string, acID string, userIP string, challenge string) (string, error) {
	url1, err := url.Parse("http://202.203.208.5/cgi-bin/srun_portal")
	if err != nil {
		return "", err
	}
	jsonp := NewJsonp()
	query := url1.Query()
	query.Set("callback", jsonp.CallbackString)
	query.Set("action", "login")
	query.Set("username", userName)
	hmd5 := pwd(passWord, challenge)
	query.Set("password", "{MD5}"+hmd5)
	query.Set("ac_id", acID)
	query.Set("ip", userIP)

	i := info(userName, passWord, acID, userIP, challenge)

	chkstr := challenge + userName
	chkstr += challenge + hmd5
	chkstr += challenge + acID
	chkstr += challenge + userIP
	chkstr += challenge + "200"
	chkstr += challenge + "1"
	chkstr += challenge + i

	query.Set("info", i)
	query.Set("chksum", chksum(chkstr))

	query.Set("n", "200")
	query.Set("type", "1")
	query.Set("os", "Windows 10")
	query.Set("name", "Windows")
	query.Set("_", timestampString())
	url1.RawQuery = query.Encode()
	Logger.Println(url1.String())
	req, err := http.NewRequest("GET", url1.String(), strings.NewReader(""))
	if err != nil {
		return "", err
	}
	setCommonHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	b = jsonp.RemoveJsonP(b)

	Logger.Println(string(b))

	var jsonData map[string]interface{}
	if err = json.Unmarshal(b, &jsonData); err == nil {
		if value, ok := jsonData["error"]; ok {
			if result, ok := value.(string); ok {
				return result, nil
			} else {
				return "", errors.New("fail to get error")
			}
		} else {
			return "", errors.New("fail to get error")
		}
	}
	return "", err
}

func getPortalLogout(client *http.Client) (string, error) {
	url1, err := url.Parse("http://202.203.208.5/cgi-bin/srun_portal")
	if err != nil {
		return "", err
	}
	jsonp := NewJsonp()
	query := url1.Query()
	query.Set("callback", jsonp.CallbackString)
	query.Set("action", "logout")

	url1.RawQuery = query.Encode()
	Logger.Println(url1.String())
	req, err := http.NewRequest("GET", url1.String(), strings.NewReader(""))
	if err != nil {
		return "", err
	}
	setCommonHeaders(req)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	b = jsonp.RemoveJsonP(b)

	Logger.Println(string(b))

	var jsonData map[string]interface{}
	if err = json.Unmarshal(b, &jsonData); err == nil {
		if value, ok := jsonData["error"]; ok {
			if result, ok := value.(string); ok {
				return result, nil
			} else {
				return "", errors.New("fail to get error")
			}
		} else {
			return "", errors.New("fail to get error")
		}
	}
	return "", err
}

func doLogout(client *http.Client) {
	portal, err := getPortalLogout(client)
	if err != nil {
		Logger.Println(err)
		return
	}

	if portal == "ok" {
		Logger.Println("Finish logout!")
		return
	}
}

func doLogin(client *http.Client, config *Config) {
	Logger.Println("Detecting...")
	jar, err := cookiejar.New(nil)
	if err != nil {
		Logger.Println(err)
		return
	}
	client.Jar = jar

	rad, err := getRad(client)
	if err != nil {
		Logger.Println(err)
		return
	}
	if value, ok := rad["error"]; ok {
		if result, ok := value.(string); ok {
			if result == "ok" {
				Logger.Println("Was login!")
				return
			}
			Logger.Println("Detect logout!")
			Logger.Println("Try to login!")

			userIP, acID, err := getUserIPAndAcID(client)
			if err != nil {
				Logger.Println(err)
				return
			}
			if len(acID) == 0 {
				acID = "0"
			}
			if len(userIP) == 0 {
				clientIP, ok := rad["client_ip"]
				if !ok {
					Logger.Println("get rad[\"client_ip\"] failed")
					return
				}
				userIP, ok = clientIP.(string)
				if !ok {
					Logger.Println("client_ip is not string")
					return
				}
			}

			challenge, err := getChallenge(client, config.UserName, userIP)
			if err != nil {
				Logger.Println(err)
				return
			}

			portal, err := getPortalLogin(client, config.UserName, config.PassWord, acID, userIP, challenge)
			if err != nil {
				Logger.Println(err)
				return
			}

			if portal == "ok" {
				Logger.Println("Finish login!")
				return
			}

		}
	}
}

func main() {
	config := NewConfig()
	err := config.Load("config.tmp.json")
	if err != nil && os.IsNotExist(err) {
		if os.IsNotExist(err) {
			err = config.Load("config.json")
		}
	}
	if err != nil {
		panic(err)
	}

	Logger.Println("Init Username:" + config.UserName)
	Logger.Println("Init Password:" + config.PassWord)
	jar, err := cookiejar.New(nil)
	if err != nil {
		Logger.Println(err)
	}
	client := &http.Client{Jar: jar}
	if len(os.Args) > 1 && strings.ToLower(os.Args[1]) == "logout" {
		doLogout(client)
		return
	}
	for {
		doLogin(client, config)
		time.Sleep(1 * time.Second)
	}
}
