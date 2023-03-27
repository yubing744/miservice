package miservice

import (
    "encoding/json"
    "errors"
    "fmt"
    "strings"
)

func MiioCommandHelp(did string, prefix string) string {
    quote := ""
    if prefix != "?" {
        quote = "'"
    }
    return fmt.Sprintf(`Get Props: %s<siid[-piid]>[,...]
           %s1,1-2,1-3,1-4,2-1,2-2,2-3
Set Props: %s<siid[-piid]=[#]value>[,...]
%s2=#60,2-2=#false,3=test
Do Action: %s<siid[-piid]> <arg1|#NA> [...]
%s2 #NA
%s5 Hello
%s5-4 Hello #1

Call MIoT: %s<cmd=prop/get|/prop/set|action> <params>
%saction %s{"did":"%s or "267090026","siid":5,"aiid":1,"in":["Hello"]}%s

Call MiIO: %s/<uri> <data>
%s/home/device_list %s{"getVirtualModel":false,"getHuamiDevices":1}%s

Devs List: %slist [name=full|name_keyword] [getVirtualModel=false|true] [getHuamiDevices=0|1]
%slist Light true 0

MIoT Spec: %sspec [model_keyword|type_urn] [format=text|python|json]
%sspec
%sspec speaker
%sspec xiaomi.wifispeaker.lx04
%sspec urn:miot-spec-v2:device:speaker:0000A015:xiaomi-lx04:1

MIoT Decode: %sdecode <ssecurity> <nonce> <data> [gzip]
`, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, prefix, cmd, quote, did, quote, cmd, cmd, quote, quote, cmd, cmd, cmd, cmd, cmd)
}

func MiioCommand(service *MiIOService, did string, text string, prefix string) (interface{}, error) {
    cmd, arg := twinsSplit(text, " ", "")
    if strings.HasPrefix(cmd, "/") {
        return service.MiIORequest(cmd, map[string]string{}) // Implement this method for the MiIOService
    }

    if strings.HasPrefix(cmd, "prop") || cmd == "action" {
        var args map[string]interface{}
        if err := json.Unmarshal([]byte(arg), &args); err != nil {
            return nil, err
        }
        return service.MiIORequest(cmd, args) // Implement this method for the MiIOService
    }

    argv := strings.Split(arg, " ")
    argc := len(argv)
    switch cmd {
    // Implement the cases for list, spec, and decode as methods for the MiIOService
    case "list":
        return service.DeviceList(argc > 0 && argv[0], argc > 1 && argv[1], argc > 2 && argv[2]) // Implement this method for the MiIOService
    case "spec":
        return service.MiotSpec(argc > 0 && argv[0], argc > 1 && argv[1]) // Implement this method for the MiIOService
    case "decode":
        if argc > 3 && argv[3] == "gzip" {
            return service.MiotDecode(argv[0], argv[1], argv[2], true) // Implement this method for the MiIOService
        }
        return service.MiotDecode(argv[0], argv[1], argv[2], false) // Implement this method for the MiIOService
    }
    if !strings.HasPrefix(did, "?") && !strings.HasPrefix(cmd, "？") && cmd != "help" && cmd != "-h" && cmd != "--help" {
        if !did.isdigit() {
            devices, err := service.DeviceList(did) // Implement this method for the MiIOService
            if err != nil {
                return nil, err
            }
            if len(devices) == 0 {
                return nil, errors.New("Device not found: " + did)
            }
            did = devices[0]["did"].(string)
        }

        var props []interface{}
        setp := true
        miot := true
        for _, item := range strings.Split(cmd, ",") {
            key, value := twinsSplit(item, "=", "")
            siid, iid := twinsSplit(key, "-", "1")
            var prop any
            if strings.HasPrefix(siid, "#") && strings.HasPrefix(iid, "#") {
                prop = []int{int(siid[1]), int(iid[1])}
            } else {
                prop = []string{key}
                miot = false
            }
            if value == "" {
                setp = false
            } else if setp {
                prop = append(prop, stringOrValue(value))
            }
            props = append(props, prop)
        }
        if miot && argc > 0 {
            args := []interface{}{}
            if arg != "#NA" {
                for _, a := range argv {
                    args = append(args, stringOrValue(a))
                }
            }
            return service.MiotAction(did, props[0], args) // Implement this method for the MiIOService
        }

        var doProps func(string, []interface{}) (interface{}, error)
        if setp {
            if miot {
                doProps = service.MiotSetProps // Implement this method for the MiIOService
            } else {
                doProps = service.HomeSetProps // Implement this method for the MiIOService
            }
        } else {
            if miot {
                doProps = service.MiotGetProps // Implement this method for the MiIOService
            } else {
                doProps = service.HomeGetProps // Implement this method for the MiIOService
            }
        }
        return doProps(did, props)
    }
    return miioCommandHelp(did, prefix), nil
}
