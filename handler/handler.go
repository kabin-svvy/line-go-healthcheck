package handler

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type UploadResponse struct {
	Filename string
	Site     []string
	Total    int
	Success  int
	Fail     int
}

func (h *UploadResponse) transformUploadCsv(data string) {
	h.Site = strings.Split(data, ",")
}

func (h *UploadResponse) removeSite(s []int) {
	for _, v := range s {
		h.Site = append(h.Site[:v], h.Site[v+1:]...)
	}
}

func (h *UploadResponse) sitePing() {
	removeIndex := []int{}
	for k, v := range h.Site {
		v = strings.Replace(strings.Replace(v, "\n", "", -1), "\r", "", -1)
		if len(v) == 0 {
			removeIndex = append(removeIndex, k)
			continue
		}
		h.Site[k] = v
		h.Total = h.Total + 1
		resp, err := http.Get(v)
		if err != nil {
			log.Printf("Test site %v result %v\n", v, err)
			h.Fail = h.Fail + 1
		}
		if resp != nil && resp.StatusCode == http.StatusOK {
			h.Success = h.Success + 1
		}
	}
	if len(removeIndex) > 0 {
		h.removeSite(removeIndex)
	}
}

func Upload(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	uploadRes := &UploadResponse{}

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			return err
		}

		data, _ := file.Open()
		buf := bytes.NewBuffer(nil)
		io.Copy(buf, data)

		uploadRes.transformUploadCsv(buf.String())

		defer src.Close()
	}
	uploadRes.sitePing()
	log.Printf("response %++v", uploadRes)
	return c.JSON(http.StatusOK, uploadRes)
}
