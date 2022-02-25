package datadragon

import (
    "io"
    "net/http"
    "encoding/json"
)

func latestVersion() string {
    resp, err := http.Get(VersionsURL)
    if err != nil {
        return currentVersion
    }

    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return currentVersion
    }

    versions := []string{}
    if err := json.Unmarshal(body, &versions); err != nil {
        return currentVersion
    }
    return versions[0]
}

func championFull() ChampionFull {
    resp, err := http.Get(ResolveChampionsURL("ko_KR"))
    if err != nil {
        panic(err)
    }

    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    championFull := ChampionFull{};
    if err := json.Unmarshal(body, &championFull); err != nil {
        panic(err)
    }
    return championFull
}

func LoadChampions() ChampionFull {
    SetVersion(latestVersion())
    return championFull()
}
