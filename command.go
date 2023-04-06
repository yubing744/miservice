package miservice

import (
    "strings"
)

var template = `Get Props: {prefix}<siid[-piid]>[,...]\n\
           {prefix}1,1-2,1-3,1-4,2-1,2-2,3\n\
Set Props: {prefix}<siid[-piid]=[#]value>[,...]\n\
           {prefix}2=#60,2-2=#false,3=test\n\
Do Action: {prefix}<siid[-piid]> <arg1|#NA> [...] \n\
           {prefix}2 #NA\n\
           {prefix}5 Hello\n\
           {prefix}5-4 Hello #1\n\n\
Call MIoT: {prefix}<cmd=prop/get|/prop/set|action> <params>\n\
           {prefix}action {quote}{{"did":"{did}","siid":5,"aiid":1,"in":["Hello"]}}{quote}\n\n\
Call MiIO: {prefix}/<uri> <data>\n\
           {prefix}/home/device_list {quote}{{"getVirtualModel":false,"getHuamiDevices":1}}{quote}\n\n\
Devs List: {prefix}list [name=full|name_keyword] [getVirtualModel=false|true] [getHuamiDevices=0|1]\n\
           {prefix}list Light true 0\n\n\
MIoT Spec: {prefix}spec [model_keyword|type_urn] [format=text|python|json]\n\
           {prefix}spec\n\
           {prefix}spec speaker\n\
           {prefix}spec xiaomi.wifispeaker.lx04\n\
           {prefix}spec urn:miot-spec-v2:device:speaker:0000A015:xiaomi-lx04:1\n\n\
MIoT Decode: {prefix}decode <ssecurity> <nonce> <data> [gzip]\n\
`

func IOCommandHelp(did string, prefix string) string {
    var quote string
    if prefix == "" {
        prefix = "?"
        quote = ""
    } else {
        quote = "'"
    }

    tmp := strings.ReplaceAll(template, "{prefix}", prefix)
    tmp = strings.ReplaceAll(tmp, "{quote}", quote)
    if did == "" {
        did = "267090026"
    }
    tmp = strings.ReplaceAll(tmp, "{did}", did)
    return tmp
}

func IOCommand(service *IOService, did string, text string, prefix string) (interface{}, error) {
    //cmd, arg := twinsSplit(text, " ", "")
    //if strings.HasPrefix(cmd, "/") {
    //    return service.Request(cmd, arg)
    //}
    //
    //if strings.HasPrefix(cmd, "prop") || cmd == "action" {
    //    var args map[string]interface{}
    //    if err := json.Unmarshal([]byte(arg), &args); err != nil {
    //        return nil, err
    //    }
    //    return service.Request(cmd, args)
    //}
    //
    //argv := strings.Split(arg, " ")
    //argc := len(argv)
    //var arg0 string
    //if argc > 0 {
    //    arg0 = argv[0]
    //}
    //var arg1 string
    //if argc > 1 {
    //    arg1 = argv[1]
    //}
    //var arg2 string
    //if argc > 2 {
    //    arg2 = argv[2]
    //}
    //switch cmd {
    //// Implement the cases for list, spec, and decode as methods for the IOService
    //case "list":
    //    a1 := false
    //    if arg1 != "" {
    //        a1, _ = strconv.ParseBool(arg1)
    //    }
    //    a2 := 0
    //    if arg2 != "" {
    //        a2, _ = strconv.Atoi(arg2)
    //    }
    //    return service.DeviceList(arg0, a1, a2) // Implement this method for the IOService
    //case "spec":
    //    return service.IotSpec(arg0, arg1) // Implement this method for the IOService
    //case "decode":
    //    if argc > 3 && argv[3] == "gzip" {
    //        return service.IotDecode(argv[0], argv[1], argv[2], true) // Implement this method for the IOService
    //    }
    //    return service.IotDecode(argv[0], argv[1], argv[2], false) // Implement this method for the IOService
    //}
    //if !strings.HasPrefix(did, "?") && !strings.HasPrefix(cmd, "？") && cmd != "help" && cmd != "-h" && cmd != "--help" {
    //    if !isDigit(did) {
    //        devices, err := service.DeviceList(did) // Implement this method for the IOService
    //        if err != nil {
    //            return nil, err
    //        }
    //        if len(devices) == 0 {
    //            return nil, errors.New("Device not found: " + did)
    //        }
    //        did = devices[0]["did"].(string)
    //    }
    //
    //    var props []interface{}
    //    setp := true
    //    miot := true
    //    for _, item := range strings.Split(cmd, ",") {
    //        key, value := twinsSplit(item, "=", "")
    //        siid, iid := twinsSplit(key, "-", "1")
    //        var prop any
    //        if strings.HasPrefix(siid, "#") && strings.HasPrefix(iid, "#") {
    //            prop = []int{int(siid[1]), int(iid[1])}
    //        } else {
    //            prop = []string{key}
    //            miot = false
    //        }
    //        if value == "" {
    //            setp = false
    //        } else if setp {
    //            prop = append(prop, stringOrValue(value))
    //        }
    //        props = append(props, prop)
    //    }
    //    if miot && argc > 0 {
    //        args := []interface{}{}
    //        if arg != "#NA" {
    //            for _, a := range argv {
    //                args = append(args, stringOrValue(a))
    //            }
    //        }
    //        return service.MiotAction(did, props[0], args) // Implement this method for the IOService
    //    }
    //
    //    var doProps func(string, []interface{}) (interface{}, error)
    //    if setp {
    //        if miot {
    //            doProps = service.MiotSetProps // Implement this method for the IOService
    //        } else {
    //            doProps = service.HomeSetProps // Implement this method for the IOService
    //        }
    //    } else {
    //        if miot {
    //            doProps = service.MiotGetProps // Implement this method for the IOService
    //        } else {
    //            doProps = service.HomeGetProps // Implement this method for the IOService
    //        }
    //    }
    //    return doProps(did, props)
    //}
    return IOCommandHelp(did, prefix), nil
}
