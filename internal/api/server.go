package api

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/fasthttp/router"
	"github.com/polundrra/shortlink/internal/service"
	"github.com/valyala/fasthttp"
	"net/http"
	"strings"
)

type LinkApi struct {
	service service.Service
}

func New(service service.Service) LinkApi {
	return LinkApi{service}
}


func (l *LinkApi) Router() fasthttp.RequestHandler {
	r := router.New()
	r.POST("/get-short-link", l.getShortLink)
	r.GET("/{code}", l.redirect)
	return r.Handler
}

type request struct {
	url string
	customEnd string
}

type response struct {
	url string
}


func (l *LinkApi) getShortLink(ctx *fasthttp.RequestCtx) {
	req := request{}
	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Write([]byte("error unmarshal request body" + err.Error()))
		return
	}

	url := strings.TrimSpace(req.url)
	if !govalidator.IsURL(url) {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		ctx.Write([]byte("invalid URL"))
		return
	}

	customEnd := strings.TrimSpace(req.customEnd)
	shortLink, err := l.service.GetShortLink(ctx, url, customEnd)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Write([]byte(err.Error()))
		return
	}

	resp, err := json.Marshal(response{shortLink})
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Write([]byte(err.Error()))
		return
	}

	if _, err := ctx.Write(resp); err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Write([]byte(err.Error()))
		return
	}
}

func (l *LinkApi) redirect(ctx *fasthttp.RequestCtx) {
	shortLink := (ctx.UserValue("code")).(string)
	url, err := l.service.GetLongLink(ctx, shortLink)
	if err != nil {
		ctx.Response.SetStatusCode(http.StatusInternalServerError)
		ctx.Write([]byte(err.Error()))
		return
	}
	ctx.Response.SetStatusCode(http.StatusPermanentRedirect)
	ctx.Response.Header.Set("Location", url)
}
