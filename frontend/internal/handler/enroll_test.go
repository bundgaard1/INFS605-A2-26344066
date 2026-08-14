package handler

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	coursecatalogue "osbourne.local/frontend/gen/course-catalogue"
	grpcclient "osbourne.local/frontend/internal/clients/grpc"
	"osbourne.local/frontend/ui"
)

type fakeCatalogueService struct {
	coursecatalogue.UnimplementedCourseCatalogueServiceServer
}

func (f *fakeCatalogueService) GetCourse(_ context.Context, req *coursecatalogue.GetCourseRequest) (*coursecatalogue.GetCourseResponse, error) {
	if req.GetCourseId() == "nope" {
		return nil, status.Error(codes.NotFound, "course not found")
	}
	return &coursecatalogue.GetCourseResponse{Course: &coursecatalogue.Course{
		Id:          req.GetCourseId(),
		Code:        "CS101",
		Title:       "Intro to CS",
		Description: "Basics",
		Credits:     10,
	}}, nil
}

func (f *fakeCatalogueService) EnrollUser(_ context.Context, _ *coursecatalogue.EnrollUserRequest) (*coursecatalogue.EnrollUserResponse, error) {
	return &coursecatalogue.EnrollUserResponse{Success: true}, nil
}

func (f *fakeCatalogueService) ListCourses(context.Context, *coursecatalogue.ListCoursesRequest) (*coursecatalogue.ListCoursesResponse, error) {
	return &coursecatalogue.ListCoursesResponse{}, nil
}

func (f *fakeCatalogueService) ListEnrolledCourses(context.Context, *coursecatalogue.ListEnrolledCoursesRequest) (*coursecatalogue.ListEnrolledCoursesResponse, error) {
	return &coursecatalogue.ListEnrolledCoursesResponse{}, nil
}

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	srv := grpc.NewServer()
	coursecatalogue.RegisterCourseCatalogueServiceServer(srv, &fakeCatalogueService{})
	go srv.Serve(lis)
	t.Cleanup(srv.Stop)

	clients, err := grpcclient.Dial("unused", "unused", lis.Addr().String(), "unused", "unused")
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(clients.Close)

	h := New(clients)
	router := h.Routes(ui.Files)

	ts := httptest.NewServer(router)
	t.Cleanup(ts.Close)
	return ts
}

func TestCoursePageNotFound(t *testing.T) {
	ts := newTestServer(t)

	resp, err := http.Get(ts.URL + "/courses/nope")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestGRPCToHTTPStatus(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want int
	}{
		{"not found", status.Error(codes.NotFound, "x"), http.StatusNotFound},
		{"permission denied", status.Error(codes.PermissionDenied, "x"), http.StatusForbidden},
		{"unauthenticated", status.Error(codes.Unauthenticated, "x"), http.StatusForbidden},
		{"unavailable", status.Error(codes.Unavailable, "x"), http.StatusServiceUnavailable},
		{"deadline exceeded", status.Error(codes.DeadlineExceeded, "x"), http.StatusServiceUnavailable},
		{"unknown default", status.Error(codes.Internal, "x"), http.StatusBadGateway},
		{"nil error", nil, http.StatusBadGateway},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := grpcToHTTPStatus(tc.err); got != tc.want {
				t.Errorf("grpcToHTTPStatus(%v) = %d, want %d", tc.err, got, tc.want)
			}
		})
	}
}

func TestCoursePageParamRoute(t *testing.T) {
	ts := newTestServer(t)

	resp, err := http.Get(ts.URL + "/courses/c1")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if !strings.Contains(string(body), "Intro to CS") {
		t.Errorf("course page missing course data (param not routed): %s", body)
	}
}

func TestHtmxStaticServed(t *testing.T) {
	ts := newTestServer(t)

	resp, err := http.Get(ts.URL + "/static/vendor/htmx.min.js")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if len(body) < 10000 {
		t.Fatalf("htmx.min.js looks truncated: %d bytes", len(body))
	}
}

func TestHandleEnrollCourseRoutes(t *testing.T) {
	ts := newTestServer(t)

	form := url.Values{"course_id": {"c1"}}
	resp, err := http.Post(ts.URL+"/api/courses/enroll", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	html := string(body)
	t.Logf("response body: %s", html)

	if !strings.Contains(html, `id="course-c1"`) {
		t.Errorf("response missing course card: %s", html)
	}
	if !strings.Contains(html, "Enroll In This Course") {
		t.Errorf("response missing enroll form: %s", html)
	}
	if strings.Contains(html, "Enrolled</span>") {
		t.Errorf("response should not contain enrolled badge: %s", html)
	}
	if !strings.Contains(html, `hx-swap-oob="beforeend:#toasts"`) {
		t.Errorf("response missing OOB toast: %s", html)
	}
	if !strings.Contains(html, "toast toast-success") {
		t.Errorf("response missing success toast: %s", html)
	}
}
