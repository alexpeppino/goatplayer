package echo // Single-instance package

import "fmt"

var (
	content       string
	beforeContent string
	afterContent  string
)

// Adds only the string
func AddSimple(s string) {
	content += s
}

// Adds "tput cup" and the string
func AddStringWithPos(row uint16, col uint16, strAdd string) {
	content += fmt.Sprintf("\x1b[%d;%dH%s", row+1, col+1, strAdd)
}

// tput cup 20 50
// \E[21;51H

/* Copies Window.contents to *s and deletes it */
func SlurpContents(s *string) {
	*s = content
	content = ""
}

func Flush() {
	// trace.Print("···· Term.Flush ····")
	fmt.Printf("%s%s%s", beforeContent, content, afterContent)
	// os.WriteFile("ciao.txt", []byte(w.contents), 0)
	content = ""
}

// // .
// func (this *vTerm) PrintTerm2() {
// 	fmt.Printf("%s%s%s", this.beforeContent, this.content, this.afterContent)
// 	// os.WriteFile("ciao.txt", []byte(w.contents), 0)
// 	this.content = ""
// }

// .
func PrintSimple() {
	fmt.Printf("%s", content)
	// os.WriteFile("ciao.txt", []byte(w.contents), 0)
	content = ""
}

func EchoTputCup(row uint16, col uint16, s string) {
	fmt.Printf("%s%s", TputCup(row, col), s)
}

func TputCup(row uint16, col uint16) string {
	return fmt.Sprintf("\x1b[%d;%dH", row, col)
}

