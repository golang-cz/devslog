package devslog

import (
	"bytes"
	"context"
	"encoding"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type developHandler struct {
	opts Options
	goas []groupOrAttrs
	mu   *sync.Mutex
	out  io.Writer
}

type Options struct {
	// You can use standard slog.HandlerOptions, that would be used in production
	*slog.HandlerOptions

	// Max number of printed elements in slice.
	MaxSlicePrintSize uint

	// If the attributes should be sorted by keys
	SortKeys bool

	// Time format for timestamp, default format is "[15:04:05]"
	TimeFormat string
}

type groupOrAttrs struct {
	group string
	attrs []slog.Attr
}

func NewHandler(out io.Writer, o *Options) *developHandler {
	h := &developHandler{out: out, mu: &sync.Mutex{}}
	if o != nil {
		h.opts = *o

		if o.HandlerOptions != nil {
			h.opts.HandlerOptions = o.HandlerOptions
			if o.Level == nil {
				h.opts.Level = slog.LevelInfo
			} else {
				h.opts.HandlerOptions.Level = o.HandlerOptions.Level
			}
		} else {
			h.opts.HandlerOptions = &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}
		}

		if o.MaxSlicePrintSize == 0 {
			h.opts.MaxSlicePrintSize = 50
		}

		if o.TimeFormat == "" {
			h.opts.TimeFormat = "[15:04:05]"
		}
	} else {
		h.opts = Options{
			HandlerOptions:    &slog.HandlerOptions{Level: slog.LevelInfo},
			MaxSlicePrintSize: 50,
			TimeFormat:        "[15:04:05]",
		}
	}

	return h
}

func (h *developHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return l >= h.opts.Level.Level()
}

func (h *developHandler) WithGroup(s string) slog.Handler {
	if s == "" {
		return h
	}

	return h.withGroupOrAttrs(groupOrAttrs{group: s})
}

func (h *developHandler) WithAttrs(as []slog.Attr) slog.Handler {
	if len(as) == 0 {
		return h
	}

	return h.withGroupOrAttrs(groupOrAttrs{attrs: as})
}

func (h *developHandler) withGroupOrAttrs(goa groupOrAttrs) *developHandler {
	h2 := *h
	h2.goas = make([]groupOrAttrs, len(h.goas)+1)
	copy(h2.goas, h.goas)
	h2.goas[len(h2.goas)-1] = goa
	return &h2
}

func (h *developHandler) Handle(ctx context.Context, r slog.Record) error {
	b := make([]byte, 0, 1024)
	b = append(b, csf([]byte(r.Time.Format(h.opts.TimeFormat)), fgWhite)...)
	b = append(b, ' ')
	b = h.formatSourceInfo(b, &r)
	b = h.levelMessage(b, &r)
	b = h.processAttributes(b, &r)

	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := h.out.Write(b)

	return err
}

func (h *developHandler) formatSourceInfo(b []byte, r *slog.Record) []byte {
	if h.opts.AddSource {
		f, _ := runtime.CallersFrames([]uintptr{r.PC}).Next()
		b = append(b, cs([]byte("@@@"), fgBlue)...)
		b = append(b, ' ')
		b = append(b, ul(cs([]byte(f.File), fgYellow))...)
		b = append(b, ':')
		b = append(b, cs([]byte(strconv.Itoa(f.Line)), fgRed)...)
		b = append(b, '\n')
	}

	return b
}

func (h *developHandler) levelMessage(b []byte, r *slog.Record) []byte {
	var bgColor backgroundColor
	var fgColor foregroundColor
	var ls string
	if h.opts.ReplaceAttr != nil {
		a := h.opts.ReplaceAttr(nil, slog.Any(slog.LevelKey, r.Level))
		ls = a.Value.String()
		if a.Key != "level" {
			r.AddAttrs(a)
		}
	} else {
		ls = r.Level.String()
	}

	lr := r.Level
	switch {
	case lr < 0:
		bgColor, fgColor = bgBlue, fgBlue
	case lr < 4:
		bgColor, fgColor = bgGreen, fgGreen
	case lr < 8:
		bgColor, fgColor = bgYellow, fgYellow
	default:
		bgColor, fgColor = bgRed, fgRed
	}

	b = append(b, csb([]byte(" "+ls+" "), fgBlack, bgColor)...)
	b = append(b, ' ')
	b = append(b, cs([]byte(r.Message), fgColor)...)
	b = append(b, '\n')

	return b
}

func (h *developHandler) processAttributes(b []byte, r *slog.Record) []byte {
	var as attributes
	if r.NumAttrs() != 0 {
		r.Attrs(func(a slog.Attr) bool {
			as = append(as, a)
			return true
		})
	}

	goas := h.goas
	if r.NumAttrs() == 0 {
		for len(goas) > 0 && goas[len(goas)-1].group != "" {
			goas = goas[:len(goas)-1]
		}
	}

	for i := len(goas) - 1; i >= 0; i-- {
		if goas[i].group != "" {
			ng := slog.Attr{
				Key:   goas[i].group,
				Value: slog.GroupValue(as...),
			}
			as = attributes{ng}
		} else {
			as = append(as, goas[i].attrs...)
		}
	}

	b = h.colorize(b, as, 0, []string{})
	b = append(b, '\n')
	return b
}

func (h *developHandler) colorize(b []byte, as attributes, l int, g []string) []byte {
	if h.opts.SortKeys {
		sort.Sort(as)
	}

	p := as.padding(fgMagenta)
	for _, a := range as {
		if h.opts.ReplaceAttr != nil {
			a = h.opts.ReplaceAttr(g, a)
		}

		k := cs([]byte(a.Key), fgMagenta)
		v := []byte(a.Value.String())
		m := []byte(" ")

		switch a.Value.Kind() {
		case slog.KindFloat64, slog.KindInt64, slog.KindUint64:
			m = cs([]byte("#"), fgYellow)
			v = cs(v, fgYellow)
		case slog.KindBool:
			m = cs([]byte("#"), fgRed)
			v = cs(v, fgRed)
		case slog.KindString:
			if len(v) == 0 {
				v = csf([]byte("empty"), fgWhite)
			} else if h.isURL(v) {
				m = cs([]byte("*"), fgBlue)
				v = ul(cs(v, fgBlue))
			}
		case slog.KindTime, slog.KindDuration:
			m = cs([]byte("@"), fgCyan)
			v = cs(v, fgCyan)
		case slog.KindAny:
			any := a.Value.Any()

			if err, ok := any.(error); ok {
				m = cs([]byte("E"), fgRed)
				v = h.formatError(err, l)
				break
			}

			if t, ok := any.(*time.Time); ok {
				m = cs([]byte("@"), fgCyan)
				v = cs([]byte(t.String()), fgCyan)
				break
			}

			if dur, ok := any.(*time.Duration); ok {
				m = cs([]byte("@"), fgCyan)
				v = cs([]byte(dur.String()), fgCyan)
				break
			}

			if textMarshaller, ok := any.(encoding.TextMarshaler); ok {
				v = atb(textMarshaller)
				break
			}

			at := reflect.TypeOf(any)
			av := reflect.ValueOf(any)
			if at == nil {
				m = cs([]byte("!"), fgRed)
				v = cs(atb("Unknown type"), fgRed)
				break
			}

			ut, uv := h.reducePointerTypeValue(at, av)

			switch ut.Kind() {
			case reflect.Array:
				m = cs([]byte("A"), fgGreen)
				v = h.formatSlice(at, av, l)
			case reflect.Slice:
				m = cs([]byte("S"), fgGreen)
				v = h.formatSlice(at, av, l)
			case reflect.Map:
				m = cs([]byte("M"), fgGreen)
				v = h.formatMap(at, av, l)
			case reflect.Struct:
				m = cs([]byte("S"), fgYellow)
				v = h.formatStruct(at, av, 0)
			case reflect.Float32, reflect.Float64:
				m = cs([]byte("#"), fgYellow)
				v = cs(atb(uv.Float()), fgYellow)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				m = cs([]byte("#"), fgYellow)
				v = cs(atb(uv.Int()), fgYellow)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				m = cs([]byte("#"), fgYellow)
				v = cs(atb(uv.Uint()), fgYellow)
			case reflect.Bool:
				m = cs([]byte("#"), fgRed)
				v = cs(atb(uv.Bool()), fgRed)
			case reflect.String:
				s := uv.String()
				if len(s) == 0 {
					v = csf([]byte("empty"), fgWhite)
				} else if h.isURL([]byte(s)) {
					v = ul(cs(v, fgBlue))
				} else {
					v = []byte(uv.String())
				}

			default:
				m = cs([]byte("!"), fgRed)
				v = cs(atb("Unknown type"), fgRed)
			}
		case slog.KindGroup:
			m = cs([]byte("G"), fgGreen)
			var ga attributes
			ga = a.Value.Group()
			g = append(g, a.Key)

			v = cs([]byte("============"), fgGreen)
			v = append(v, '\n')
			v = append(v, h.colorize(nil, ga, l+1, g)...)
		}

		b = append(b, bytes.Repeat([]byte(" "), l*2)...)
		b = append(b, m...)
		b = append(b, ' ')
		b = append(b, k...)
		b = append(b, bytes.Repeat([]byte(" "), p-len(k))...)
		b = append(b, ':')
		b = append(b, ' ')
		b = append(b, v...)
		if a.Value.Kind() != slog.KindGroup {
			b = append(b, '\n')
		}
	}

	return b
}

func (h *developHandler) isURL(u []byte) bool {
	_, err := url.ParseRequestURI(string(u))
	return err == nil
}

func (h *developHandler) formatError(err error, l int) (b []byte) {
	if err == nil {
		b = append(b, ul(cs([]byte("<nil>"), fgRed))...)
		return
	}

	for i := 0; err != nil; i++ {
		b = append(b, '\n')
		b = append(b, bytes.Repeat([]byte(" "), l*2+4)...)
		tb := strconv.Itoa(i)
		b = append(b, cs([]byte(tb), fgRed)...)
		b = append(b, cs([]byte(": "), fgWhite)...)

		errMsg := err.Error()
		ue := errors.Unwrap(err)
		if ue != nil {
			errMsg, _ = strings.CutSuffix(errMsg, ue.Error())
			errMsg, _ = strings.CutSuffix(errMsg, ": ")
		}
		if errMsg == "" {
			errMsg = fmt.Sprintf("[%T]", err)
		}
		b = append(b, cs([]byte(errMsg), fgRed)...)

		for j, fileLine := range getFileLineFromPC(extractPCFromError(err)) {
			b = append(b, '\n')
			tb := strconv.Itoa(j)
			b = append(b, bytes.Repeat([]byte(" "), l*2+6)...)
			b = append(b, bytes.Repeat([]byte(" "), len(tb))...)
			b = append(b, cs([]byte(tb), fgBlue)...)
			b = append(b, []byte(": ")...)
			b = append(b, ul(cs([]byte(fileLine), fgBlue))...)
		}

		err = ue
	}

	return b
}

func (h *developHandler) formatSlice(st reflect.Type, sv reflect.Value, l int) (b []byte) {
	ts := h.buildTypeString(st.String())
	st, sv = h.reducePointerTypeValue(st, sv)

	b = append(b, cs([]byte(strconv.Itoa(sv.Len())), fgBlue)...)
	b = append(b, ' ')
	b = append(b, ts...)
	d := len(strconv.Itoa(sv.Len()))
	if len(strconv.Itoa(int(h.opts.MaxSlicePrintSize))) < d {
		d = len(strconv.Itoa(int(h.opts.MaxSlicePrintSize)))
	}

	for i := 0; i < sv.Len(); i++ {
		if i == int(h.opts.MaxSlicePrintSize) {
			b = append(b, '\n')
			b = append(b, bytes.Repeat([]byte(" "), l*2+4)...)
			b = append(b, bytes.Repeat([]byte(" "), d+2)...)
			b = append(b, cs([]byte("..."), fgBlue)...)
			b = append(b, cs([]byte("]"), fgGreen)...)
			break
		}

		v := sv.Index(i)
		t := v.Type()

		tb := strconv.Itoa(i)
		b = append(b, '\n')
		b = append(b, bytes.Repeat([]byte(" "), l*2+4)...)
		b = append(b, bytes.Repeat([]byte(" "), d-len(tb))...)
		b = append(b, cs([]byte(tb), fgGreen)...)
		b = append(b, ':')
		b = append(b, ' ')
		b = append(b, h.elementType(t, v, l)...)

	}

	return b
}

func (h *developHandler) formatMap(st reflect.Type, sv reflect.Value, l int) (b []byte) {
	ts := h.buildTypeString(st.String())
	st, sv = h.reducePointerTypeValue(st, sv)

	p := h.mapKeyPadding(sv, fgGreen)
	b = append(b, cs([]byte(strconv.Itoa(sv.Len())), fgBlue)...)
	b = append(b, ' ')
	b = append(b, ts...)
	sk := h.sortMapKeys(sv)
	for _, k := range sk {
		v := sv.MapIndex(k)
		v = h.reducePointerValue(v)
		k = h.reducePointerValue(k)

		tb := cs(atb(k.Interface()), fgGreen)
		b = append(b, '\n')
		b = append(b, bytes.Repeat([]byte(" "), l*2+4)...)
		b = append(b, tb...)
		b = append(b, bytes.Repeat([]byte(" "), p-len(tb))...)
		b = append(b, ':')
		b = append(b, ' ')
		b = append(b, h.elementType(v.Type(), v, l)...)
	}

	return b
}

func (h *developHandler) formatStruct(st reflect.Type, sv reflect.Value, l int) (b []byte) {
	b = h.buildTypeString(st.String())

	st, sv = h.reducePointerTypeValue(st, sv)
	p := h.structKeyPadding(sv, fgGreen)

	for i := 0; i < sv.NumField(); i++ {
		v := sv.Field(i)
		t := v.Type()

		tb := cs([]byte(sv.Type().Field(i).Name), fgGreen)
		b = append(b, '\n')
		b = append(b, bytes.Repeat([]byte(" "), l*2+4)...)
		b = append(b, tb...)
		b = append(b, bytes.Repeat([]byte(" "), p-len(tb))...)
		b = append(b, ':')
		b = append(b, ' ')
		b = append(b, h.elementType(t, v, l)...)
	}

	return b
}

var marshalTextInterface = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()

func (h *developHandler) elementType(t reflect.Type, v reflect.Value, l int) (b []byte) {
	if t.Implements(marshalTextInterface) {
		return atb(v)
	}

	switch v.Kind() {
	case reflect.Array:
		b = h.formatSlice(t, v, l+1)
	case reflect.Slice:
		b = h.formatSlice(t, v, l+1)
	case reflect.Map:
		b = h.formatMap(t, v, l+1)
	case reflect.Struct:
		b = h.formatStruct(t, v, l+1)
	case reflect.Pointer:
		if v.IsNil() {
			b = cs([]byte("<"), fgRed)
			b = append(b, cs([]byte("nil"), fgYellow)...)
			b = append(b, cs([]byte(">"), fgRed)...)
		} else {
			b = h.elementType(t, v.Elem(), l)
		}
	case reflect.Float32, reflect.Float64:
		b = cs(atb(v.Float()), fgYellow)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		b = cs(atb(v.Int()), fgYellow)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		b = cs(atb(v.Uint()), fgYellow)
	case reflect.Bool:
		b = cs(atb(v.Bool()), fgRed)
	case reflect.String:
		s := v.String()
		if len(s) == 0 {
			b = csf([]byte("empty"), fgWhite)
		} else if h.isURL([]byte(s)) {
			b = ul(cs([]byte(s), fgBlue))
		} else {
			b = atb(s)
		}
	default:
		b = atb("Unknown type: ")
		b = append(b, atb(v.Kind())...)
		b = cs(b, fgRed)
	}

	return b
}

func (h *developHandler) buildTypeString(ts string) (b []byte) {
	t := []byte(ts)

	for len(t) > 0 {
		switch t[0] {
		case '*':
			b = append(b, cs([]byte{t[0]}, fgRed)...)
		case '[', ']':
			b = append(b, cs([]byte{t[0]}, fgGreen)...)
		default:
			b = append(b, cs([]byte{t[0]}, fgYellow)...)
		}

		t = t[1:]
	}

	return b
}

func (h *developHandler) sortMapKeys(rv reflect.Value) []reflect.Value {
	ks := make([]reflect.Value, 0, rv.Len())
	for _, k := range rv.MapKeys() {
		ks = append(ks, k)
	}

	sort.Slice(ks, func(i, j int) bool {
		return fmt.Sprint(ks[i].Interface()) < fmt.Sprint(ks[j].Interface())
	})

	return ks
}

func (h *developHandler) mapKeyPadding(rv reflect.Value, fgColor foregroundColor) (p int) {
	for _, k := range rv.MapKeys() {
		k = h.reducePointerValue(k)
		c := len(cs(atb(k.Interface()), fgColor))
		if c > p {
			p = c
		}
	}

	return p
}

func (h *developHandler) structKeyPadding(sv reflect.Value, fgColor foregroundColor) (p int) {
	st := sv.Type()
	for i := 0; i < sv.NumField(); i++ {
		c := len(cs([]byte(st.Field(i).Name), fgColor))
		if c > p {
			p = c
		}
	}

	return p
}

func (h *developHandler) reducePointerValue(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	return v
}

func (h *developHandler) reducePointerTypeValue(t reflect.Type, v reflect.Value) (reflect.Type, reflect.Value) {
	if t == nil {
		return t, v
	}

	for v.Kind() == reflect.Pointer {
		v = v.Elem()
		if isNilValue(v) {
			return t, v
		}
		t = v.Type()
	}

	return t, v
}

// Any to []byte using fmt.Sprintf
func atb(a any) []byte {
	return []byte(fmt.Sprintf("%v", a))
}

func isNilValue(v reflect.Value) bool {
	nilValue := reflect.ValueOf(nil)
	return v == nilValue
}
