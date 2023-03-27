package miservice

import (
    "encoding/json"
    "os"
)

type MiTokenStore struct {
    tokenPath string
}

func NewMiTokenStore(tokenPath string) *MiTokenStore {
    return &MiTokenStore{tokenPath: tokenPath}
}

func (mts *MiTokenStore) loadToken() (map[string]string, error) {
    var token map[string]string
    if _, err := os.Stat(mts.tokenPath); os.IsNotExist(err) {
        return nil, err
    }
    data, err := os.ReadFile(mts.tokenPath)
    if err != nil {
        return nil, err
    }
    err = json.Unmarshal(data, &token)
    return token, err
}

func (mts *MiTokenStore) saveToken(token map[string]string) error {
    var err error
    if token != nil {
        data, err := json.MarshalIndent(token, "", "  ")
        if err != nil {
            return err
        }
        err = os.WriteFile(mts.tokenPath, data, 0644)
        if err != nil {
            return err
        }
    } else {
        err = os.Remove(mts.tokenPath)
        if os.IsNotExist(err) {
            err = nil
        }
    }
    return err
}
