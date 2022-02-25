package datadragon

type ChampionFull struct {
    Type string `json:"type"`
    Format string `json:"format"`
    Version string `json:"version"`
    Data map[string]Champion `json:"data"`
}

type Champion struct {
    Id string `json:"id"`
    Key string `json:"key"`
    Name string `json:"name"`
    Title string `json:"title"`
    Image Image `json:"image"`
    Blurb string `json:"blurb"`
    Tags []string `json:"tags"`
    Stats map[string]float64 `json:"stats"`
    Spells []Spell `json:"spells"`
    Passive Passive `json:"passive"`
}

type Passive struct {
    Name string `json:"name"`
    Description string `json:"description"`
    Image Image `json:"image"`
}

type Leveltip struct {
    Label []string `json:"label"`
    Effect []string `json:"effect"`
}

type Spell struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Description string `json:"description"`
    Cooldown []float64 `json:"cooldown"`
    Leveltip Leveltip `json:"leveltip"`
    Maxrank int32 `json:"maxrank"`
    Effect [][]float64 `json:"effect"`
    Range []float64 `json:"range"`
    Image Image `json:"image"`
}

type Image struct {
    Full string `json:"full"`
    Sprite string `json:"sprite"`
    Group string `json:"group"`
    X int `json:"x"`
    Y int `json:"y"`
    W int `json:"w"`
    H int `json:"h"`
}
