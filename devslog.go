package devslog

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
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

	// Time format for timestamp, default format is "[15:06:05]"
	TimeFormat string
}

type groupOrAttrs struct {
	group string
	attrs []slog.Attr
}

func NewHandler(out io.Writer, opts *Options) *developHandler {
	h := &developHandler{out: out, mu: &sync.Mutex{}}
	if opts != nil {
		h.opts = *opts

		if opts.HandlerOptions != nil {
			h.opts.HandlerOptions = opts.HandlerOptions
			if opts.Level == nil {
				h.opts.Level = slog.LevelInfo
			} else {
				h.opts.HandlerOptions.Level = opts.HandlerOptions.Level
			}
		} else {
			h.opts.HandlerOptions = &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}
		}

		if opts.MaxSlicePrintSize == 0 {
			h.opts.MaxSlicePrintSize = 50
		}

		if opts.TimeFormat == "" {
			h.opts.TimeFormat = "[15:06:05]"
		}
	} else {
		h.opts = Options{
			HandlerOptions:    &slog.HandlerOptions{Level: slog.LevelInfo},
			MaxSlicePrintSize: 50,
			TimeFormat:        "[15:06:05]",
		}
	}

	return h
}

func (h *developHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *developHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	return h.withGroupOrAttrs(groupOrAttrs{group: name})
}

func (h *developHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	return h.withGroupOrAttrs(groupOrAttrs{attrs: attrs})
}

func (h *developHandler) withGroupOrAttrs(goa groupOrAttrs) *developHandler {
	h2 := *h
	h2.goas = make([]groupOrAttrs, len(h.goas)+1)
	copy(h2.goas, h.goas)
	h2.goas[len(h2.goas)-1] = goa
	return &h2
}

func (h *developHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)
	buf = fmt.Appendf(buf, "%s ", csf(r.Time.Format(h.opts.TimeFormat), fgWhite))
	buf = h.formatSourceInfo(buf, &r)
	buf = h.levelMessage(buf, &r)
	buf = h.processAttributes(buf, &r)

	h.mu.Lock()
	defer h.mu.Unlock()

	_, err := h.out.Write(buf)

	return err
}

func (h *developHandler) formatSourceInfo(buf []byte, r *slog.Record) []byte {
	if h.opts.AddSource {
		at := cs("@@@", fgBlue)
		frame, _ := runtime.CallersFrames([]uintptr{r.PC}).Next()
		source := ul(cs(frame.File, fgYellow))
		line := cs(strconv.Itoa(frame.Line), fgRed)
		buf = fmt.Appendf(buf, "%s %s:%s\n", at, source, line)
	}

	return buf
}

func (h *developHandler) levelMessage(buf []byte, r *slog.Record) []byte {
	var bgColor color
	var fgColor color
	var lvlStr string
	if h.opts.ReplaceAttr != nil {
		a := h.opts.ReplaceAttr(nil, slog.Any(slog.LevelKey, r.Level))
		lvlStr = a.Value.String()
		if a.Key != "level" {
			r.AddAttrs(a)
		}
	} else {
		lvlStr = r.Level.String()
	}

	level := r.Level
	switch {
	case level < 0:
		bgColor, fgColor = bgBlue, fgBlue
	case level < 4:
		bgColor, fgColor = bgGreen, fgGreen
	case level < 8:
		bgColor, fgColor = bgYellow, fgYellow
	default:
		bgColor, fgColor = bgRed, fgRed
	}

	lvl := csb(" "+lvlStr+" ", fgBlack, bgColor)
	msg := cs(r.Message, fgColor)

	buf = fmt.Appendf(buf, "%s %s\n", lvl, msg)

	return buf
}

func (h *developHandler) processAttributes(buf []byte, r *slog.Record) []byte {
	var attrs attributes
	if r.NumAttrs() != 0 {
		r.Attrs(func(a slog.Attr) bool {
			attrs = append(attrs, a)
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
			newGroup := slog.Attr{
				Key:   goas[i].group,
				Value: slog.GroupValue(attrs...),
			}
			attrs = attributes{newGroup}
		} else {
			attrs = append(attrs, goas[i].attrs...)
		}
	}

	buf = h.colorize(buf, attrs, 0, []string{})
	buf = append(buf, '\n')
	return buf
}

func (h *developHandler) colorize(buf []byte, as attributes, level int, groups []string) []byte {
	if h.opts.SortKeys {
		sort.Sort(as)
	}

	keyColor := fgMagenta
	padding := as.padding(keyColor)

	for _, a := range as {
		if h.opts.ReplaceAttr != nil {
			a = h.opts.ReplaceAttr(groups, a)
		}

		key := cs(a.Key, keyColor)
		val := a.Value.String()
		var mark string
		switch a.Value.Kind() {
		case slog.KindFloat64, slog.KindInt64, slog.KindUint64:
			mark = cs("#", fgYellow)
			val = cs(val, fgYellow)
		case slog.KindBool:
			mark = cs("#", fgRed)
			val = cs(val, fgRed)
		case slog.KindString:
			if len(val) == 0 {
				val = csf("empty", fgWhite)
			} else if isURL(val) {
				mark = cs("*", fgBlue)
				val = cs(val, fgBlue)
			}
		case slog.KindTime, slog.KindDuration:
			mark = cs("@", fgCyan)
			val = cs(val, fgCyan)
		case slog.KindAny:
			a := a.Value.Any()
			err, isError := a.(error)
			if isError {
				mark = cs("E", fgRed)
				val = csb(fmt.Sprintf(" %v ", err), fgBlack, bgRed)
				break
			}

			jsonBytes, err := json.Marshal(a)
			if err != nil {
				break
			}

			var decodedData interface{}
			err = json.Unmarshal(jsonBytes, &decodedData)
			if err != nil {
				break
			}

			switch decoded := decodedData.(type) {
			case []interface{}:
				mark = cs("S", fgGreen)
				val = h.formatSlice(decoded, level)
			case map[string]interface{}:
				mark = cs("M", fgGreen)
				val = h.formatMap(decoded, level)
			}
		case slog.KindGroup:
			mark = cs("G", fgGreen)
			var groupAttrs attributes
			groupAttrs = a.Value.Group()
			groups = append(groups, a.Key)
			val = fmt.Sprintf("%v\n%s", cs("group", fgGreen), h.colorize(nil, groupAttrs, level+1, groups))
		}

		buf = fmt.Appendf(buf, "%*v%1v %-*s : %s", level*2, "", mark, padding, key, val)
		if a.Value.Kind() != slog.KindGroup {
			buf = append(buf, '\n')
		}
	}

	return buf
}

func isURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

func isSlice(s string) bool {
	return strings.HasPrefix(s, "slice[") && strings.HasSuffix(s, "]")
}

func isMap(s string) bool {
	return strings.HasPrefix(s, "map[") && strings.HasSuffix(s, "]")
}

// Splitting is done by SliceElementDivider
func (h *developHandler) formatSlice(s []interface{}, l int) string {
	if len(s) == 0 {
		return fmt.Sprintf("%v %v", cs("0", fgYellow), cs("slice[]", fgGreen))
	}

	length := cs(strconv.Itoa(len(s)), fgYellow)
	digits := len(strconv.Itoa(len(s)))
	if digits > 3 {
		digits = 3
	}

	elementColor := fgBlue
	res := fmt.Sprintf("%s %s", length, cs("slice[", fgGreen))
	for i, e := range s {
		if i == int(h.opts.MaxSlicePrintSize) {
			res += fmt.Sprintf("\n%*v%*s  %s%s", l*2+4, "", digits, "", cs("...", elementColor), cs("]", fgGreen))
			break
		}

		res += fmt.Sprintf("\n%*v%*s: %s", l*2+4, "", digits, cs(strconv.Itoa(i), fgGreen), cs(fmt.Sprint(e), elementColor))
		if i == len(s)-1 {
			res += fmt.Sprintf(" %s", cs("]", fgGreen))
		}
	}

	return res
}

// Splitting is done by SliceElementDivider,
func (h *developHandler) formatMap(s map[string]interface{}, level int) string {
	if len(s) == 0 {
		return fmt.Sprintf("%v %v", cs("0", fgYellow), cs("map[]", fgGreen))
	}

	var padding int
	for key := range s {
		color := len(cs(key, fgGreen))
		if color > padding {
			padding = color
		}
	}

	sortedKeys := sortDataMap(s)
	length := cs(strconv.Itoa(len(sortedKeys)), fgYellow)
	res := fmt.Sprintf("%s %s", length, cs("map[", fgGreen))
	for i, key := range sortedKeys {
		res += fmt.Sprintf("\n%*v%-*s : %s", level*2+4, "", padding, cs(key, fgGreen), cs(fmt.Sprint(s[key]), fgBlue))
		if i == len(sortedKeys)-1 {
			res += fmt.Sprintf(" %s", cs("]", fgGreen))
		}
	}

	return res
}

func sortDataMap(data map[string]interface{}) []string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}
