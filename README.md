# Alert-cli

### Features
检测监控问题并通过企业微信将问题信息告警给业务方。
目前支持的功能：
* 检测不收敛的 metrics，并统计 label 信息发送给业务方  

### Getting start
```shell script
$ git clone https://github.com/YiyongHuang/Alert-cli.git
```
```shell script
# 检测 metrics 问题
$ go run main.go metrics
--report-path=<url> # 企业微信通知接口
--report-backup-path=<url> # 企业微信备用通知接口
--thanos-query-url=<url> # thanos query 服务路径
--service-path=<url> # 查询服务对应项目信息的接口（主要用于查询通知接收人）
```