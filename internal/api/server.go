package api

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"github.com/fasthttp/router"
	"github.com/polundrra/shortlink/internal/service"
	"github.com/valyala/fasthttp"
	"log"
	"regexp"
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
	Url string
	CustomEnd string
}

type response struct {
	Url string
}

func (l *LinkApi) getShortLink(ctx *fasthttp.RequestCtx) {
	req := request{}
	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("error unmarshal request body" + err.Error()))
		return
	}

	if req.CustomEnd != "" {
		if !validateEnding(req.CustomEnd) {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.WriteString("invalid custom link")
			return
		}
	}

	url := strings.TrimSpace(req.Url)
	if !govalidator.IsURL(url) {
		log.Println(url)
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("invalid URL"))
		return
	}

	customEnd := strings.TrimSpace(req.CustomEnd)
	shortLink, err := l.service.CreateShortLink(ctx, url, customEnd)
	if err != nil {
		if err == service.ErrCodeConflict {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
		} else {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		}
		ctx.Write([]byte(err.Error()))
		return
	}

	resp, err := json.Marshal(&response{shortLink})
	if err != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(err.Error()))
		return
	}

	if _, err := ctx.Write(resp); err != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(err.Error()))
		return
	}
}

func (l *LinkApi) redirect(ctx *fasthttp.RequestCtx) {
	shortLink := (ctx.UserValue("code")).(string)
	url, err := l.service.GetLongLink(ctx, shortLink)
	if err != nil {
		if err == service.ErrLongLinkNotFound {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(err.Error()))
		return
	}
	ctx.Response.SetStatusCode(fasthttp.StatusPermanentRedirect)
	ctx.Response.Header.Set("Location", url)
}

func validateEnding(ending string) bool {
	regex:= regexp.MustCompile("^[a-zA-Z0-9-_]{1,32}$")

	return regex.MatchString(ending)
}