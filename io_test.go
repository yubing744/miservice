package miservice

import (
    "os"
    "testing"
)

func TestDeviceList(t *testing.T) {
    user := os.Getenv("MI_USER")
    if user == "" {
        t.Skip("MI_USER not set")
    }
    pwd := os.Getenv("MI_PASS")
    acc := NewAccount(user, pwd, NewTokenStore("token.json"))
    service := NewIOService(acc, nil)
    devices, _, err := service.DeviceList(false, 0)
    if err != nil {
        t.Error(err)
    }
    t.Log(devices)
}
