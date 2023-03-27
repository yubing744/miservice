package miservice

import (
    "crypto/md5"
    "crypto/sha1"
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
)

type MiAccount struct {
    client     *http.Client
    username   string
    password   string
    tokenStore *MiTokenStore
    token      map[string]string
}

func NewMiAccount(client *http.Client, username string, password string, token_store *MiTokenStore) *MiAccount {
    return &MiAccount{
        client:     client,
        username:   username,
        password:   password,
        tokenStore: token_store,
    }
}

func (ma *MiAccount) Login(sid string) error {
    if ma.token == nil {
        ma.token = make(map[string]string)
        ma.token["deviceId"] = strings.ToUpper(getRandom(16))
        if ma.tokenStore != nil {
            loadedToken, err := ma.tokenStore.loadToken()
            if err == nil {
                ma.token = loadedToken
            }
        }
    }

    resp, err := ma.serviceLogin(fmt.Sprintf("serviceLogin?sid=%s&_json=true", sid))
    if err != nil {
        return err
    }

    if resp["code"] != "0" {
        data := url.Values{
            "_json":    {"true"},
            "qs":       {resp["qs"]},
            "sid":      {resp["sid"]},
            "_sign":    {resp["_sign"]},
            "callback": {resp["callback"]},
            "user":     {ma.username},
            "hash":     {strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(ma.password))))},
        }
        resp, err = ma.serviceLogin("serviceLoginAuth2", data)
        if err != nil {
            return err
        }
        if resp["code"] != "0" {
            return errors.New(fmt.Sprintf("Error: %v", resp))
        }
    }

    ma.token["userId"] = resp["userId"]
    ma.token["passToken"] = resp["passToken"]
    location := resp["location"]
    nonce := resp["nonce"]
    ssecurity := resp["ssecurity"]

    serviceToken, err := ma.securityTokenService(location, nonce, ssecurity)
    if err != nil {
        return err
    }
    ma.token[sid] = serviceToken

    if ma.tokenStore != nil {
        ma.tokenStore.saveToken(ma.token)
    }

    return nil
}

func (ma *MiAccount) serviceLogin(uri string, data ...url.Values) (map[string]string, error) {
    req, err := http.NewRequest(http.MethodGet, "https://account.xiaomi.com/pass/"+uri, nil)
    if err != nil {
        return nil, err
    }

    headers := http.Header{
        "User-Agent": []string{"APP/com.xiaomi.mihome APPV/6.0.103 iosPassportSDK/3.9.0 iOS/14.4 miHSTS"},
    }

    req.Header = headers
    req.URL.RawQuery = data[0].Encode()

    cookies := []*http.Cookie{
        {Name: "sdkVersion", Value: "3.9"},
        {Name: "deviceId", Value: ma.token["deviceId"]},
    }

    if passToken, ok := ma.token["passToken"]; ok {
        cookies = append(cookies, &http.Cookie{Name: "userId", Value: ma.token["userId"]})
        cookies = append(cookies, &http.Cookie{Name: "passToken", Value: passToken})
    }

    for _, cookie := range cookies {
        req.AddCookie(cookie)
    }

    resp, err := ma.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var jsonResponse map[string]string
    err = json.Unmarshal(body[11:], &jsonResponse)
    if err != nil {
        return nil, err
    }

    return jsonResponse, nil
}

func (ma *MiAccount) securityTokenService(location, nonce, ssecurity string) (string, error) {
    nsec := "nonce=" + nonce + "&" + ssecurity
    sum := sha1.Sum([]byte(nsec))
    clientSign := base64.StdEncoding.EncodeToString(sum[:])

    requestUrl := fmt.Sprintf("%s&clientSign=%s", location, url.QueryEscape(clientSign))
    resp, err := ma.client.Get(requestUrl)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    cookies := resp.Cookies()
    var serviceToken string

    for _, cookie := range cookies {
        if cookie.Name == "serviceToken" {
            serviceToken = cookie.Value
            break
        }
    }

    if serviceToken == "" {
        body, _ := io.ReadAll(resp.Body)
        return "", errors.New(string(body))
    }

    return serviceToken, nil
}

func (ma *MiAccount) MiRequest(sid, url string, data interface{}, headers http.Header, relogin bool) (map[string]interface{}, error) {
    if _, ok := ma.token[sid]; ok || ma.Login(sid) == nil {
        cookies := []*http.Cookie{
            {Name: "userId", Value: ma.token["userId"]},
            {Name: "serviceToken", Value: ma.token[sid]},
        }

        var req *http.Request
        var err error
        if data == nil {
            req, err = http.NewRequest(http.MethodGet, url, nil)
        } else {
            jsonData, _ := json.Marshal(data)
            req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(string(jsonData)))
        }
        if err != nil {
            return nil, err
        }

        req.Header = headers
        for _, cookie := range cookies {
            req.AddCookie(cookie)
        }

        resp, err := ma.client.Do(req)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()

        if resp.StatusCode == 200 {
            var result map[string]interface{}
            err = json.NewDecoder(resp.Body).Decode(&result)
            if err != nil {
                return nil, err
            }

            if code, ok := result["code"].(float64); ok && code == 0 {
                return result, nil
            }

            if message, ok := result["message"].(string); ok && strings.Contains(strings.ToLower(message), "auth") {
                resp.StatusCode = 401
            }
        }

        if resp.StatusCode == 401 && relogin {
            ma.token = nil
            if ma.tokenStore != nil {
                ma.tokenStore.saveToken(nil)
            }
            return ma.MiRequest(sid, url, data, headers, false)
        }

        body, _ := io.ReadAll(resp.Body)
        return nil, errors.New(fmt.Sprintf("Error %s: %s", url, string(body)))
    }

    return nil, errors.New("login failed")
}
