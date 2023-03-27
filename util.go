package miservice

import (
    "bytes"
    "compress/gzip"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/binary"
    "encoding/json"
    "fmt"
    "io"
    "math/rand"
    "strconv"
    "strings"
    "time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func getRandom(length int) string {
    charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    randomStr := make([]byte, length)
    for i := range randomStr {
        randomStr[i] = charset[r.Intn(len(charset))]
    }
    return string(randomStr)
}

func signNonce(ssecurity string, nonce string) (string, error) {
    decodedSsecurity, err := base64.StdEncoding.DecodeString(ssecurity)
    if err != nil {
        return "", err
    }

    decodedNonce, err := base64.StdEncoding.DecodeString(nonce)
    if err != nil {
        return "", err
    }

    hash := sha256.New()
    hash.Write(decodedSsecurity)
    hash.Write(decodedNonce)
    return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}

func signData(uri string, data map[string]interface{}, ssecurity string) map[string]interface{} {
    dataStr, err := json.Marshal(data)
    if err != nil {
        return nil
    }

    nonce := make([]byte, 12)
    _, err = rand.Read(nonce[:8])
    if err != nil {
        return nil
    }
    binary.BigEndian.PutUint32(nonce[8:], uint32(time.Now().Unix()/60))
    encodedNonce := base64.StdEncoding.EncodeToString(nonce)

    snonce, err := signNonce(ssecurity, encodedNonce)
    if err != nil {
        return nil
    }

    msg := fmt.Sprintf("%s&%s&%s&data=%s", uri, snonce, encodedNonce, dataStr)
    sign := hmac.New(sha256.New, []byte(snonce))
    sign.Write([]byte(msg))
    signature := base64.StdEncoding.EncodeToString(sign.Sum(nil))

    return map[string]interface{}{
        "_nonce":    encodedNonce,
        "data":      string(dataStr),
        "signature": signature,
    }
}

func twinsSplit(str, sep string, def string) (string, string) {
    pos := strings.Index(str, sep)
    if pos == -1 {
        return str, def
    }
    return str[0:pos], str[pos+1:]
}

func stringToValue(str string) interface{} {
    switch str {
    case "null", "none":
        return nil
    case "false":
        return false
    case "true":
        return true
    default:
        if intValue, err := strconv.Atoi(str); err == nil {
            return intValue
        }
        return str
    }
}

func stringOrValue(str string) interface{} {
    if str[0] == '#' {
        return stringToValue(str[1:])
    }
    return str
}

func unzip(data []byte) ([]byte, error) {
    reader, err := gzip.NewReader(bytes.NewReader(data))
    if err != nil {
        return nil, err
    }
    defer reader.Close()

    return io.ReadAll(reader)
}
