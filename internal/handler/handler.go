package handler

import (
	"crud/internal/domain"
	"crud/internal/pkg/authclient"
	"crud/internal/service"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

func ServerHandler(ctx *fasthttp.RequestCtx) {
	log.Debug().
		Str("method", string(ctx.Method())).
		Str("uri", ctx.Request.URI().String()).
		Str("from", string(ctx.RemoteAddr().String())).
		Msg("request")

	ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowOrigin, "*")
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowMethods, fasthttp.MethodPost)
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowMethods, fasthttp.MethodGet)
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowMethods, fasthttp.MethodDelete)
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowHeaders, fasthttp.HeaderContentType)
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowHeaders, fasthttp.HeaderAuthorization)

	if ctx.IsOptions() {
		return
	}

	switch {
	case ctx.IsGet():
		GetHandler(ctx)
	case ctx.IsPost():
		PostHandler(ctx)
	case ctx.IsPut():
		PutHandler(ctx)
	case ctx.IsDelete():
		DeleteHandler(ctx)
	default:
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	}
}

func GetHandler(ctx *fasthttp.RequestCtx) {
	id := ctx.QueryArgs().Peek("id")
	if len(id) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	rec, err := service.Get(string(id))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	if err := json.NewEncoder(ctx).Encode(rec); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func DeleteHandler(ctx *fasthttp.RequestCtx) {
	id := ctx.QueryArgs().Peek("id")
	if len(id) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	userID, err := service.GetAuthorID(string(id))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	if !isAuth(ctx, userID) {
		return
	}

	if err := service.Delete(string(id)); err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func PostHandler(ctx *fasthttp.RequestCtx) {
	var rec domain.Recipe
	if err := json.Unmarshal(ctx.PostBody(), &rec); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if !isAuth(ctx, rec.AuthorID) {
		return
	}

	rec.ID = ""
	if err := service.AddOrUpd(&rec); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	resp := IdResponse{ID: rec.ID}
	if err := json.NewEncoder(ctx).Encode(resp); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func PutHandler(ctx *fasthttp.RequestCtx) {
	id := ctx.QueryArgs().Peek("id")
	if len(id) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	userID, err := service.GetAuthorID(string(id))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	var rec domain.Recipe
	if err := json.Unmarshal(ctx.PostBody(), &rec); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if rec.AuthorID != userID || rec.ID != string(id) {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if !isAuth(ctx, userID) {
		return
	}

	if err := service.AddOrUpd(&rec); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	resp := IdResponse{ID: rec.ID}
	if err := json.NewEncoder(ctx).Encode(resp); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func isAuth(ctx *fasthttp.RequestCtx, userID string) bool {
	token := ctx.Request.Header.Peek(fasthttp.HeaderAuthorization)
	validateResult := authclient.ValidateToken(string(token), userID)
	if string(token) == "" || !validateResult {
		log.Trace().
			Any("validating result", validateResult).
			Str("token", string(token)).
			Msg("not authorized request")
		ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		return false
	}
	return true
}
