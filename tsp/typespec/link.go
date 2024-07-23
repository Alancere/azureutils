package typespec

import "fmt"

type Link struct {
	link            string
	additionalLinks []string
}

func (l *Link) Links() []string {
	links := make([]string, 0, 1+len(l.additionalLinks))
	links = append(links, l.link)
	links = append(links, l.additionalLinks...)

	return links
}

func (l *Link) String() string {
	output := fmt.Sprintf("directory: %s\n", l.link)
	for i, ad := range l.additionalLinks {
		if i == 0 {
			output += fmt.Sprintf("additional directory:\n")
		}
		output += fmt.Sprintf("  - %s\n", ad)
	}
	return output
}
