package hostname

import (
	"errors"
	"net"
	"regexp"
	"strings"

	"github.com/coredns/coredns/middleware"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

// Hostname is a middleware that serve ec2 style hostname
type Hostname struct{}

var re = regexp.MustCompile("^ip-(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])-){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$")

// ServeDNS implements the middleware.Handler interface.
func (wh Hostname) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	opt := middleware.Options{}
	state := request.Request{W: w, Req: r}

	parts := dns.SplitDomainName(state.Name())
	if len(parts) < 2 {
		return dns.RcodeServerFailure, middleware.Error(wh.Name(), errors.New("invalid query"))
	}
	query := parts[0]
	zone := dns.Name(strings.Join(parts[1:], ".")).String() + "."

	switch state.Type() {
	case "A":
		if len(re.FindAllString(query, -1)) > 0 {
			a := new(dns.Msg)
			a.SetReply(r)
			a.Authoritative, a.RecursionAvailable, a.Compress = true, true, true
			rr := new(dns.A)
			rr.Hdr = dns.RR_Header{
				Name:   state.QName(),
				Rrtype: dns.TypeA,
				Class:  state.QClass(),
				Ttl:    3600,
			}
			rr.A = net.ParseIP(strings.Replace(query[3:], "-", ".", -1)).To4()
			a.Answer = append(a.Answer, rr)

			state.SizeAndDo(a)
			w.WriteMsg(a)
		} else {
			return middleware.BackendError(nil, zone, dns.RcodeNameError, state, nil, nil, opt)
		}

	default:
		return middleware.BackendError(nil, zone, dns.RcodeNameError, state, nil, nil, opt)
	}

	return 0, nil
}

// Name implements the Handler interface.
func (wh Hostname) Name() string { return "hostname" }
