package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/chinaboard/brewing/pkg/bininfo"
	"github.com/chinaboard/brewing/pkg/videodl"
	"github.com/chinaboard/brewing/pkg/whisper"
	"log"
	"os"
	"strings"
)

const (
	filePath = "/tmp/temp.m4a"
)

func main() {
	whisperEndpoint := ""
	videoUrl := ""
	version := false
	language := "zh"
	flag.StringVar(&whisperEndpoint, "whisperEndpoint", "http://whisper:9000", "e.g http://whisper:9000")
	flag.StringVar(&language, "language", "zh", "Available values : af, am, ar, as, az, ba, be, bg, bn, bo, br, bs, ca, cs, cy, da, de, el, en, es, et, eu, fa, fi, fo, fr, gl, gu, ha, haw, he, hi, hr, ht, hu, hy, id, is, it, ja, jw, ka, kk, km, kn, ko, la, lb, ln, lo, lt, lv, mg, mi, mk, ml, mn, mr, ms, mt, my, ne, nl, nn, no, oc, pa, pl, ps, pt, ro, ru, sa, sd, si, sk, sl, sn, so, sq, sr, su, sv, sw, ta, te, tg, th, tk, tl, tr, tt, uk, ur, uz, vi, yi, yo, zh")
	flag.StringVar(&videoUrl, "videoUrl", "", "e.g https://www.youtube.com/watch?v=xxxxx")
	flag.BoolVar(&version, "v", false, "print version")
	flag.Parse()

	if version {
		fmt.Println(bininfo.StringifySingleLine())
		os.Exit(0)
	}

	if whisperEndpoint == "" {
		fmt.Println("Error: whisperEndpoint must be set")
		os.Exit(1)
	}
	if videoUrl == "" {
		fmt.Println("Error: videoUrl must be set")
		os.Exit(1)
	}

	whisperEndpoint = strings.TrimSuffix(whisperEndpoint, "/")

	if err := videodl.Download(videoUrl, filePath); err != nil {
		log.Fatalln(err)
	}

	resp, err := whisper.Asr(whisperEndpoint, filePath, language)
	if err != nil {
		log.Fatalln(err)
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(data))
}
