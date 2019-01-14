package parse

import ("strings"
		"os"
		"bufio"
		"fmt"
		"reflect"
		
		"github.com/grokify/html-strip-tags-go"
		"github.com/turnage/graw/reddit"
		"github.com/mvdan/xurls"
		"buysale/mylog"
)

var (
		infohandler = os.Stdout
		logger = mylog.GetInstance()
	)
	
type Shoe struct {
	User			string
	Make			string
	Leather			string
	Sole			string
	ImageLink		string
	Notes			string
	PermLink		string
	Condition		string
	Price			string
	Size			string
}

//identify the reoccuring post and optional arg of automod authored
func IdentifyPost(postTitle, title string, automod bool) bool {
	logger.Print("Post title: ", postTitle, "; FlagTitle: ", title)
	if automod && strings.Contains(postTitle, title) {return true}
	if !automod {
		logger.Print("automod flag is: ", automod)
		if strings.Contains(postTitle, title) {return true}
	}
	logger.Print("This post is not a Buy/Sell Thread")
	return false
}

func GetSaleItems(p *reddit.Comment) map[int]*Shoe {
	scanner := bufio.NewScanner(strings.NewReader(p.Body))
	html_scanner := bufio.NewScanner(strings.NewReader(p.BodyHTML))
	i := 0
	ss := make(map[int]*Shoe)
	ss[i] = &Shoe{}
	prevShoe := false
	logger.Print("User is ", p.Author)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), "maker/model:") ||
			strings.Contains(strings.ToLower(line), "maker:") || 
			strings.Contains(strings.ToLower(line), "model:") ||
			strings.Contains(strings.ToLower(line), "maker/model**:") ||
			strings.Contains(strings.ToLower(line), "maker**:") || 
			strings.Contains(strings.ToLower(line), "model**:") {
			line = strip.StripTags(line)
			loc := strings.Index(strings.ToLower(line), ":")
			line = line[loc+1:]
			line = remove(line, "*")
			line = strings.TrimLeft(line, " ")
			line = strings.TrimRight(line, " ")
			if ss[i].Make != "" {
				i = i + 1
				ss[i] = &Shoe{}
			}
			ss[i].Make = line
			prevShoe = true
			ss[i].User = p.Author
			logger.Print("Found maker/model: ", line)
		} else if prevShoe && strings.Contains(strings.ToLower(line), "size:") ||
			strings.Contains(strings.ToLower(line), "size**:") ||
			strings.Contains(strings.ToLower(line), "size:**") {
			line = strip.StripTags(line)
			loc := strings.Index(strings.ToLower(line), ":")
			line = line[loc+1:]
			line = remove(line, "*")
			line = strings.TrimLeft(line, " ")
			line = strings.TrimRight(line, " ")
			logger.Print("Found size: ", line)
			ss[i].Size = line
		} else if prevShoe && strings.Contains(strings.ToLower(line), "leather:") ||
			strings.Contains(strings.ToLower(line), "leather**:") ||
			strings.Contains(strings.ToLower(line), "leather:**") {
			line = strip.StripTags(line)
			loc := strings.Index(strings.ToLower(line), ":")
			line = line[loc+1:]
			line = remove(line, "*")
			line = strings.TrimLeft(line, " ")
			line = strings.TrimRight(line, " ")
			ss[i].Leather = line
			logger.Print("Found leather: ", line)
		} else if prevShoe && strings.Contains(strings.ToLower(line), "sole:") ||
			strings.Contains(strings.ToLower(line), "sole**:") ||
			strings.Contains(strings.ToLower(line), "sole:**") {
			line = strip.StripTags(line)
			loc := strings.Index(strings.ToLower(line), ":")
			line = line[loc+1:]
			line = remove(line, "*")
			line = strings.TrimLeft(line, " ")
			line = strings.TrimRight(line, " ")
			ss[i].Sole = line
			logger.Print("Found sole: ", line)
		} else if prevShoe && strings.Contains(strings.ToLower(line), "price:") ||
			strings.Contains(strings.ToLower(line), "price**:") ||
			strings.Contains(strings.ToLower(line), "price:**") {
			if ss[i].Price != "" { continue }
			line = strip.StripTags(line)
			loc := strings.Index(strings.ToLower(line), ":")
			line = line[loc+1:]
			line = remove(line, "*")
			line = strings.TrimLeft(line, " ")
			line = strings.TrimRight(line, " ")
			ss[i].Price = line
			logger.Print("Found price: ", line)
		} else if prevShoe && strings.Contains(strings.ToLower(line), "wears/condition:") ||
			strings.Contains(strings.ToLower(line), "wears/condition**:") ||
			strings.Contains(strings.ToLower(line), "wears/condition:**") {
			line = strip.StripTags(line)
			loc := strings.Index(strings.ToLower(line), ":")
			line = line[loc+1:]
			line = remove(line, "*")
			line = strings.TrimLeft(line, " ")
			line = strings.TrimRight(line, " ")
			ss[i].Condition = line
			logger.Print("Found wears/condition: ", line)
		} else if prevShoe && strings.Contains(strings.ToLower(line), "notes:") ||
			strings.Contains(strings.ToLower(line), "notes**:") ||
			strings.Contains(strings.ToLower(line), "notes:**") {
			line = strip.StripTags(line)
			loc := strings.Index(strings.ToLower(line), ":")
			line = line[loc+1:]
			line = remove(line, "*")
			line = strings.TrimLeft(line, " ")
			line = strings.TrimRight(line, " ")
			ss[i].Notes = line
			logger.Print("Found notes: ", line)
		} else if prevShoe && strings.Contains(strings.ToLower(line), "images:") ||
			strings.Contains(strings.ToLower(line), "images**:") ||
			strings.Contains(strings.ToLower(line), "images:**") {
			html_line := ""
			for html_scanner.Scan() {
				if l := html_scanner.Text();
				strings.Contains(strings.ToLower(l), "images") {
					html_line = html_scanner.Text()
					logger.Print("Found html_line: ", html_line)
					break
				}
			}
			html_line = xurls.Relaxed().FindString(html_line)
			ss[i].ImageLink = html_line
			logger.Print("Found images: ", html_line)
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Print(err)
	}
	if len(ss) == 1 && ss[0].Make == "" {
		return nil
	}
	return ss
}

func remove(s string, sub string) string{
	i := strings.Index(s, "*")
	for i > -1 {
		z := fmt.Sprint(s[:i])
		y := fmt.Sprint(s[i+1:])
		s = z+y
		fmt.Println(s)
		i = strings.Index(s, sub)
	}
	return s
}

func (s *Shoe) Reflectz(st string) string{
	if v := reflect.ValueOf(s).Elem().FieldByName(st); v.IsValid() {
		return fmt.Sprintf("%s", v)
	}
	return ""
}