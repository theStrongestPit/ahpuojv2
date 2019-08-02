package utils

import (
	"encoding/base64"
	"encoding/xml"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/Unknwon/goconfig"
)

type TimeLimit struct {
	Unit    string `xml:"where,attr"`
	Content string `xml:",chardata"`
}

type MemoryLimit struct {
	Unit    string `xml:"where,attr"`
	Content string `xml:",chardata"`
}

type Solution struct {
	XMLName  xml.Name `xml:"solution"`
	Language string   `xml:"language,attr"`
	Content  string   `xml:",chardata"`
}

type Item struct {
	XMLName      xml.Name    `xml:"item"`
	Title        string      `xml:"title"`
	TimeLimit    TimeLimit   `xml:"time_limit"`
	MemoryLimit  MemoryLimit `xml:"memory_limit"`
	Description  string      `xml:"description"`
	Input        string      `xml:"input"`
	Output       string      `xml:"output"`
	SampleInput  string      `xml:"sample_input"`
	SampleOutput string      `xml:"sample_output"`
	TestInput    []string    `xml:"test_input"`
	TestOutput   []string    `xml:"test_output"`
	Hint         string      `xml:"hint"`
	Source       string      `xml:"source"`
	Solution     []Solution  `xml:"solution"`
}

type Fps struct {
	XMLName xml.Name `xml:"fps"`
	Item    []Item   `xml:"item"`
}

func Mkdata(pid int, filename string, input string) error {
	cfg, _ := goconfig.LoadConfigFile("config/config.ini")
	dataDir, _ := cfg.GetValue("project", "datadir")
	basedir := dataDir + "/" + strconv.Itoa(pid)
	file, err := os.Create(basedir + "/" + filename)
	if err != nil {
		Consolelog(err.Error())
	}
	Consolelog(basedir + "/" + filename)
	defer file.Close()
	if err != nil {
		return err
	}
	strings.Replace(input, "\n\r", "\n", -1)
	_, err = file.WriteString(input)
	return err
}
func ImageSaveFile(filePath string, base64EncodedImage string) error {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return err
	}
	base64DecodedImage, err := base64.StdEncoding.DecodeString(base64EncodedImage)
	file.Write(base64DecodedImage)
	return err
}

func ImportFps(tempfile multipart.File) (Fps, error) {
	defer tempfile.Close()
	fps := Fps{}
	data, _ := ioutil.ReadAll(tempfile)
	err := xml.Unmarshal(data, &fps)
	if err != nil {
		return fps, err
	}
	return fps, nil
}
