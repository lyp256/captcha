package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/lyp256/captcha"
	cache2 "github.com/lyp256/captcha/pkg/kv"
)

func newServer() (*captchaServer, error) {
	provider, err := captcha.LoadDirImage("img")
	if err != nil {
		return nil, err
	}
	c := cache2.NewCURD(time.Minute)

	return &captchaServer{
		c: captcha.NewCaptcha(provider, c),
	}, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	s, err := newServer()
	if err != nil {
		log.Fatalln(err)
	}
	http.HandleFunc("/captcha", s.Captcha)
	http.HandleFunc("/image", s.Image)
	http.HandleFunc("/test", s.test)
	http.HandleFunc("/validate", s.Validate)
	http.Handle("/", http.FileServer(http.Dir("html")))
	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

type captchaServer struct {
	c *captcha.Captcha
}

// Captcha 生成一个验证码请求
func (s captchaServer) Captcha(writer http.ResponseWriter, request *http.Request) {
	en := json.NewEncoder(writer)
	key := strconv.FormatInt(time.Now().UnixNano(), 36)
	err := s.c.Generate(request.Context(), key)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_ = en.Encode(responseError{
			Message: err.Error(),
		})
		return
	}
	_ = en.Encode(responseCaptchaInfo{
		Key:      key,
		ImageURI: fmt.Sprintf("/image?key=%s", key),
	})
}

// Image 验证码图片
func (s captchaServer) Image(writer http.ResponseWriter, request *http.Request) {
	key := request.URL.Query().Get("key")
	img, err := s.c.DrawCaptcha(request.Context(), key)
	if err != nil {

		if errors.Is(err, cache2.ErrNotExist) {
			writer.WriteHeader(http.StatusNotFound)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
		}
		en := json.NewEncoder(writer)
		_ = en.Encode(responseError{
			Message: err.Error(),
		})
		return
	}
	writer.Header().Set("Content-Type", "image/png")
	_ = png.Encode(writer, img)
}

//  验证码图片
func (s captchaServer) test(writer http.ResponseWriter, request *http.Request) {
	img, err := s.c.Draw(request.Context(), "img/logo.png", randRadian())
	if err != nil {
		log.Println(err)
		return
	}
	writer.Header().Set("Content-Type", "image/png")
	_ = png.Encode(writer, img)
}

// Validate 验证验证码结果
func (s captchaServer) Validate(writer http.ResponseWriter, request *http.Request) {
	en := json.NewEncoder(writer)
	key := request.URL.Query().Get("key")
	radian, err := strconv.ParseFloat(request.URL.Query().Get("radian"), 63)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		_ = en.Encode(responseError{
			Message: fmt.Sprintf("radian:%s", err),
		})
		return
	}

	diff, err := s.c.Compare(request.Context(), key, radian)
	if err != nil {
		if errors.Is(err, cache2.ErrNotExist) {
			writer.WriteHeader(http.StatusNotFound)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
		}
		_ = en.Encode(responseError{
			Message: err.Error(),
		})
		return
	}

	// validate diff
	_ = en.Encode(responseValidate{
		Diff: diff,
	})
}

type responseCaptchaInfo struct {
	Key      string `json:"key"`
	ImageURI string `json:"imageURI"`
}

type responseValidate struct {
	Diff float64 `json:"diff"`
}

type responseError struct {
	Message string `json:"message"`
}

// return [0.14,2π-0.14]
func randRadian() float64 {
	const r = 2*math.Pi - 0.28
	return 0.14 + (rand.Float64() * r)
}
