package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/jbenzshawel/go-sandbox/common/auth"
	crest "github.com/jbenzshawel/go-sandbox/common/rest"
	"github.com/jbenzshawel/go-sandbox/identity/app/query"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user/permission"
)

func (h *HttpHandler) GetUsers(ctx *gin.Context) {
	h.authorize(ctx, permission.ViewUsers, h.getUsers)
}

func (h *HttpHandler) getUsers(ctx *gin.Context, authUser *auth.User) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	pageSize, _ := strconv.Atoi(ctx.Query("pageSize"))
	if pageSize == 0 {
		pageSize = 20
	}

	users, err := h.app.Queries.Users.Handle(ctx, query.Users{Page: page, PageSize: pageSize})
	if err != nil {
		crest.HandleErrorResponse(ctx, err)
		return
	}

	resp := make([]*userResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, mapUserResponse(u))
	}
	ctx.JSON(http.StatusOK, resp)
}
