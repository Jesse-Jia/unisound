package main

import (
	"gree.com/utils"
	"net/http"
	"log"
	"gree.com/semantic"
	"io/ioutil"
	"encoding/json"
	"runtime"
	"fmt"
	"github.com/tidwall/gjson"
	"time"
	"strings"
)
// 日志文件存储位置
//var logfile = "/roobo/logs/gree.service/unisound/unisound_logfile_" + time.Now().Format("200601")
var logfile = "unisound_logfile_" + time.Now().Format("200601")
var logReqResp utils.LogData
// tencent fm(story) domain ("\"Gjson 会自动删除转义字符)
func fmDomain(jsonData []byte) ([]byte) {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	items := ""
	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items = `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}
	resultUrl := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.audio.stream.url")
	resultTitle := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.title")
	resultContent := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.textContent")
	resultText := gjson.GetBytes(jsonData, "payload.response_text")
	resultLen := len(resultUrl.Array())

	if resultLen == len(resultTitle.Array()) && resultLen == len(resultContent.Array()) {
		for i := 0; i < resultLen; i++ {
			if i == 0 {
				items = `{"url":`+ resultUrl.Array()[i].Raw +`,"title":`+ resultTitle.Array()[i].Raw +
					`,"content":` + resultContent.Array()[i].Raw + `}`
			} else {
				items += `,{"url":`+ resultUrl.Array()[i].Raw +`,"title":`+ resultTitle.Array()[i].Raw +
					`,"content":` + resultContent.Array()[i].Raw + `}`
			}
		}
	} else {
		items = `{"code":501,"errorType":"result length of url/title/content does not math!"}`
		return []byte(items)
	}
	jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultText.Raw + `,"listItems":[` + items + `]}`
	return []byte(jsonStr)
}
// tencent poem domain
func poemDomain(jsonData []byte) ([]byte) {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	items := ""
	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items = `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}

	resultUrl := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.htmlView")
	resultTitle := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.title")
	resultContent := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.textContent")
	resultText := gjson.GetBytes(jsonData, "payload.response_text")
	resultLen := len(resultUrl.Array())
	if resultLen == len(resultTitle.Array()) && resultLen == len(resultContent.Array()) {
		for i := 0; i < resultLen; i++ {
			if i == 0 {
				items = `{"url":`+ resultUrl.Array()[i].Raw +`,"title":`+ resultTitle.Array()[i].Raw +
					`,"content":` + resultContent.Array()[i].Raw + `}`
			} else {
				items += `,{"url":`+ resultUrl.Array()[i].Raw +`,"title":`+ resultTitle.Array()[i].Raw +
					`,"content":` + resultContent.Array()[i].Raw + `}`
			}
		}
	} else {
		items = `{"code":501,"errorType":"result length of url/title/content does not math!"}`
		return []byte(items)
	}

	//item := `{"url":"`+ resultUrl.String() +`","title":"`+ resultTitle.String() +
	//	`","content":"` + resultContent.String() + `"}`
	jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultText.Raw + `,"listItems":[`+ items +`]}`
	return []byte(jsonStr)
}
// tencent weather domain
func weatherDomain(jsonData []byte) ([]byte) {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	//植入广告
	ad := "格力空调祝您生活愉快。"

	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items := `{"code":500,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}

	//resultText := gjson.GetBytes(jsonData, "payload.response_text")
	//jsonStr := `{"header":` + header.String() + `,"response_text":"`+ resultText.String() + ad +`"}`
	resultContent := gjson.GetBytes(jsonData, "payload.data.json.listItems.0.textContent")
	test := resultContent.String() + ad
	jsonStr := `{"header":` + header.Raw + `,"response_text":"`+ test +`"}`
	return []byte(jsonStr)
}
// tencent music domain
func musicDomain(jsonData []byte) ([]byte, []semantic.TencentReport) {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	items := ""

	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items = `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items), nil
	}
	speakInfo := gjson.GetBytes(jsonData, "payload.data.json.controlInfo.textSpeak")
	resultUrl := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.audio.stream.url")
	resultSinger := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.selfData.singer")
	resultSong := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.selfData.song")
	resultText := gjson.GetBytes(jsonData, "payload.response_text")
	resultLen := len(resultUrl.Array())
	var reportDatas []semantic.TencentReport
	var reportData semantic.TencentReport
	resultMediaId := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.mediaId")
	if resultLen != 0 && resultLen == len(resultSinger.Array()) && resultLen == len(resultSong.Array()) {
		for i:=0; i < resultLen; i++ {
			if i == 0 {
				items = `{"url":`+ resultUrl.Array()[i].Raw +`,"singer":`+ resultSinger.Array()[i].Raw +
					`,"song":` + resultSong.Array()[i].Raw + `}`
			} else {
				items += `,{"url":`+ resultUrl.Array()[i].Raw +`,"singer":`+ resultSinger.Array()[i].Raw +
					`,"song":` + resultSong.Array()[i].Raw + `}`
			}
			reportData.UserId = "test"
			reportData.Domain = "music"
			reportData.Intent = gjson.GetBytes(jsonData, "header.semantic.intent").String()
			reportData.ResourceId = resultMediaId.Array()[i].String()
			reportData.DataSource = "..."
			if reportData.ResourceId != "" {
				reportDatas = append(reportDatas, reportData)
			}
			//reportStatus := semantic.ReportTencentStatus(reportData)
			//if nil == reportStatus {
			//	fmt.Println("report error!")
			//} else {
			//	fmt.Println(string(reportStatus))
			//}
		}
	} else {
		items = `{"code":501,"errorType":"result length of url/singer/song does not math!"}`
		return []byte(items), reportDatas
	}
	if speakInfo.String() == "true" {
		resultContent := gjson.GetBytes(jsonData, "payload.data.json.listItems.0.textContent")
		jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultContent.Raw + `,"listItems":[` + items + `]}`
		return []byte(jsonStr), reportDatas
	}
	jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultText.Raw + `,"listItems":[` + items + `]}`
	return []byte(jsonStr), reportDatas
}
// Tencent news domain
func newsDomain(jsonData []byte) ([]byte) {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	items := ""

	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items = `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}

	resultUrl := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.audio.stream.url")
	resultFrom := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.selfData.newsFrom")
	resultType := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.selfData.type")
	//resultMediaId := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.mediaId")
	//resultSource := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.selfData.source")
	//resultAbstract := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.selfData.newsAbstract")
	resultText := gjson.GetBytes(jsonData, "payload.response_text")
	resultTitle := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.title")
	resultLen := len(resultUrl.Array())
	//var reportData semantic.TencentReport
	//reportData.UserId = ""
	//reportData.Domain = "news"
	//reportData.Intent = gjson.GetBytes(jsonData, "header.semantic.intent").String()
	//reportData.ResourceId = ""
	//reportData.DataSource = ""
	if resultLen != 0 && resultLen == len(resultFrom.Array()) && resultLen == len(resultType.Array()) {
		for i:=0; i < resultLen; i++ {
			/*
			reportData.ResourceId = resultMediaId.Array()[i].String()
			reportData.DataSource = resultSource.Array()[i].String()
			reportStatus := semantic.ReportTencentStatus(reportData)
			if nil == reportStatus {
				fmt.Println("report error!")
			} else {
				fmt.Println(string(reportStatus))
			}
			*/
			if i == 0 {
				items = `{"url":`+ resultUrl.Array()[i].Raw +`,"newsFrom":`+ resultFrom.Array()[i].Raw +
					`,"type":` + resultType.Array()[i].Raw + `,"resultTitle":`+ resultTitle.Array()[i].Raw + `}`
			} else {
				items += `,{"url":`+ resultUrl.Array()[i].Raw +`,"newsFrom":`+ resultFrom.Array()[i].Raw +
					`,"type":` + resultType.Array()[i].Raw + `,"resultTitle":`+ resultTitle.Array()[i].Raw + `}`
			}
		}
	} else {
		items = `{"code":501,"errorType":"result length of url/singer/song does not math!"}`
		return []byte(items)
	}
	jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultText.Raw + `,"listItems":[` + items + `]}`
	return []byte(jsonStr)
}

// Tencent sports domain
func sportsDomain(jsonData []byte) ([]byte) {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	items := ""

	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items = `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}
	resultText := gjson.GetBytes(jsonData, "payload.response_text")
	jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultText.Raw + `}`
	return []byte(jsonStr)
}

// Tencent joke domain
func jokeDomain(jsonData []byte) ([]byte){
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	items := ""

	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items = `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}
	resultUrl := gjson.GetBytes(jsonData,"payload.data.json.listItems.#.audio.stream.url")
	for i, item := range resultUrl.Array() {
		if i == 0 {
			items = `{"url":`+ item.Raw + `}`
		} else {
			items += `,{"url":`+ item.Raw + `}`
		}
	}
	//resultMediaId := gjson.GetBytes(jsonData,"payload.data.json.sMongoNewId")
	resultText := gjson.GetBytes(jsonData, "payload.response_text")
	jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultText.Raw + `,"listItems":[` + items + `]}`
	return []byte(jsonStr)
}

// Tencent astro domain
func astroDomain(jsonData []byte) ([]byte) {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	items := ""

	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items = `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}
	resultText := gjson.GetBytes(jsonData, "payload.response_text")
	jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultText.Raw + `}`
	return []byte(jsonStr)
}

// Tencent holiday domain
func holidayDomain(jsonData []byte) ([]byte){
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	items := ""

	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items = `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}
	resultContent := gjson.GetBytes(jsonData, "payload.data.json.listItems.0.textContent")
	jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultContent.Raw + `}`
	return []byte(jsonStr)
}

// Tencent stock domain
func stockDomain(jsonData []byte) ([]byte) {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	items := ""

	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items = `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}
	//resultText := gjson.GetBytes(jsonData, "payload.response_text")
	resultContent := gjson.GetBytes(jsonData, "payload.data.json.listItems.0.textContent")
	jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultContent.Raw + `}`
	return []byte(jsonStr)
}

// Tencent translate domain
func translateDomain(jsonData []byte) ([]byte)  {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	items := ""

	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items = `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}
	resultContent := gjson.GetBytes(jsonData, "payload.data.json.listItems.0.textContent")
	jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultContent.Raw + `}`
	return []byte(jsonStr)
}
// Tencent sound baike domain
func soundDomain(jsonData []byte) ([]byte)  {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	items := ""

	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items = `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}
	resultText := gjson.GetBytes(jsonData, "payload.response_text")
	resultUrl := gjson.GetBytes(jsonData, "payload.data.json.listItems.#.audio.stream.url")
	for i:=0; i < len(resultUrl.Array()); i++  {
		if i == 0 {
			items = `{"url":`+ resultUrl.Array()[i].Raw + `}`
		} else {
			items += `,{"url":`+ resultUrl.Array()[i].Raw + `}`
		}
	}
	jsonStr := `{"header":` + header.String() + `,"response_text":`+ resultText.Raw + `,"listItems":[` + items + `]}`
	return []byte(jsonStr)
}
// Tencent almanac domain
func almanacDomain(jsonData []byte) ([]byte)  {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	items := ""

	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items = `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}
	resultText := gjson.GetBytes(jsonData, "payload.response_text")
	jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultText.Raw + `}`
	return []byte(jsonStr)
}
// Tencent finance domain
func financeDomain(jsonData []byte) ([]byte)  {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	items := ""

	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items = `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}
	resultText := gjson.GetBytes(jsonData, "payload.response_text")
	jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultText.Raw + `}`
	return []byte(jsonStr)
}
// Tencent food domain
func foodDomain(jsonData []byte) ([]byte) {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	items := ""

	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items = `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}
	resultText := gjson.GetBytes(jsonData, "payload.response_text")
	jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultText.Raw + `}`
	return []byte(jsonStr)
}
// Tencent general_question_answering domain 
func generalQADomain(jsonData []byte) ([]byte) {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	items := ""

	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items = `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}
	resultContent := gjson.GetBytes(jsonData, "payload.data.json.listItems.0.textContent")
	jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultContent.Raw + `}`
	return []byte(jsonStr)
}
// Tencent baike domain
func baikeDomain(jsonData []byte) ([]byte) {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	items := ""

	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items = `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}
	//resultGInfo := gjson.GetBytes(jsonData, "payload.data.json.globalInfo.seeMore")
	resultListItems := gjson.GetBytes(jsonData, "payload.data.json.listItems")
	//resultText := gjson.GetBytes(jsonData, "payload.response_text")
	resultContent := gjson.GetBytes(jsonData, "payload.data.json.listItems.0.textContent")
	index := strings.Index(resultContent.Raw, "。")
	test := string(resultContent.Raw[0:index]) + "\""
	jsonStr := `{"header":` + header.Raw + `,"response_text":`+ test + `,"listItems":`+ resultListItems.Raw +`}`
	return []byte(jsonStr)
}
// Tencent chenyu domain
func chengyuDomain(jsonData []byte) ([]byte) {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	items := ""

	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items = `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}
	//resultText := gjson.GetBytes(jsonData, "payload.response_text")
	resultContent := gjson.GetBytes(jsonData, "payload.data.json.listItems.0.textContent")
	jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultContent.Raw + `}`
	return []byte(jsonStr)
}
// Tencent science domain
func scienceDomain(jsonData []byte) ([]byte) {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items := `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}
	resultText := gjson.GetBytes(jsonData, "payload.response_text")
	jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultText.Raw + `}`
	return []byte(jsonStr)
}
// Tencent recipe domain
func recipeDomain(jsonData []byte) ([]byte) {
	header := gjson.GetBytes(jsonData, "header")
	//resultText := gjson.GetBytes(jsonData, "payload.response_text")
	resultText := "一道好菜很难通过三言两语描述清楚，建议您浏览相关菜谱网站。"
	jsonStr := `{"header":` + header.Raw + `,"response_text":"`+ resultText + `"}`
	return []byte(jsonStr)
}
// Tencent chat domain
func chatDomain(jsonData []byte) ([]byte) {
	header := gjson.GetBytes(jsonData, "header")
	payload := gjson.GetBytes(jsonData, "payload")
	if !header.Exists() || !payload.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items := `{"code":501,"errorType":"result header/payload dose not exists!"}`
		return []byte(items)
	}
	code := gjson.GetBytes(jsonData, "header.semantic.code").String()
	if code == "0" {
		resultContent := gjson.GetBytes(jsonData, "payload.data.json.listItems.0.textContent")
		jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultContent.Raw + `}`
		return []byte(jsonStr)
	} else {
		resultText := gjson.GetBytes(jsonData, "payload.response_text")
		jsonStr := `{"header":` + header.Raw + `,"response_text":`+ resultText.Raw + `}`
		return []byte(jsonStr)
	}
}
// unsupported domain
func otherDomain(jsonData []byte) ([]byte) {
	header := gjson.GetBytes(jsonData, "header")
	resultText := "unsupported"
	jsonStr := `{"header":` + header.Raw + `,"response_text":"`+ resultText + `"}`
	return []byte(jsonStr)
}
// module test 其中jsonData为提取后的格式
// 语音板解析http数据包的格式非常固定，json内容不要改，否则可能出现解析异常错误
func moduleTest(jsonData []byte, query string) ([]byte) {
	header := gjson.GetBytes(jsonData, "header")
	response_text := gjson.GetBytes(jsonData, "response_text")
	listItems := gjson.GetBytes(jsonData, "listItems")
	if !header.Exists() || !response_text.Exists() {
		//fmt.Println("json data header/payload dose not exists!")
		items := `{"code":501,"errorType":"result header/response_text dose not exists!"}`
		return []byte(items)
	} else if response_text.String() == "unsupported" {
		items := `{"code":501,"errorType":"unsupported domain!"}`
		return []byte(items)
	} else if len(listItems.Array()) > 0 {
		jsonStr := `{"header":` + header.Raw + `,"response_text":`+ response_text.Raw +`,"asr_recongize":"`+
			query +`","listItems":`+ listItems.Raw +`}`
		return []byte(jsonStr)
	} else {
		jsonStr := `{"header":` + header.Raw + `,"response_text":` + response_text.Raw + `,"asr_recongize":"` + query + `"}`
		return []byte(jsonStr)
	}
}

/*
func moduleTest(jsonData []byte, query string) ([]byte) {
	var jsonMap map[string]interface{}
	json.Unmarshal(jsonData, &jsonMap)
	if jsonMap != nil {
		jsonMap["asr_recongize"] = query
	}
	jsonStr, _ := json.Marshal(jsonMap)
	return []byte(jsonStr)
}
*/
func ParseGreeJson(jsonData []byte) ([]byte, error) {
	var structData semantic.SemanticResp
	err := json.Unmarshal(jsonData, &structData)
	if nil != err {
		return nil, err
	}
	switch structData.Semantic.Service {
	case "Aftersales":
		structData.Result.Hint = "您可以拨打格力售后电话400-836-5315，或扫面机身二维码预约维修。"
	case "eQuery":
		structData.Result.Hint = "您的电费数据尚未导入"
	case "Company":
		structData.Result.Hint = "让世界爱上中国造"
	case "Product":
		structData.Result.Hint = "我们不生产空调，我们只是大自然空气能源的搬运工"
	default:
		//responseStructData.Asr_recongize = ""
	}
	return json.Marshal(structData)
}
func ParseTencentJson(jsonData []byte) ([]byte, []semantic.TencentReport) {
	domain := gjson.GetBytes(jsonData, "header.semantic.domain")
	//items := ""
	var jsonByte []byte
	var reportDatas []semantic.TencentReport
	// 音乐
	// 新闻
	// 天气
	// 唐诗
	// 故事
	// 体育
	// 笑话
	// 星座
	// 节日
	// 股市
	// 翻译
	// 黄历
	// 财经
	// 食物
	// 一般问答
	// 百科
	// 成语
	// 科学
	// 菜谱
	// 聊天
	// 其他
	if domain.String() == "music" {
		//jsonByte = musicDomain(jsonData)
		jsonByte,reportDatas = musicDomain(jsonData)
	} else if domain.String() == "news" {
		jsonByte = newsDomain(jsonData)
	} else if domain.String() == "weather" {
		jsonByte = weatherDomain(jsonData)
	} else if domain.String() == "ancient_poem" {
		jsonByte = poemDomain(jsonData)
	} else if domain.String() == "fm" {
		jsonByte = fmDomain(jsonData)
	} else if domain.String() == "sports" {
		jsonByte = sportsDomain(jsonData)
	} else if domain.String() == "joke" {
		jsonByte = jokeDomain(jsonData)
	} else if domain.String() == "astro" {
		jsonByte = astroDomain(jsonData)
	} else if domain.String() == "holiday" {
		jsonByte = holidayDomain(jsonData)
	} else if domain.String() == "stock" {
		jsonByte = stockDomain(jsonData)
	} else if domain.String() == "translate" {
		jsonByte = translateDomain(jsonData)
	} else if domain.String() == "sound" {
		jsonByte = soundDomain(jsonData)
	} else if domain.String() == "almanac" {
		jsonByte = almanacDomain(jsonData)
	} else if domain.String() == "finance" {
		jsonByte = financeDomain(jsonData)
	} else if domain.String() == "food" {
		jsonByte = foodDomain(jsonData)
	} else if domain.String() == "general_question_answering" {
		jsonByte = generalQADomain(jsonData)
	} else if domain.String() == "baike" {
		jsonByte = baikeDomain(jsonData)
	} else if domain.String() == "chengyu" {
		jsonByte = chengyuDomain(jsonData)
	} else if domain.String() == "science" {
		jsonByte = scienceDomain(jsonData)
	} else if domain.String() == "recipe" {
		jsonByte = recipeDomain(jsonData)
	} else if domain.String() == "chat" {
		jsonByte = chatDomain(jsonData)
	} else {
		// other domains
		//return jsonData
		jsonByte = otherDomain(jsonData)
	}
	//fmt.Println(reportDatas)

	//for _, reportData := range reportDatas {
		//腾讯上报数据
	//	reportStatus := semantic.ReportTencentStatus(reportData)
	//	fmt.Println(string(reportStatus))

	//}


	return jsonByte, reportDatas
}

func getHTTPReqBody(req *http.Request) (*semantic.SemanticReq, error) {
	// get client request HTTPS body data
	clientReqBodyBytes, err := ioutil.ReadAll(req.Body)
	if nil != err {
		log.Println("service: Get request body data error")
		return nil, err
	}
	fmt.Println(string(clientReqBodyBytes))
	logdata := utils.RemoveSpaceAndLineBreaks(string(clientReqBodyBytes))
	/*
		if err := utils.LogStoreReqToFile(logfile, logdata); nil != err {
			log.Printf("Store access data to file error: %s", err)
			return nil, err
		}
	*/
	logReqResp.RequestData = logdata

	// clientReqBodyStruct variable store access semantic platform data.
	var clientReqBodyStruct semantic.SemanticReq
	if err := json.Unmarshal(clientReqBodyBytes, &clientReqBodyStruct); nil != err {
		log.Println("service: Unmarshal request data error")
		log.Println(err)
		return nil, err
	}
	return &clientReqBodyStruct, nil
}

func handleService(w http.ResponseWriter, req *http.Request)  {
	var reportDatas []semantic.TencentReport
	reqTime := time.Now()
	// accessSemanticData variable store access semantic platform data.
	accessSemanticData, err := getHTTPReqBody(req)
	if nil != err {
		log.Println("service function: Get client request data error")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.Write([]byte(`{"code":501,"errorType":"request data error, check it and try again!"}`))
		return
	}
	// check uid & token
	if accessSemanticData.UID != "unisound" ||
		accessSemanticData.Token != "9ff9874dd2f8b6d9e0343c22c23f4248543eec156303703b42a38488e581be42" {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.Write([]byte(`{"code":501,"errorType":"request unauthorized, check it and try again!"}`))
		return
	}
	// check timestamp
	/*
	timeStamp := accessSemanticData.TimeStamp
	recTime, err := time.ParseInLocation("2006-01-02 15:04:05", timeStamp, time.Local)

	if reqTime.Unix() - recTime.Unix()  > 60.0 {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.Write([]byte(`{"code":501,"errorType":"time expire!"}`))
		return
	} else if reqTime.Unix() - recTime.Unix()  < -30.0 {
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.Write([]byte(`{"code":501,"errorType":"time illegal!"}`))
		return
	}
	*/
	// request remote IP
	if accessSemanticData.Ip == "" {
		//fmt.Println(req.RemoteAddr)
		//fmt.Println(req.Header.Get("Remote_addr"))
		accessSemanticData.Ip = req.Header.Get("Remote_addr")
	}
	//semanticRespByteData, err := semantic.GetSemanticData(accessSemanticData.Query)
	//
	var semanticRespByteData []byte
	reqQuery := accessSemanticData.Query
	if accessSemanticData.Classify == "gree" || accessSemanticData.Classify == "Gree" {
		greeRespByteData, err := semantic.GetGreeNLUData(reqQuery)
		if nil != err {
			log.Println("HandleWearher: Get semantic platform Gree data error")
			log.Println(err)
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			w.Write([]byte(`{"code":501,"errorType":"transfer data error, check it and try again!"}`))
			return
		}
		semanticRespByteData, err = ParseGreeJson(greeRespByteData)
	} else if accessSemanticData.Classify == "tencent" || accessSemanticData.Classify == "Tencent" {
		//semanticRespByteData, err = semantic.GetTencentNLUData(accessSemanticData.Query)
		//tencentRespByteData, err := semantic.GetTencentNLUData(accessSemanticData.Query)
		tencentRespByteData, err := semantic.GetTencentNLUData(accessSemanticData)
		if nil != err {
			log.Println("HandleWearher: Get semantic platform Tencent data error")
			log.Println(err)
			w.Header().Set("Content-Type", "application/json;charset=utf-8")
			w.Write([]byte(`{"code":501,"errorType":"transfer data error, check it and try again!"}`))
			return
		}
		semanticRespByteData, reportDatas = ParseTencentJson(tencentRespByteData)
		// for module test
		semanticRespByteData = moduleTest(semanticRespByteData, reqQuery)

	} else {
		log.Println("HandleService: no such classify!")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.Write([]byte(`{"code":501,"errorType":"classify error, check it and try again!"}`))
	}
	/*
	if nil != err {
		log.Println("HandleWearher: Get semantic platform ROS.AI data error")
		log.Println(err)
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.Write([]byte(`{"code":501,"errorType":"transfer data error, check it and try again!"}`))
		return
	}
	*/
	//fmt.Println(string(semanticRespByteData))
	/*
	var semanticRespStructData []byte
	if err := json.Unmarshal(semanticRespByteData, &semanticRespStructData); nil != err {
		log.Println("Unmarshal semantic data error")
		log.Println(err)
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.Write([]byte(`{"code":501,"errorType":"response data error, check it and try again!"}`))
		return
	}
	*/
	// return semantic platform response data.
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(semanticRespByteData)

	// 数据上报
	if len(reportDatas) != 0 {
		/*
		for _, reportData := range reportDatas {
			//腾讯上报数据
			reportStatus := semantic.ReportTencentStatus(reportData)
			fmt.Println(string(reportStatus))
		}
		*/
	} else {
		fmt.Println("nothing to report!")
	}

	// write logfile
	resTime := time.Now()
	costTime := resTime.Sub(reqTime)
	logReqResp.CostTime = fmt.Sprint(costTime)
	logReqResp.ResponseData = utils.RemoveSpaceAndLineBreaks(string(semanticRespByteData))

	utils.StoreToLogfile(logfile, logReqResp)
}

func main()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//调试的路径，服务器需要更改--云知声http://api.gree.com:8088/unisound/v1/query -> http://10.2.45.70:9997
	http.HandleFunc("/", handleService)
	log.Fatalln(http.ListenAndServe(":9997", nil))

}