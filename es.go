package main

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strconv"
	"strings"
)

type WithEsClientOption func(client *resty.Client)

type ES struct {
	client  *resty.Client
	ctx     context.Context
	version *uint8
	logger  Logger
}

func (e *ES) Version(v *uint8) (uint8, error) {
	if e.logger != nil {
		e.logger.Trace("es: Version")
	}

	if v != nil {
		if *v < 6 || *v > 8 {
			return 0, fmt.Errorf("es: major supported versions are 6, 7, 8")
		}

		e.version = v

		return *e.version, nil
	}

	var clusterVersion version
	r, err := e.client.R().SetResult(&clusterVersion).Get("/")
	if err != nil {
		return 0, fmt.Errorf("es: version http: %s", err.Error())
	}

	if r.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("es: version http: %s", r.Status())
	}

	n := clusterVersion.Version.Number
	if n == nil {
		return 0, fmt.Errorf("es: version undefined: expected string, got nil")
	}
	pn := *n

	dPosition := strings.Index(pn, ".")
	if dPosition == -1 || dPosition == 0 {
		return 0, fmt.Errorf("es: version parsing: version %s is not valid", pn)
	}

	stringVersion := pn[:strings.IndexByte(pn, '.')]
	if stringVersion == "" {
		return 0, fmt.Errorf("es: version parsing: got empty string")
	}

	intVersion, err := strconv.Atoi(stringVersion)
	if err != nil {
		return 0, fmt.Errorf("es: version parsing: %s", err.Error())
	}

	if intVersion < 6 || intVersion > 8 {
		return 0, fmt.Errorf("es: major supported versions are 6, 7, 8")
	}

	uint8Version := uint8(intVersion)
	e.version = &uint8Version

	return uint8Version, nil
}

func (e *ES) SetLogger(logger Logger) {
	e.logger = logger
}

func (e *ES) SqlQuery(query string) (*SqlResponse, error) {
	if e.logger != nil {
		e.logger.Trace("es: SqlQuery")
	}

	if e.version == nil {
		return nil, fmt.Errorf("es: sql query: version is not defined")
	}

	var path string
	if *e.version == 6 {
		path = "_xpack/sql?format=json"
	} else {
		path = "_sql?format=json"
	}

	q := sqlQuery{
		Query:    query,
		Leniency: true,
	}

	res, err := e.client.R().SetBody(q).Post(path)
	if err != nil {
		return nil, fmt.Errorf("es: sql query: %s", err.Error())
	}

	if res.StatusCode() != http.StatusOK {
		if e.logger != nil {
			e.logger.Infof("es: sql query: invalid response body: %s", string(res.Body()))
		}

		return nil, fmt.Errorf("es: sql query: invalid response %s", res.Status())
	}

	sqlResponse, err := parseJsonResponse(res.Body())
	if err != nil {
		return nil, fmt.Errorf("es: sql query: response parse: %s", err.Error())
	}

	return sqlResponse, nil
}

func Create(ctx context.Context, hc *http.Client, opts ...WithEsClientOption) ES {
	var c *resty.Client

	if hc == nil {
		c = resty.New()
	} else {
		c = resty.NewWithClient(hc)
	}

	es := ES{
		client: c,
		ctx:    ctx,
	}

	for _, v := range opts {
		v(es.client)
	}

	return es
}

func CreateWithBaseUrl(ctx context.Context, baseUrl string, hc *http.Client, opts ...WithEsClientOption) ES {
	opts = append(opts, func(c *resty.Client) {
		c.SetBaseURL(baseUrl)
	})

	return Create(ctx, hc, opts...)
}
