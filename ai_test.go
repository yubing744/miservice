package miservice

import (
    "os"
    "testing"
)

func TestAIDeviceList(t *testing.T) {
    user := os.Getenv("MI_USER")
    if user == "" {
        t.Skip("MI_USER not set")
    }
    pwd := os.Getenv("MI_PASS")
    acc := NewAccount(user, pwd, NewTokenStore("token.json"))
    service := NewAIService(acc)
    devices, err := service.DeviceList(0)
    if err != nil {
        t.Error(err)
    }
    t.Log(devices)
    //status, err := service.PlayerGetStatus("")
    //if err != nil {
    //    t.Error(err)
    //}
    //t.Log(status)
}
