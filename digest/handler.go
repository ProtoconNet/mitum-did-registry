package digest

import (
	"context"
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-currency/v3/types"
	isaacnetwork "github.com/ProtoconNet/mitum2/isaac/network"
	"github.com/ProtoconNet/mitum2/network/quicmemberlist"
	"github.com/ProtoconNet/mitum2/network/quicstream"
	"net/http"
	"time"

	"github.com/ProtoconNet/mitum-currency/v3/digest/network"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/launch"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/logging"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/singleflight"
)

var (
	HandlerPathDIDDesign   = `/did-registry/{contract:(?i)` + types.REStringAddressString + `}`
	HandlerPathDIDData     = `/did-registry/{contract:(?i)` + types.REStringAddressString + `}/did/{pubKey:` + types.ReSpecialCh + `}`
	HandlerPathDIDDocument = `/did-registry/{contract:(?i)` + types.REStringAddressString + `}/document`
)

func init() {
	if b, err := currencydigest.JSON.Marshal(currencydigest.UnknownProblem); err != nil {
		panic(err)
	} else {
		currencydigest.UnknownProblemJSON = b
	}
}

type Handlers struct {
	*zerolog.Logger
	networkID       base.NetworkID
	encoders        *encoder.Encoders
	encoder         encoder.Encoder
	database        *currencydigest.Database
	cache           currencydigest.Cache
	nodeInfoHandler currencydigest.NodeInfoHandler
	send            func(interface{}) (base.Operation, error)
	client          func() (*isaacnetwork.BaseClient, *quicmemberlist.Memberlist, []quicstream.ConnInfo, error)
	router          *mux.Router
	routes          map[string]*mux.Route
	itemsLimiter    func(string) int64
	rg              *singleflight.Group
	expireNotFilled time.Duration
}

func NewHandlers(
	ctx context.Context,
	networkID base.NetworkID,
	encs *encoder.Encoders,
	enc encoder.Encoder,
	st *currencydigest.Database,
	cache currencydigest.Cache,
	router *mux.Router,
	routes map[string]*mux.Route,
) *Handlers {
	var log *logging.Logging
	if err := util.LoadFromContextOK(ctx, launch.LoggingContextKey, &log); err != nil {
		return nil
	}

	return &Handlers{
		Logger:          log.Log(),
		networkID:       networkID,
		encoders:        encs,
		encoder:         enc,
		database:        st,
		cache:           cache,
		router:          router,
		routes:          routes,
		itemsLimiter:    currencydigest.DefaultItemsLimiter,
		rg:              &singleflight.Group{},
		expireNotFilled: time.Second * 3,
	}
}

func (hd *Handlers) Initialize() error {
	//cors := handlers.CORS(
	//	handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
	//	handlers.AllowedHeaders([]string{"content-type"}),
	//	handlers.AllowedOrigins([]string{"*"}),
	//	handlers.AllowCredentials(),
	//)
	//hd.router.Use(cors)

	hd.setHandlers()

	return nil
}

func (hd *Handlers) SetLimiter(f func(string) int64) *Handlers {
	hd.itemsLimiter = f

	return hd
}

func (hd *Handlers) Cache() currencydigest.Cache {
	return hd.cache
}

func (hd *Handlers) Router() *mux.Router {
	return hd.router
}

func (hd *Handlers) Routes() map[string]*mux.Route {
	return hd.routes
}

func (hd *Handlers) Handler() http.Handler {
	return network.HTTPLogHandler(hd.router, hd.Logger)
}

func (hd *Handlers) setHandlers() {
	get := 1000
	_ = hd.setHandler(HandlerPathDIDData, hd.handleDIDData, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDIDDesign, hd.handleDIDDesign, true, get, get).
		Methods(http.MethodOptions, "GET")
	_ = hd.setHandler(HandlerPathDIDDocument, hd.handleDIDDocument, true, get, get).
		Methods(http.MethodOptions, "GET")
}

func (hd *Handlers) setHandler(prefix string, h network.HTTPHandlerFunc, useCache bool, rps, burst int) *mux.Route {
	var handler http.Handler
	if !useCache {
		handler = http.HandlerFunc(h)
	} else {
		ch := currencydigest.NewCachedHTTPHandler(hd.cache, h)

		handler = ch
	}

	var name string
	if prefix == "" || prefix == "/" {
		name = "root"
	} else {
		name = prefix
	}

	var route *mux.Route
	if r := hd.router.Get(name); r != nil {
		route = r
	} else {
		route = hd.router.Name(name)
	}

	handler = currencydigest.RateLimiter(rps, burst)(handler)

	/*
		if rules, found := hd.rateLimit[prefix]; found {
			handler = process.NewRateLimitMiddleware(
				process.NewRateLimit(rules, limiter.Rate{Limit: -1}), // NOTE by default, unlimited
				hd.rateLimitStore,
			).Middleware(handler)

			hd.Log().Debug().Str("prefix", prefix).Msg("ratelimit middleware attached")
		}
	*/

	route = route.
		Path(prefix).
		Handler(handler)

	hd.routes[prefix] = route

	return route
}

func (hd *Handlers) combineURL(path string, pairs ...string) (string, error) {
	if n := len(pairs); n%2 != 0 {
		return "", errors.Errorf("failed to combine url; uneven pairs to combine url")
	} else if n < 1 {
		u, err := hd.routes[path].URL()
		if err != nil {
			return "", errors.Wrap(err, "failed to combine url")
		}
		return u.String(), nil
	}

	u, err := hd.routes[path].URLPath(pairs...)
	if err != nil {
		return "", errors.Wrap(err, "failed to combine url")
	}

	return u.String(), nil
}

//func (hd *Handlers) combineQuery(path string, pairs ...string) (string, error) {
//	if n := len(pairs); n%2 != 0 {
//		return "", errors.Errorf("failed to combine url; uneven pairs to combine url")
//	} else if n < 1 {
//		u, err := hd.routes[path].URL()
//		if err != nil {
//			return "", errors.Wrap(err, "failed to combine url")
//		}
//		return u.String(), nil
//	}
//
//	u, err := hd.routes[path].URL()
//	if err != nil {
//		return "", errors.Wrap(err, "failed to combine query")
//	}
//
//	query := u.Query()
//
//	// Add query parameters from pairs
//	for i := 0; i < len(pairs); i += 2 {
//		query.Add(pairs[i], pairs[i+1])
//	}
//
//	// Encode the query back into the URL
//	u.RawQuery = query.Encode()
//
//	return u.String(), nil
//}
