package topic

import (
	topicHandler "e-learn/internal/handlers/topic"
	"github.com/gin-gonic/gin"
)

func TopicRoutes(r *gin.Engine) {

	r.GET("/topics/pref/:slug", topicHandler.StoreTopicPreference)

}
