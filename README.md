[golang实现每天自动签到和免费抽奖](https://juejin.cn/post/7047311719659470885)

go run main.go "cookie值"  普通运行  
main_scf.go  云函数运行

#### workflow定时任务
* 仓库Settings->secrets>点击New repository secret按钮新建Name:JUEJIN_COOKIE,
* 在输入框中粘贴juejin网址已登录后的cookie值，cookie值开头例如：MONITOR_WEB_ID=a17...
* 点击Actions->Go->Run workflow按钮测试

> [go.yml](https://github.com/jangworn/juejin-checkin/blob/master/.github/workflows/go.yml) 中cron的值`0 2 * * *`表示每天02:00 UTC运行，实际测试中会有延迟北京时间上午11点多才运行

active 20240709
