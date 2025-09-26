package init

type Template struct {
    Name  string `json:"name"`
    Items []Item `json:"items"`
}

type Item struct {
    Name    string `json:"name"`
    Message string `json:"message"`
    Preset  string `json:"preset"`
    Path    string `json:"path"`
    Link    string `json:"link"`
    Extras  *Extras `json:"extras,omitempty"`
}

// Extras represents optional additional information sent by the templates API
// Currently supports type "env" with a list of inputs to be collected from the user
type Extras struct {
    Type   string       `json:"type"`
    Inputs []ExtraInput `json:"inputs"`
}

// ExtraInput represents a single input for the extras block
// For type "env", each input contains a key and a user-facing prompt text
type ExtraInput struct {
    Key  string `json:"key"`
    Text string `json:"text"`
}
