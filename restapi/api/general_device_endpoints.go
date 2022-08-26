package api

import (
	"github.com/danielpaulus/go-ios/ios"
	"github.com/danielpaulus/go-ios/ios/instruments"
	"github.com/danielpaulus/go-ios/ios/screenshotr"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GetBookByISBN locates the book whose ISBN value matches the isbn
// GetBookByISBN                godoc
// @Summary      Get single book by isbn
// @Description  Returns the book whose ISBN value matches the isbn.
// @Tags         books
// @Produce      json
// @Param        isbn  path      string  true  "search book by isbn"
// @Success      200  {object}  map[string]interface{}
// @Router       /books/{isbn} [get]
func Info(c *gin.Context) {
	uid := c.Param("uid")
	device, _ := ios.GetDevice(uid)

	allValues, err := ios.GetValuesPlist(device)
	if err != nil {
		print(err)
	}
	svc, err := instruments.NewDeviceInfoService(device)
	if err != nil {
		log.Debugf("could not open instruments, probably dev image not mounted %v", err)
	}
	if err == nil {
		info, err := svc.NetworkInformation()
		if err != nil {
			log.Debugf("error getting networkinfo from instruments %v", err)
		} else {
			allValues["instruments:networkInformation"] = info
		}
		info, err = svc.HardwareInformation()
		if err != nil {
			log.Debugf("error getting hardwareinfo from instruments %v", err)
		} else {
			allValues["instruments:hardwareInformation"] = info
		}
	}
	c.IndentedJSON(http.StatusOK, allValues)
}

func Screenshot(c *gin.Context) {
	uid := c.Param("uid")
	device, _ := ios.GetDevice(uid)

	conn, err := screenshotr.New(device)
	log.Error(err)
	b, _ := conn.TakeScreenshot()

	c.Header("Content-Type", "image/png")
	c.Data(http.StatusOK, "application/octet-stream", b)
}