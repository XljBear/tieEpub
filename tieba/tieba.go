package tieba

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-shiori/go-epub"
	"github.com/gocolly/colly/v2"
	"github.com/spf13/viper"
)

func extractPostID(url string) (string, error) {
	// 定义匹配帖子 ID 的正则表达式
	re := regexp.MustCompile(`https://tieba\.baidu\.com/p/(\d+)`)

	// 查找匹配项
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("未找到帖子 ID")
	}

	return matches[1], nil
}

type TieRequest struct {
	Url         string
	MinimumWord int
	OnlyLZ      bool
	FilterLink  bool
	FilterImg   bool
	ErrorChan   chan string
	ProcessChan chan int
	SuccessChan chan int
}

var totalContent []string
var title = ""
var author = ""
var coverData = ""

type TieContent struct {
	TotalContent []string
	Title        string
	Author       string
}

func GetTieData() *TieContent {
	if len(totalContent) == 0 || title == "" || author == "" {
		return nil
	}
	return &TieContent{
		TotalContent: totalContent,
		Title:        title,
		Author:       author,
	}
}
func removeImgTags(html string) string {
	re := regexp.MustCompile(`<img\b[^>]*>`)
	return re.ReplaceAllString(html, "")
}
func removeHrefTags(html string) string {
	re := regexp.MustCompile(`<a\b[^>]*>(.*?)</a>`)
	return re.ReplaceAllString(html, "")
}
func GetTie(tieRequest *TieRequest) {

	title = ""
	author = ""
	coverData = ""
	totalContent = make([]string, 0)

	cookie := viper.Get("cookie").(string)
	id, err := extractPostID(tieRequest.Url)
	if err != nil || id == "" {
		tieRequest.ErrorChan <- "请输入正确贴吧的链接"
		return
	}
	totalPage := 1
	nowPage := 1
	errFlag := false
	for {
		if nowPage > totalPage || errFlag {
			break
		}
		urlStr := "https://tieba.baidu.com/p/" + id + "?pn=" + strconv.Itoa(nowPage)
		if tieRequest.OnlyLZ {
			urlStr += "&see_lz=1"
		}
		fmt.Println(urlStr)
		c := colly.NewCollector()
		c.OnHTML(".core_title_txt:first-of-type", func(e *colly.HTMLElement) {
			if title != "" {
				return
			}
			title = e.Attr("title")
		})
		c.OnHTML(".louzhubiaoshi:first-of-type", func(e *colly.HTMLElement) {
			if author != "" {
				return
			}
			author = e.Attr("author")
		})
		c.OnHTML(".pb_list_pager a:last-child", func(e *colly.HTMLElement) {
			lastPageLink := e.Attr("href")
			u, _ := url.Parse(lastPageLink)
			query := u.RawQuery
			parsedQuery, _ := url.ParseQuery(query)
			totalPage, _ = strconv.Atoi(parsedQuery.Get("pn"))
		})
		c.OnHTML(".j_d_post_content", func(e *colly.HTMLElement) {
			contentText := e.Text
			contentText = strings.TrimSpace(contentText)
			contentText = strings.Replace(contentText, "<br>", "\n", -1)
			contentText = removeImgTags(contentText)
			chineseLen := len([]rune(contentText))
			if chineseLen < tieRequest.MinimumWord {
				return
			}
			if tieRequest.FilterLink && strings.Contains(contentText, "</a>") {
				return
			}
			totalContent = append(totalContent, contentText)
		})
		c.OnError(func(r *colly.Response, err error) {
			tieRequest.ErrorChan <- "发生错误，可能被百度屏蔽。请尝试更新cookie"
			errFlag = true
		})
		c.OnRequest(func(r *colly.Request) {
			r.Headers.Set("cookie", cookie)
		})
		c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"
		c.Visit(urlStr)
		tieRequest.ProcessChan <- int(float32(nowPage) / float32(totalPage) * 100)
		nowPage += 1
		time.Sleep(1 * time.Second)
	}
	if errFlag {
		return
	}
	tieRequest.SuccessChan <- 1
}

func SaveEpub(savePath string, enableCover bool) error {
	e, err := epub.NewEpub(title)
	if err != nil {
		return err
	}
	e.SetAuthor(author)

	if enableCover && coverData != "" {
		coverPath, err := e.AddImage(coverData, "cover.png")
		if err != nil {
			return err
		}
		err = e.SetCover(coverPath, "")
		if err != nil {
			return err
		}
	}

	cssData := "data:text/css;base64,Ym9keSB7CiAgICBmb250LWZhbWlseTogQXJpYWwsIHNhbnMtc2VyaWY7CiAgICBsaW5lLWhlaWdodDogMS42OwogICAgbWFyZ2luOiAyZW07Cn0KaDEgewogICAgY29sb3I6ICMzMzM7Cn0KcCB7CiAgICB0ZXh0LWFsaWduOiBqdXN0aWZ5Owp9"
	cssPath, err := e.AddCSS(cssData, "styles.css")
	if err != nil {
		return err
	}
	for i, chapterContent := range totalContent {
		chapterContent = strings.ReplaceAll(chapterContent, "&nbsp;", " ")
		chapterContent = strings.ReplaceAll(chapterContent, "\u0018", "")
		chapterName := fmt.Sprintf("第%d章", i+1)
		addChapter(e, chapterName, chapterContent, cssPath)
	}

	// 把title中的特殊字符去掉
	//title = strings.ReplaceAll(title, "/", "-")
	//title = strings.ReplaceAll(title, ":", "-")
	//title = strings.ReplaceAll(title, "*", "-")
	//title = strings.ReplaceAll(title, "?", "-")
	//title = strings.ReplaceAll(title, "\"", "-")

	// 保存EPUB文件
	err = e.Write(savePath)
	if err != nil {
		return err
	}
	return nil
}
func StartGetAiImg(keyword string, chapterIndex int) (imgBase64 string, err error) {
	apiKey := viper.GetString("ai-api-key")
	if apiKey == "" {
		err = errors.New("请先完成API Key的设置配置")
		return
	}
	if chapterIndex > len(totalContent) {
		err = errors.New("AI创作封面章节选择错误")
		return
	}
	randSeed := rand.Int32()
	keyword += "。请根据下面内容绘制合适的场景，该图像将作为书籍封面，请按照书籍封面的风格来绘制。"
	keyword += removeHrefTags(removeImgTags(totalContent[chapterIndex]))
	fmt.Println(keyword)
	apiUrl := "https://api.siliconflow.cn/v1/images/generations"
	method := "POST"
	keyword = strings.Replace(keyword, "&nbsp;", " ", -1)
	keyword = strings.Replace(keyword, "\r", "", -1)
	keyword = strings.Replace(keyword, "\n", "", -1)

	type requestBody struct {
		Model             string  `json:"model"`
		Prompt            string  `json:"prompt"`
		NegativePrompt    string  `json:"negative_prompt"`
		ImageSize         string  `json:"image_size"`
		BatchSize         int     `json:"batch_size"`
		Seed              int32   `json:"seed"`
		NumInferenceSteps int     `json:"num_inference_steps"`
		GuidanceScale     float64 `json:"guidance_scale"`
	}

	reqBody := requestBody{
		Model:             "Kwai-Kolors/Kolors",
		Prompt:            keyword,
		NegativePrompt:    "",
		ImageSize:         "768x1024",
		BatchSize:         1,
		Seed:              randSeed,
		NumInferenceSteps: 25,
		GuidanceScale:     7.5,
	}
	payloadData, _ := json.Marshal(reqBody)
	payload := bytes.NewReader(payloadData)

	client := &http.Client{}
	req, err := http.NewRequest(method, apiUrl, payload)

	if err != nil {
		return
	}
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	if res.StatusCode != 200 {
		fmt.Println(string(body))
		err = errors.New("服务错误，请检查API Key是否配置正确")
		return
	}
	resp := &aiImgResponseData{}
	err = json.Unmarshal(body, resp)
	fmt.Println(resp)
	if resp.Code != 0 {
		err = errors.New(resp.Message)
		return
	}
	imgBase64, err = ImageURLToBase64(resp.Images[0].Url)
	coverData = imgBase64
	return
}
func ImageURLToBase64(imageURL string) (string, error) {
	// 创建HTTP客户端，设置超时时间
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 发送GET请求获取图片
	resp, err := client.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("获取图片失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP错误: %s", resp.Status)
	}

	// 读取图片数据
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取图片数据失败: %v", err)
	}

	// 获取Content-Type
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		// 如果响应头没有Content-Type，根据URL后缀猜测
		if strings.HasSuffix(strings.ToLower(imageURL), ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(strings.ToLower(imageURL), ".jpg") ||
			strings.HasSuffix(strings.ToLower(imageURL), ".jpeg") {
			contentType = "image/jpeg"
		} else if strings.HasSuffix(strings.ToLower(imageURL), ".gif") {
			contentType = "image/gif"
		} else if strings.HasSuffix(strings.ToLower(imageURL), ".webp") {
			contentType = "image/webp"
		} else {
			contentType = "image/jpeg" // 默认值
		}
	}

	// 转换为Base64
	base64String := base64.StdEncoding.EncodeToString(imageData)

	// 构建Data URL格式
	dataURL := fmt.Sprintf("data:%s;base64,%s", contentType, base64String)

	return dataURL, nil
}

type aiImgResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Images  []struct {
		Url string `json:"url"`
	} `json:"images"`
	Timings struct {
		Inference float64 `json:"inference"`
	} `json:"timings"`
	Seed     int64  `json:"seed"`
	SharedId string `json:"shared_id"`
	Data     []struct {
		Url string `json:"url"`
	} `json:"data"`
	Created int `json:"created"`
}

func addChapter(e *epub.Epub, title, content, cssPath string) {
	contents := strings.Split(content, "\n")
	contentHtml := ""
	for _, c := range contents {
		contentHtml += fmt.Sprintf("<p>%s</p>", c)
	}
	htmlContent := fmt.Sprintf(`
		<h1>%s</h1>
		%s
	`, title, contentHtml)
	_, err := e.AddSection(htmlContent, title, "", cssPath)
	if err != nil {
		log.Fatal(err)
	}
}
