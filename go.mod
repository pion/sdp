module github.com/pion/sdp/v2

go 1.12

require (
	github.com/pion/randutil v0.1.0
	github.com/stretchr/testify v1.6.1
        github.com/pion/ice/v2 v2.2.24
)

replace (
	github.com/pion/ice/v2 => /home/kory/Desktop/Code/RTC/pion_ice/ice
)
