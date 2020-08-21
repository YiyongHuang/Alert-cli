package utils

import "strconv"

var seriesToMany = `count by (opsservice,pod,__name__)({__name__=~".+",prometheus_replica=~"prom-biz-.*",opsservice!=""}) > `

//qa env ： var SeriesToMany = `count by (opsservice,pod,__name__)({__name__=~".+",prometheus_replica=~"prom.*",opsservice!=""}) > `

func SeriesToMany(threshold int) string {
	return seriesToMany + strconv.Itoa(threshold)
}
