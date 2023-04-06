package miservice

import (
    "os"
    "testing"
)

func TestLogin(t *testing.T) {
    user := os.Getenv("MI_USER")
    if user == "" {
        t.Skip("MI_USER not set")
    }
    pwd := os.Getenv("MI_PASS")
    acc := NewAccount(user, pwd, NewTokenStore("token.json"))
    err := acc.Login("xiaomiio")
    if err != nil {
        t.Error(err)
    }
}

func Test_secureUrl(t *testing.T) {
    type args struct {
        location  string
        nonce     int64
        ssecurity string
    }
    tests := []struct {
        name string
        args args
        want string
    }{
        {
            name: "1",
            args: args{
                location:  "https://sts.api.io.mi.com/sts?d=GBYUMAMUCF5RT0CF&ticket=0&pwd=1&p_ts=1680084122000&fid=0&p_lm=1&auth=ewsNVmAm42tnK0XyjoNJfFbbgvn2Ms3ocvxa7OmbfQ1OgrhNiVnykEZyPNilysnLBtb%2FBRR%2Fee3LoMF0lCVuxcdsnxxL%2BREEzw6PfcvOgFhuBKpbYGN%2Fn1qPMfbnuoZQ2J%2B70YLRgdCecq0U24yCZtpDZtqEsznsWWt8n6fNS2g%3D&m=1&_group=DEFAULT&tsl=0&p_ca=0&p_ur=CN&p_idc=China&nonce=zIwEuLSiq68Bq0R6&_ssign=EgvboyyZT6RXTY3y8uBmWPrUEmo%3D",
                nonce:     1108575485263259648,
                ssecurity: "TJNUSMjTUuDxv4oFKmTpDw==",
            },
            want: "https://sts.api.io.mi.com/sts?d=GBYUMAMUCF5RT0CF&ticket=0&pwd=1&p_ts=1680084122000&fid=0&p_lm=1&auth=ewsNVmAm42tnK0XyjoNJfFbbgvn2Ms3ocvxa7OmbfQ1OgrhNiVnykEZyPNilysnLBtb%2FBRR%2Fee3LoMF0lCVuxcdsnxxL%2BREEzw6PfcvOgFhuBKpbYGN%2Fn1qPMfbnuoZQ2J%2B70YLRgdCecq0U24yCZtpDZtqEsznsWWt8n6fNS2g%3D&m=1&_group=DEFAULT&tsl=0&p_ca=0&p_ur=CN&p_idc=China&nonce=zIwEuLSiq68Bq0R6&_ssign=EgvboyyZT6RXTY3y8uBmWPrUEmo%3D&clientSign=VjbGIM%2FyX0DRCnXxf0T%2BpioPcDM%3D"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := secureUrl(tt.args.location, tt.args.ssecurity, tt.args.nonce); got != tt.want {
                t.Errorf("secureUrl() = \n got  %v\n want %v", got, tt.want)
            }
        })
    }
}
