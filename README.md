# vRealize Operations Team Slack Bot

The vR Ops Slack Bot runs in the configured (#vrops-cat) channel on https://vmware-vrops.slack.com.

## Pull Commands

### help

* `help`

### cat

* `cat recommended [branch]`
* `cat build [branch] [target] [result]`
* `cat testrun [branch] [area] [result]`

### vro

* `vro workflow` (`vro wf`)

### admin

* `admin uptime`

### hidden

* `cat areas`
* `cat branches`
* `cat results`
* `cat slas`
* `cat targets`
* `thank you`

### todo

* `admin config`
* `admin logs`
* `admin reset`
* `admin stats`

## Push Plugins

These plugins publish results periodically at a configured rate to a configured slack channel.

* recommended build change
* build status change - i.e. vrops, vcopssuitevm, vcopssuitepak
* test status change - i.e. SingleVA_PostCheckin, MultiVA_PostCheckin, SingleVA_B2B_PostCheckin

## Reference

* CAT REST API - https://cat.eng.vmware.com/urls
* MBU CAT Home - https://mbu-cat.eng.vmware.com
* vRO SDK - https://www.vmware.com/support/pubs/orchestrator_pubs.html
* Local vRO Instance - https://10.25.37.28:8281/vco/api/docs/index.html
* Slack API - https://api.slack.com

## Dependencies

* https://github.com/nlopes/slack - Slack API
* https://github.com/jmcvetta/napping - REST API encapsulation
* https://github.com/mattn/go-sqlite3 - sqlite3 database driver

## Setup

### Configure ssh in order to use a private repo (gitlab only)

```
$ cat $HOME/.ssh/config
# qeconfig.wp.fsi
Host qeconfig.wp.fsi
Hostname qeconfig.wp.fsi
User git
```

### Pull the vropsbot code

```
export GOPATH=$HOME/go
go get github.com/bruceadowns/vropsbot
```

Ensure that *$GOPATH/bin* is in your path.

### Customize Config File

`config.json`

## Future Ideas

* vR Ops Suite API integration
* Manage salesforce escalations
