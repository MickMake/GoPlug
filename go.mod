module github.com/MickMake/GoPlug

go 1.19

// replace github.com/MickMake/GoUnify => ../../GoUnify

replace github.com/MickMake/GoPlug => ./

require (
	github.com/Masterminds/semver v1.5.0
	github.com/MickMake/GoUnify/Only v0.0.0-20221125023651-ff4a37b1928a
	github.com/MickMake/GoUnify/Unify v0.0.0-20221125023651-ff4a37b1928a
	github.com/MickMake/GoUnify/cmdConfig v0.0.0-20221125023651-ff4a37b1928a
	github.com/MickMake/GoUnify/cmdHelp v0.0.0-20221125023651-ff4a37b1928a
	github.com/briandowns/openweathermap v0.19.0
	github.com/frankban/quicktest v1.14.6
	github.com/h2non/filetype v1.1.3
	github.com/hashicorp/go-hclog v1.5.0
	github.com/hashicorp/go-plugin v1.5.2
	github.com/spf13/cobra v1.7.0
	github.com/spf13/viper v1.16.0
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/MichaelMure/go-term-markdown v0.1.4 // indirect
	github.com/MichaelMure/go-term-text v0.3.1 // indirect
	github.com/MickMake/GoUnify/cmdCron v0.0.0-20221125020154-1b15b8d20735 // indirect
	github.com/MickMake/GoUnify/cmdDaemon v0.0.0-20221125023651-ff4a37b1928a // indirect
	github.com/MickMake/GoUnify/cmdExec v0.0.0-20221125015223-b8c165efd0ec // indirect
	github.com/MickMake/GoUnify/cmdLog v0.0.0-20221125015223-b8c165efd0ec // indirect
	github.com/MickMake/GoUnify/cmdPath v0.0.0-20221125023651-ff4a37b1928a // indirect
	github.com/MickMake/GoUnify/cmdShell v0.0.0-20221125023651-ff4a37b1928a // indirect
	github.com/MickMake/GoUnify/cmdVersion v0.0.0-20221125023651-ff4a37b1928a // indirect
	github.com/abiosoft/ishell v2.0.0+incompatible // indirect
	github.com/abiosoft/ishell/v2 v2.0.2 // indirect
	github.com/abiosoft/readline v0.0.0-20180607040430-155bce2042db // indirect
	github.com/alecthomas/chroma v0.7.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/briandowns/spinner v1.23.0 // indirect
	github.com/danwakefield/fnmatch v0.0.0-20160403171240-cbb64ac3d964 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/dlclark/regexp2 v1.1.6 // indirect
	github.com/eliukblau/pixterm/pkg/ansimage v0.0.0-20191210081756-9fb6cf8c2f75 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/flynn-archive/go-shlex v0.0.0-20150515145356-3f9db97f8568 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-co-op/gocron v1.18.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/gomarkdown/markdown v0.0.0-20191123064959-2c17d62f5098 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-github/v30 v30.1.0 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/yamux v0.0.0-20180604194846-3520598351bb // indirect
	github.com/inconshreveable/go-update v0.0.0-20160112193335-8152e7eb6ccf // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/ivanpirog/coloredcobra v1.0.1 // indirect
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/kyokomi/emoji/v2 v2.2.8 // indirect
	github.com/logrusorgru/aurora v2.0.3+incompatible // indirect
	github.com/lucasb-eyer/go-colorful v1.0.3 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.12 // indirect
	github.com/mitchellh/go-testing-interface v0.0.0-20171004221916-a61a99592b77 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rhysd/go-github-selfupdate v1.2.3 // indirect
	github.com/rivo/uniseg v0.1.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/sevlyar/go-daemon v0.1.6 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/tcnksm/go-gitconfig v0.1.2 // indirect
	github.com/ulikunitz/xz v0.5.9 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/image v0.0.0-20191206065243-da761ea9ff43 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/oauth2 v0.7.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/term v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/grpc v1.55.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
