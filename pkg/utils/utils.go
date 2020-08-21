package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	clientapi "github.com/prometheus/client_golang/api"
	promapi "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"k8s.io/klog"
)

const (
	timeout = 100
	//sendMessageDelay = 1
)

type weChatResponse struct {
	Data         int32  `json:"data"`
	Code         int32  `json:"code"`
	BusinessCode int32  `json:"businessCode"`
	Message      string `json:"message"`
}

type projectInfo struct {
	Project project `json:"project"`
}

type project struct {
	DepartmentMember []departmentMember `json:"department_member"`
}

type departmentMember struct {
	Username string `json:"username"`
}

func NewPromClient(ThanosQueryURL string) (promapi.API, error) {
	client, err := clientapi.NewClient(clientapi.Config{Address: ThanosQueryURL})
	if err != nil {
		klog.Errorf("new prometheus client err: %v", err)
		return nil, err
	}

	promClient := promapi.NewAPI(client)

	return promClient, nil
}

func Query(promClient promapi.API, promql string) (model.Vector, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	result, _, err := promClient.Query(ctx, promql, time.Now())
	if err != nil {
		klog.Errorf(fmt.Sprintf("prom api query err: %v", err))
		return nil, err
	}

	resultVector, ok := result.(model.Vector)
	if !ok {
		klog.Errorf("type of result err")
		return nil, err
	}

	return resultVector, nil
}

func SendMessage(body []byte, url string) error {
	client := &http.Client{
		Timeout: time.Second * 2,
	}

	postReq, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("new post request err: %w", err)
	}

	postResponse, err := client.Do(postReq)
	if err != nil {
		return fmt.Errorf("http client do post request err: %w", err)
	}
	defer postResponse.Body.Close()

	b, err := ioutil.ReadAll(postResponse.Body)
	if err != nil {
		return fmt.Errorf("ioutil read response body err: %w", err)
	}
	var weChatRes weChatResponse
	if err := json.Unmarshal(b, &weChatRes); err != nil {
		return fmt.Errorf("unmarshal weChat response body err: %w", err)
	}
	if weChatRes.Code != 200 {
		return fmt.Errorf("request interface err: [response code: %d, message: %s]", weChatRes.Code, weChatRes.Message)
	}

	//time.Sleep(time.Second * sendMessageDelay)
	return nil
}

func GetAlertBody(message, ServicePath string, labels *model.Sample) ([]byte, error) {
	members, err := getTos(string(labels.Metric["opsservice"]), ServicePath)
	if err != nil {
		return nil, fmt.Errorf("get Tos err: %w", err)
	}

	var postBody MessageBody
	if err := postBody.BuildMessageBody(message, members, []string{}); err != nil {
		return nil, fmt.Errorf("build message error: %w", err)
	}

	body, err := json.Marshal(postBody)
	if err != nil {
		return nil, fmt.Errorf("marshal post body err: %w", err)
	}

	return body, nil
}

// 查询对应服务接口获取对应接收通知的业务人员
func getTos(opsservice, ServicePath string) ([]string, error) {
	client := &http.Client{
		Timeout: time.Second * 2,
	}

	var projectInfo projectInfo
	var words []string
	var pre = 0
	for i := 0; i < len(opsservice); i++ {
		if opsservice[i] == '-' || opsservice[i] == '_' {
			words = append(words, opsservice[pre:i])
			pre = i + 1
		} else if i == (len(opsservice) - 1) {
			words = append(words, opsservice[pre:i+1])
		}
	}

	var opsUrls []string
	for _, item := range words {
		if len(opsUrls) == 0 {
			opsUrls = append(opsUrls, item)
		} else {
			for _, v := range opsUrls {
				opsUrls = append(opsUrls, v)
			}
			for idx := 0; idx < len(opsUrls); idx++ {
				if idx < len(opsUrls)/2 {
					opsUrls[idx] = opsUrls[idx] + "-" + item
				} else {
					opsUrls[idx] = opsUrls[idx] + "_" + item
				}
			}
		}
	}

	for _, opsUrl := range opsUrls {
		getUrl := fmt.Sprintf("%s%s", ServicePath, opsUrl)
		getReq, err := http.NewRequest("GET", getUrl, nil)
		if err != nil {
			return nil, fmt.Errorf("new get request err: %w", err)
		}

		getResponse, err := client.Do(getReq)
		if err != nil {
			return nil, fmt.Errorf("http client do get request err: %w", err)
		}

		getBody, err := ioutil.ReadAll(getResponse.Body)
		if err != nil {
			getResponse.Body.Close()
			return nil, fmt.Errorf("ioutil read response body err: %w", err)
		}

		if err := json.Unmarshal(getBody, &projectInfo); err != nil {
			getResponse.Body.Close()
			return nil, fmt.Errorf("unmarshal service %s project info err: %w", opsservice, err)
		}
		if len(projectInfo.Project.DepartmentMember) != 0 {
			getResponse.Body.Close()
			break
		}
	}
	var members []string
	//members = append(members, "huangyiyong")
	for _, usr := range projectInfo.Project.DepartmentMember {
		members = append(members, usr.Username)
	}
	if len(members) == 0 {
		return nil, fmt.Errorf("%s no member to send ", opsservice)
	}
	return members, nil
}
