package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockResponseWriter struct{}

func TestRouterGET(t *testing.T) {
	router := New()
	routed := false
	router.GET("/test", func(*Context) {
		routed = true
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	if !routed {
		t.Fatal("routing get failed")
	}
}
func TestRouterPOST(t *testing.T) {
	router := New()
	routed := false
	router.POST("/test", func(*Context) {
		routed = true
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/test", nil)
	router.ServeHTTP(w, req)

	if !routed {
		t.Fatal("routing post failed")
	}
}

func TestRouterGroup(t *testing.T) {
	router := New()
	routed := false
	router.Group("api", func() {
		router.GET("test", func(*Context) {
			routed = true
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/test", nil)
	router.ServeHTTP(w, req)

	if !routed {
		t.Fatal("routing group failed")
	}
}

func TestRouterUse(t *testing.T) {
	router := New()
	routed := false
	used := false
	router.Use(func(c *Context) {
		used = true
		c.Next()
	})
	router.Group("api", func() {
		router.GET("test", func(*Context) {
			routed = true
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/test", nil)
	router.ServeHTTP(w, req)

	if !used || !routed {
		t.Fatal("routing use failed")
	}
}

func TestRouterNotFound(t *testing.T) {
	router := New()
	router.POST("/test", func(*Context) {
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/test1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatal("NotFound handling route failed")
	}
}

func TestRouterMethodUnsupported(t *testing.T) {
	router := New()
	router.POST("/test", func(*Context) {
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Fatal("Method unsupported handling route failed")
	}
}
