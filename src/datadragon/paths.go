package datadragon

const BaseURL = "http://ddragon.leagueoflegends.com/cdn"
const VersionsURL = "https://ddragon.leagueoflegends.com/api/versions.json"
var currentVersion = "12.4.1"

func SetVersion(version string) {
     currentVersion = version
}

func ResolveBaseURL() string {
    return BaseURL + "/" + currentVersion
}

func ResolveChampionsURL(lang string) string {
    return ResolveBaseURL() + "/data/" + lang + "/championFull.json"
}

func ResolveImageURL(image Image) string {
    return ResolveBaseURL() + "/img/" + image.Group + "/" + image.Full
}

func ResolveSpriteURL(image Image) string {
    return ResolveBaseURL() + "/img/sprite/" + image.Sprite
}

func ResolveChampionTileImageURL(key string) string {
    return ResolveBaseURL() + "/img/champion/tiles/" + key + "_0.png"
}
