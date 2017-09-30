package router

import (
	"fmt"
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
	req, _ := http.NewRequest("GET", "/test", nil)
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
	req, _ := http.NewRequest("POST", "/test", nil)
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
	req, _ := http.NewRequest("GET", "/api/test", nil)
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
	req, _ := http.NewRequest("GET", "/api/test", nil)
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
	req, _ := http.NewRequest("POST", "/test1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatal("NotFound handling route failed")
	}
}

//func TestRouterMethodUnsupported(t *testing.T) {
//	router := New()
//	router.POST("/test", func(*Context) {
//	})
//
//	w := httptest.NewRecorder()
//	req, _ := http.NewRequest("GET", "/test", nil)
//	router.ServeHTTP(w, req)
//	if w.Code != http.StatusNotImplemented {
//		t.Fatal("Method unsupported handling route failed")
//	}
//}

func TestRouterWithSameURL(t *testing.T) {
	router := New()
	router.Group("/api", func() {
		router.GET("/test")
		router.POST("/test")
	})

	for k, r := range router.routers {
		fmt.Printf("method is %s, path is %v\n", r.method, k)
	}
}
