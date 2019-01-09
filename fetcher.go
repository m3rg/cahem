package main

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/text/encoding/htmlindex"
)

type Course struct {
	Name        string
	Description []string
	Index       int
}

const (
	dataFile = "cahem.dat"
)

func Fetch() error {
	resp, err := http.Get(config.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	htm, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//@TODO: Auto detect page encoding
	e, err := htmlindex.Get("iso-8859-9")
	if err != nil {
		return err
	}
	body := e.NewDecoder().Reader(strings.NewReader(string(htm)))
	doc, _ := html.Parse(body)
	contentNode, err := GetContentContainer(doc)
	if err != nil {
		return err
	}
	content := RenderNode(contentNode)
	updateDate := ParseUpdateDate(&content)
	courses := ParseCourseList(&content)

	mailContent := Diff(&updateDate, courses)
	if mailContent != nil {
		log.Println("Sending e-mail.")
		err = SendMail(mailContent)
		if err != nil {
			log.Printf("smtp error: %s", err)
		}
	}
	WriteCourses(&updateDate, courses)
	return nil
}

func GetContentContainer(doc *html.Node) (*html.Node, error) {
	var b *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			var attributes []html.Attribute = n.Attr
			for _, attr := range attributes {
				if attr.Key == "id" && attr.Val == "araSayfaNrmlicerik" {
					b = n
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if b != nil {
		return b, nil
	}
	return nil, nil
}

func ParseUpdateDate(content *string) string {
	re := regexp.MustCompile("\\(Güncelleme Tarihi:\\s*(.*?)\\s*\\)")
	if re.MatchString(*content) {
		groups := re.FindStringSubmatch(*content)
		return groups[1]
	}
	return ""
}

func ParseCourseList(content *string) *[]Course {
	re := regexp.MustCompile("<p class=\"MsoNormal\" style=\"margin-bottom: \\.0001pt; line-height: normal;\"><strong.*?>(.*)</strong></p>")
	reLine := regexp.MustCompile("^([0-9]+)\\s*-\\s*(.*)$")
	groups := re.FindAllStringSubmatch(*content, -1)
	var course Course
	var description []string
	var courses []Course
	for _, group := range groups {
		if group[1] == " " {
			continue
		}
		lineGroups := reLine.FindStringSubmatch(group[1])
		if len(lineGroups) > 0 {
			index, err := strconv.Atoi(lineGroups[1])
			if err != nil {
				continue
			}
			if len(description) > 0 {
				course.Description = description
				courses = append(courses, course)
			}
			course = Course{Index: index, Name: lineGroups[2]}
			description = []string{}
		} else {
			description = append(description, group[1])
		}
	}
	course.Description = description
	courses = append(courses, course)
	return &courses
}

func WriteCourses(updateDate *string, courses *[]Course) {
	var content string = *updateDate + "\r\n"
	for _, course := range *courses {
		content += "###" + strconv.Itoa(course.Index) + "- " + course.Name + "\r\n"
		content += strings.Join(course.Description, "\r\n") + "\r\n"
	}
	ioutil.WriteFile(dataFile, []byte(content), 0644)
}

func Diff(updateDate *string, courses *[]Course) *string {
	file, err := os.Open(dataFile)
	if err != nil {
		log.Println(err)
		return GetMailContent(updateDate, courses, nil)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	courseMap := map[string]string{}
	re := regexp.MustCompile("^###([0-9]+)\\s*-\\s*(.*)$")
	i := -1
	for scanner.Scan() {
		i++
		line := scanner.Text()
		if i == 0 {
			if line == *updateDate { //no change
				return nil
			}
			continue
		}
		if strings.HasPrefix(line, "###") {
			groups := re.FindStringSubmatch(line)
			courseMap[groups[2]] = strings.Replace(line, "###", "", -1)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return GetMailContent(updateDate, courses, courseMap)
}

func GetMailContent(updateDate *string, courses *[]Course, courseMap map[string]string) *string {
	var content string = "<html><head><meta http-equiv=\"Content-Type\" content=\"text/html; charset=iso-8859-9\"/></head><body><b>Güncelleme Tarihi: " + *updateDate + "</b>"
	for _, course := range *courses {
		_, exists := courseMap[course.Name]
		if exists || courseMap == nil {
			content += "<p>"
			if exists {
				delete(courseMap, course.Name)
			}
		} else {
			content += "<p style=\"background: #9EF7AC;\">"
		}
		content += "<b>" + strconv.Itoa(course.Index) + "- " + course.Name + "</b><br/>"
		content += strings.Join(course.Description, "<br/>") + "<br/>"
	}
	if len(courseMap) > 0 {
		content += "<br/><p style=\"font-weight: bold;\">Kontenjanı Dolan veya Kayıt Süresi Geçen Kurslar</p>"
	}
	index := 1
	re := regexp.MustCompile("^[0-9]+")
	for _, line := range courseMap {
		line = re.ReplaceAllString(line, strconv.Itoa(index))
		content += "<p style=\"background-color: #FF9494;\">" + line + "</p>"
		index++
	}
	content += "</body></html>"
	return &content
}

func RenderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}
