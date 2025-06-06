package config

import (
    "crypto/sha256"
    "strings"
    "time"

    "claude2api/logger"
    "github.com/imroc/req/v3"
)

func CheckAndFilterProxies(proxies []string) []string {
    usable := make([]string, 0, len(proxies))

    for _, p := range proxies {
        client := req.C().ImpersonateChrome().
            SetTimeout(10 * time.Second).
            SetProxyURL(p)

        resp, err := client.R().Get("https://claude.ai/login")
        if err != nil {
            logger.Warn("proxy %s unavailable: %v", p, err)
            continue
        }

        body, _ := resp.ToString()
        if strings.Contains(body, "Just a moment...") {
            logger.Warn("proxy %s blocked (cloudflare), dropping", p)
            continue
        }

        logger.Info("proxy %s OK", p)
        usable = append(usable, p)
    }

    return usable
}


func SelectProxyBySessionKey(sessionKey string, proxies []string) string {
    if len(proxies) == 0 {
        return ""
    }
    sum := sha256.Sum256([]byte(sessionKey))
    idx := int(sum[0]) % len(proxies)
    logger.Info("SelectProxy: sessionKey=%s -> proxyIdx=%d (%s)",sessionKey, idx, proxies[idx])
    return proxies[idx]
}
