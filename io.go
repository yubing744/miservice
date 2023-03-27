package miservice

import (
    "fmt"
    "net/http"
    "strconv"
)

type MiIOService struct {
    account *MiAccount
    server  string
}

func NewMiIOService(account *MiAccount, region *string) *MiIOService {
    server := "https://"
    if region != nil && *region != "cn" {
        server += *region + "."
    }
    server += "api.io.mi.com/app"
    return &MiIOService{account: account, server: server}
}

func (m *MiIOService) MiIORequest(uri string, data map[string]interface{}) (map[string]interface{}, error) {
    prepareData := func(token map[string]string, cookies map[string]string) map[string]interface{} {
        cookies["PassportDeviceId"] = token["deviceId"]
        return signData(uri, data, token["xiaomiio"])
    }

    headers := http.Header{
        "User-Agent":                 []string{"iOS-14.4-6.0.103-iPhone12,3--D7744744F7AF32F0544445285880DD63E47D9BE9-8816080-84A3F44E137B71AE-iPhone"},
        "x-xiaomi-protocal-flag-cli": []string{"PROTOCAL-HTTP2"},
    }

    resp, err := m.account.MiRequest("xiaomiio", m.server+uri, prepareData, headers, true)
    if err != nil {
        return nil, err
    }

    result, ok := resp["result"].(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("error %s: %v", uri, resp)
    }

    return result, nil
}

func (s *MiIOService) HomeRequest(did, method string, params interface{}) (map[string]interface{}, error) {
    data := map[string]interface{}{
        "id":        1,
        "method":    method,
        "accessKey": "IOS00026747c5acafc2",
        "params":    params,
    }
    return s.MiIORequest("/home/rpc/"+did, data)
}

func (s *MiIOService) HomeGetProps(did string, props []string) (map[string]interface{}, error) {
    return s.HomeRequest(did, "get_prop", props)
}

func (s *MiIOService) HomeSetProps(did string, props []string) ([]interface{}, error) {
    results := make([]interface{}, len(props))
    for i, prop := range props {
        result, err := s.HomeSetProp(did, prop, []interface{}{prop})
        if err != nil {
            return nil, err
        }
        results[i] = result
    }
    return results, nil
}

func (s *MiIOService) HomeGetProp(did, prop string) (interface{}, error) {
    results, err := s.HomeGetProps(did, []string{prop})
    if err != nil {
        return nil, err
    }
    return results[prop], nil
}

func (s *MiIOService) HomeSetProp(did, prop string, value interface{}) (int, error) {
    result, err := s.HomeRequest(did, "set_"+prop, value)
    if err != nil {
        return 0, err
    }
    if result["result"] == "ok" {
        return 0, nil
    }
    return -1, nil
}

func (s *MiIOService) MiotRequest(cmd string, params interface{}) (map[string]interface{}, error) {
    return s.MiIORequest("/miotspec/"+cmd, map[string]interface{}{"params": params})
}
func (s *MiIOService) MiotGetProps(did string, iids [][]int) ([]interface{}, error) {
    params := make([]map[string]interface{}, len(iids))
    for i, iid := range iids {
        params[i] = map[string]interface{}{
            "did":  did,
            "siid": iid[0],
            "piid": iid[1],
        }
    }
    result, err := s.MiotRequest("prop/get", params)
    if err != nil {
        return nil, err
    }

    values := make([]interface{}, len(result))
    for i, it := range result {
        index, _ := strconv.Atoi(i)
        itm := it.(map[string]interface{})
        if code, ok := itm["code"].(int); ok && code == 0 {
            values[index] = itm["value"]
        } else {
            values[index] = nil
        }
    }
    return values, nil
}

func (s *MiIOService) MiotSetProps(did string, props [][]interface{}) ([]int, error) {
    params := make([]map[string]interface{}, len(props))
    for i, prop := range props {
        params[i] = map[string]interface{}{
            "did":   did,
            "siid":  prop[0],
            "piid":  prop[1],
            "value": prop[2],
        }
    }
    result, err := s.MiotRequest("prop/set", params)
    if err != nil {
        return nil, err
    }

    codes := make([]int, len(result))
    for i, it := range result {
        index, _ := strconv.Atoi(i)
        itm := it.(map[string]interface{})
        codes[index] = itm["code"].(int)
    }
    return codes, nil
}

func (s *MiIOService) MiotGetProp(did string, iid []int) (interface{}, error) {
    results, err := s.MiotGetProps(did, [][]int{iid})
    if err != nil {
        return nil, err
    }
    return results[0], nil
}

func (s *MiIOService) MiotSetProp(did string, iid []int, value interface{}) (int, error) {
    results, err := s.MiotSetProps(did, [][]interface{}{{iid[0], iid[1], value}})
    if err != nil {
        return 0, err
    }
    return results[0], nil
}

func (s *MiIOService) MiotAction(did string, iid []int, args []interface{}) (int, error) {
    result, err := s.MiotRequest("action", map[string]interface{}{
        "did":  did,
        "siid": iid[0],
        "aiid": iid[1],
        "in":   args,
    })
    if err != nil {
        return -1, err
    }
    return result["code"].(int), nil
}

type DeviceInfo struct {
    Name  string
    Model string
    Did   string
    Token string
}

func (s *MiIOService) DeviceList(name string, getVirtualModel bool, getHuamiDevices int) ([]DeviceInfo, error) {
    data := map[string]interface{}{
        "getVirtualModel": getVirtualModel,
        "getHuamiDevices": getHuamiDevices,
    }
    result, err := s.MiIORequest("/home/device_list", data)
    if err != nil {
        return nil, err
    }

    deviceList := result["list"].([]interface{})
    devices := make([]DeviceInfo, len(deviceList))
    for i, item := range deviceList {
        device := item.(map[string]interface{})
        devices[i] = DeviceInfo{
            Name:  device["name"].(string),
            Model: device["model"].(string),
            Did:   device["did"].(string),
            Token: device["token"].(string),
        }
    }
    return devices, nil
}
