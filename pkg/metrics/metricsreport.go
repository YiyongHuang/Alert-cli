package metrics

import (
	"fmt"
	"time"

	"github.com/YiyongHuang/Alert-cli/pkg/utils"
	"github.com/prometheus/common/model"
	"k8s.io/klog"
)

const (
	labelThreshold = 2
	prometheusDoc  = "https://prometheus.io/docs/prometheus/latest/getting_started/"
	threshold      = 2000
	sleepPeriod    = 24
)

type MetricsCli struct {
}

func (m *MetricsCli) buildMetricMessage(metrics map[string]int, labels *model.Sample) string {
	message := fmt.Sprintf("### 服务单 pod 单条 prometheus metrics 时间序列条数不符合规范(大于 `%d` 条)，详情如下：\n\n", threshold)

	message = fmt.Sprintf("%s>**PAAS 项目名:** %s\n", message, labels.Metric["opsservice"])
	message = fmt.Sprintf("%s>**Pod 名:** %s\n", message, labels.Metric["pod"])
	message = fmt.Sprintf("%s>**Metrics 名:** %s\n", message, labels.Metric["__name__"])
	message = fmt.Sprintf("%s>**总时间序列条数:** %f\n\n", message, labels.Value)
	message = fmt.Sprintf("%s>**标签统计详情如下:**\n", message)
	for k, v := range metrics {
		message = fmt.Sprintf("%s>**%s:** %d\n", message, k, v)
	}
	message = fmt.Sprintf("%s\n>使用 qms 、pedestal 框架请联系***协助升级框架！\n", message)
	message = fmt.Sprintf("%s>未使用框架请参考 [prometheus文档](%s) 收敛 prometheus metrics 标签值 ！\n", message, prometheusDoc)
	return message
}

func (m *MetricsCli) getBadMetrics(vector model.Vector) map[string]int {
	var badMetrics = make(map[string]int)
	metricsMap := make(map[string]map[string]int32)

	for _, item := range vector {
		for k, v := range item.Metric {
			if _, ok := metricsMap[string(k)]; !ok {
				metricsMap[string(k)] = make(map[string]int32)
			}
			if _, ok := metricsMap[string(k)][string(v)]; !ok {
				metricsMap[string(k)][string(v)] = 1
			}
		}
	}

	for k, v := range metricsMap {
		if len(v) > labelThreshold {
			badMetrics[k] = len(v)
		}
	}

	return badMetrics
}

func (m *MetricsCli) HandleMetrics(cfg *MetricCfg) {
	promClient, err := utils.NewPromClient(cfg.ThanosQueryURL)
	if err != nil {
		klog.Errorf("New promClient error, quit alert cli ...")
		return
	}

	for {
		resultVector, err := utils.Query(promClient, utils.SeriesToMany(threshold))
		if err != nil {
			klog.Errorf("msg: ", "get SeriesToMany error, quit!")
			break
		}
	
		if len(resultVector) == 0 {
			klog.Errorf("msg: ", "null SeriesToMany result vector, quit!")
			break
		}
	
		klog.Info("msg: ", fmt.Sprintf("%d metrics need to process", len(resultVector)))
		for _, item := range resultVector {
			result, err := utils.Query(promClient, item.Metric.String())
			if err != nil {
				klog.Errorf("get metrics error")
				continue
			}
	
			metrics := m.getBadMetrics(result)
			message := m.buildMetricMessage(metrics, item)
			body, err := utils.GetAlertBody(message, cfg.ServicePath, item)
			if err != nil {
				klog.Errorf(fmt.Sprintf("get alert body err: %v", err))
				continue
			}
	
			if err := utils.SendMessage(body, cfg.ReportPath); err != nil {
				klog.Errorf(fmt.Sprintf("send message [%s] err: %v, try backup url again", item.Metric.String(), err))
				if err := utils.SendMessage(body, cfg.ReportPathBak); err != nil {
					klog.Errorf(fmt.Sprintf("try err again: %v", err))
				} else {
					klog.Info("msg: ", "try again success !")
				}
			}
		}

		time.Sleep(sleepPeriod*time.Hour)
	}

	klog.Info("msg: ", "metrics report quit !")
	return
}
