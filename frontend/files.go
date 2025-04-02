package frontend

import (
	"fileupbackendv2/internal/storage"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const PageSize = 50

func getOffsetLimit(pageNo string) (int, int) {
	page, err := strconv.Atoi(pageNo)

	if page < 1 || err != nil {
		page = 1 // Default to first page
	}
	limit := PageSize
	offset := (page - 1) * PageSize
	return offset, limit
}
func listFiles(c echo.Context) error {
	folder := c.QueryParam("folder")
	page := c.QueryParam("page")
	offset, limit := getOffsetLimit(page)
	cookie, err := c.Cookie("api-key")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/app/form/api-key/")
	}
	files, folderSize, err := storage.FindFiles(cookie.Value, folder)
	if err != nil {
		return err
	}
	var totalPages int
	if len(files)%PageSize == 0 {
		totalPages = len(files) / PageSize
	} else {
		totalPages = (len(files) / PageSize) + 1
	}

	if len(files) > (offset + limit) {
		files = files[offset : offset+limit]
	} else {
		files = files
	}
	return renderTemplate(c.Response(), "listfiles", map[string]interface{}{
		"files": files,
		"pages": totalPages,
		"size":  folderSize,
	},
		c,
		false,
	)
}

func uploadFileForm(c echo.Context) error {
	apiKey, err := c.Cookie("api-key")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/app/form/api-key/")
	}
	folders, err := storage.GetAllFolders(apiKey.Value)
	if err != nil {
		return err
	}
	log.Logger.Println(folders)
	return renderTemplate(c.Response(), "uploads", map[string]interface{}{
		"folders": folders,
	},
		c,
		false,
	)
}

func addFolderForm(c echo.Context) error {
	apiKey, err := c.Cookie("api-key")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/app/form/api-key/")
	}
	folders, err := storage.GetAllFolders(apiKey.Value)
	if err != nil {
		return err
	}
	return renderTemplate(c.Response(), "addFolder", map[string]interface{}{
		"folders": folders,
	},
		c,
		false,
	)
}
