package app

import (
	"github.com/gin-gonic/gin"
	"log"
)


func (app *Application) healthcheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (app *Application) GetCategoriesHandler(c *gin.Context) {
	categories := app.service.GetCategoriesService()
	c.JSON(200, gin.H{"data":categories})
}

func (app *Application) GetCategoryTreeHandler(c *gin.Context) {
	categoryID, ok := c.GetQuery("id")
	if !ok {
		c.JSON(400, gin.H{"error": "invalid query parameter"})
		return
	}
	categoryTree := app.service.GetCategoryTreeService(categoryID)
	c.JSON(200, gin.H{"data":categoryTree})
}

func (app *Application) GetContentHandler(c *gin.Context) {
	ContentId, ok := c.GetQuery("id")
	if !ok {
		c.JSON(400, gin.H{"error": "invalid query parameter"})
		return
	}
	content := app.service.GetContentService(ContentId)
	c.JSON(200, gin.H{"data":content})
}


func (app *Application) webhookHandler(c *gin.Context) {
	// 处理 Webhook 请求的逻辑
	token := c.GetHeader("X-Codeup-Token")
	if token != app.WebhookSecret {
		// 403 Forbidden: 密钥不对
		c.JSON(403, gin.H{"error": "Unauthorized"})
		return
	}
	go func() {
		select {
			case app.WebhookChan <- struct{}{}:
				log.Println("[Webhook] Sync signal sent successfully.")
				return
			default:
				log.Println("[Webhook] WARNING: Worker is busy, dropping sync request.")
				return
		}

	}()
	c.JSON(200, gin.H{"status": "Webhook received"})
}