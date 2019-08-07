package model

type People struct {
	Name     string   `json:"name"`
	FilmsURL []string `json:"films"`
}

func (p People) String() string {
	return p.Name
}

type PeopleContainer struct {
	Count   int      `json:"count"`
	Results []People `json:"results"`
}

type Film struct {
	Title string `json:"title"`
}

func (f Film) String() string {
	return f.Title
}
