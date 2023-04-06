package miservice

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type AIService struct {
    account *Account
}

func NewAIService(account *Account) *AIService {
    return &AIService{
        account: account,
    }
}

func (mnas *AIService) Request(uri string, data map[string]interface{}) (map[string]interface{}, error) {
    requestId := "app_ios_" + getRandom(30)
    if data != nil {
        data["requestId"] = requestId
    } else {
        uri += "&requestId=" + requestId
    }

    headers := http.Header{
        "User-Agent": []string{"MiHome/6.0.103 (com.xiaomi.mihome; build:6.0.103.1; iOS 14.4.0) Alamofire/6.0.103 MICO/iOSApp/appStore/6.0.103"},
    }

    var res map[string]interface{}
    mnas.account.Request("micoapi", "https://api2.mina.mi.com"+uri, nil /*data*/, nil, headers, true, &res)
    return res, nil
}

func (mnas *AIService) DeviceList(master int) ([]map[string]interface{}, error) {
    result, err := mnas.Request(fmt.Sprintf("/admin/v2/device_list?master=%d", master), nil)
    if err != nil {
        return nil, err
    }

    data, ok := result["data"].([]map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("failed to parse device list")
    }

    return data, nil
}

func (mnas *AIService) UbusRequest(deviceId, method, path string, message map[string]interface{}) (map[string]interface{}, error) {
    messageJSON, _ := json.Marshal(message)
    data := map[string]interface{}{
        "deviceId": deviceId,
        "message":  string(messageJSON),
        "method":   method,
        "path":     path,
    }

    return mnas.Request("/remote/ubus", data)
}

func (mnas *AIService) TextToSpeech(deviceId, text string) (map[string]interface{}, error) {
    return mnas.UbusRequest(deviceId, "text_to_speech", "mibrain", map[string]interface{}{"text": text})
}

func (mnas *AIService) PlayerSetVolume(deviceId string, volume int) (map[string]interface{}, error) {
    return mnas.UbusRequest(deviceId, "player_set_volume", "mediaplayer", map[string]interface{}{"volume": volume, "media": "app_ios"})
}

func (mnas *AIService) PlayerPause(deviceId string) (map[string]interface{}, error) {
    return mnas.UbusRequest(deviceId, "player_play_operation", "mediaplayer", map[string]interface{}{"action": "pause", "media": "app_ios"})
}

func (mnas *AIService) PlayerPlay(deviceId string) (map[string]interface{}, error) {
    return mnas.UbusRequest(deviceId, "player_play_operation", "mediaplayer", map[string]interface{}{"action": "play", "media": "app_ios"})
}

func (mnas *AIService) PlayerGetStatus(deviceId string) (
    map[string]interface{}, error) {
    return mnas.UbusRequest(deviceId, "player_get_play_status", "mediaplayer", map[string]interface{}{"media": "app_ios"})
}

func (mnas *AIService) PlayByUrl(deviceId, url string) (map[string]interface{}, error) {
    return mnas.UbusRequest(deviceId, "player_play_url", "mediaplayer", map[string]interface{}{"url": url, "type": 1, "media": "app_ios"})
}

func (mnas *AIService) SendMessage(devices []map[string]interface{}, devno int, message string, volume *int) (bool, error) {
    result := false
    for i, device := range devices {
        if devno == -1 || devno != i+1 || device["capabilities"].(map[string]interface{})["yunduantts"] != nil {
            deviceId := device["deviceID"].(string)
            if volume != nil {
                res, err := mnas.PlayerSetVolume(deviceId, *volume)
                result = err == nil && res != nil
            }
            if message != "" {
                res, err := mnas.TextToSpeech(deviceId, message)
                result = err == nil && res != nil
            }
            if devno != -1 || !result {
                break
            }
        }
    }

    return result, nil
}
