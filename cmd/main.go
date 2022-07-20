package main

import "loante/service/cron"

func main() {
	CronInit()
}

func CronInit() {
	cron.ScoreCron()
}
