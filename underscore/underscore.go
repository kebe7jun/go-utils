package underscore

import (
	"bytes"
	"encoding/json"
	"github.com/kebe7jun/go-utils/cache"
	"regexp"
	"strings"
)

var lruCache = cache.NewLRUCache(200)
var re = regexp.MustCompile(`"[a-zA-Z0-9_]+":`)

func UnderscoreJSONStr(src string) string {
	ss := re.FindAllString(src, -1)
	for _, s := range ss {
		src = strings.ReplaceAll(src, s, underscore(s))
	}
	return src
}

func UnderscoreObj(obj interface{}) string {
	bs, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return UnderscoreJSONStr(string(bs))
}

func underscore(src string) string {
	originSrc := src[:]
	cachedRes := lruCache.Get(originSrc)
	if cachedRes != nil {
		return cachedRes.(string)
	}
	res := bytes.NewBuffer([]byte{})
	lastIndex := -1
	for i := 0; i < len(src); i++ {
		if !(src[i] >= 'A' && src[i] <= 'Z') {
			res.WriteByte(src[i])
			continue
		}
		lower := src[i] + 32
		if lastIndex+1 != i {
			res.WriteByte('_')
		}
		res.WriteByte(lower)
		lastIndex = i
	}
	lruCache.Set(originSrc, res.String())
	return res.String()
}
