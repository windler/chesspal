module github.com/windler/chesspal

go 1.13

require (
	github.com/gorilla/websocket v1.5.0
	github.com/gosuri/uilive v0.0.4
	github.com/jacobsa/go-serial v0.0.0-20180131005756-15cf729a72d4
	github.com/labstack/echo/v4 v4.7.2
	github.com/notnil/chess v1.8.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/notnil/chess => github.com/windler/chess v1.8.0
