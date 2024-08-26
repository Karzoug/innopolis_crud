package authclient

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

var c *fasthttp.HostClient

func Init(host string) {
	c = &fasthttp.HostClient{
		Addr: host,
	}
}

func ValidateToken(token, userID string) bool {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI("http://" + c.Addr + "/v2/get_user_info")
	req.URI().QueryArgs().Add("user_id", userID)

	req.Header.Set(fasthttp.HeaderAuthorization, token)
	req.Header.SetHost(c.Addr)
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := c.Do(req, resp); err != nil {
		return false
	}
	log.Trace().Any("auth client response", resp)

	return resp.StatusCode() == http.StatusOK
}
