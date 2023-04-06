package miservice

import (
    "net/url"
    "reflect"
    "testing"
)

func Test_signNonce(t *testing.T) {
    type args struct {
        ssecurity string
        nonce     string
    }
    tests := []struct {
        name string
        args args
        want string
    }{
        {name: "1", args: args{
            ssecurity: "TJNUSMjTUuDxv4oFKmTpDw==",
            nonce:     "5BefBVR/SbABq0Wz",
        }, want: "TCOrc8R5WAUQ0UUfuQt07Ou3IM8VkCa1rLSo3ZeceQM="},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := signNonce(tt.args.ssecurity, tt.args.nonce)
            if err != nil {
                t.Error(err)
            }
            if got != tt.want {
                t.Errorf("signNonce() got = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test_signData(t *testing.T) {
    type args struct {
        uri       string
        data      any
        ssecurity string
    }
    tests := []struct {
        name string
        args args
        want url.Values
    }{
        {name: "1",
            args: args{
                uri:       "/home/device_list",
                data:      `{"getVirtualModel": false, "getHuamiDevices": 0}`,
                ssecurity: "TJNUSMjTUuDxv4oFKmTpDw==",
            },
            want: url.Values{
                "_nonce":    {"5BefBVR/SbABq0Wz"},
                "data":      {`{"getVirtualModel": false, "getHuamiDevices": 0}`},
                "signature": {"lSeqFg9S0JRsAt4SKjmpLF2vzzfdmvfYtMB9nvEfk7o="},
            },
        }}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            genNonce = func() string {
                return "5BefBVR/SbABq0Wz"
            }
            if got := signData(tt.args.uri, tt.args.data, tt.args.ssecurity); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("signData() = %v, want %v", got, tt.want)
            }
        })
    }
}
