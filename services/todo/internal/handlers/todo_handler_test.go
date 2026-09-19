package handlers

// import (
// 	"context"
// 	"database/sql"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
// 	"time"

// 	dtomodels "droplets_mini_todoservice/internal/models/dto_models"

// 	"github.com/go-chi/chi"
// )

// type fakeTaskService struct {
// 	createFunc    func(ctx context.Context, req *dtomodels.CreateRequest) (*dtomodels.CreateResponse, error)
// 	readAllFunc   func(ctx context.Context) (*dtomodels.GetAllResponse, error)
// 	readByPIDFunc func(ctx context.Context, pid string) (*dtomodels.GetTaskResponse, error)
// 	updateFunc    func(ctx context.Context, task *dtomodels.UpdateRequest, pid string) (*dtomodels.UpdateResponse, error)
// 	deleteFunc    func(ctx context.Context, pid string) error
// }

// func (f *fakeTaskService) CreateNewTaskItem(ctx context.Context, req *dtomodels.CreateRequest) (*dtomodels.CreateResponse, error) {
// 	return f.createFunc(ctx, req)
// }

// func (f *fakeTaskService) ReadAllTasks(ctx context.Context) (*dtomodels.GetAllResponse, error) {
// 	return f.readAllFunc(ctx)
// }

// func (f *fakeTaskService) ReadTaskByPID(ctx context.Context, pid string) (*dtomodels.GetTaskResponse, error) {
// 	return f.readByPIDFunc(ctx, pid)
// }

// func (f *fakeTaskService) UpdateTask(ctx context.Context, task *dtomodels.UpdateRequest, pid string) (*dtomodels.UpdateResponse, error) {
// 	return f.updateFunc(ctx, task, pid)
// }

// func (f *fakeTaskService) DeleteTask(ctx context.Context, pid string) error {
// 	return f.deleteFunc(ctx, pid)
// }

// func newTodoRouter(h *TodoHandler) http.Handler {
// 	r := chi.NewRouter()
// 	r.Post("/tasks", h.PostCreateTaskHandler)
// 	r.Get("/tasks", h.GetAllTaskHandler)
// 	r.Put("/tasks/{pid}", h.PutUpdateTaskHandler)
// 	r.Delete("/tasks/{pid}", h.DeleteTaskHandler)
// 	return r
// }

// func serve(h http.Handler, method, target, body string) *httptest.ResponseRecorder {
// 	var req *http.Request
// 	if body == "" {
// 		req = httptest.NewRequest(method, target, nil)
// 	} else {
// 		req = httptest.NewRequest(method, target, strings.NewReader(body))
// 	}
// 	rec := httptest.NewRecorder()
// 	h.ServeHTTP(rec, req)
// 	return rec
// }

// func decodeResponse[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
// 	t.Helper()
// 	var v T
// 	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil {
// 		t.Fatalf("unmarshal: %v", err)
// 	}
// 	return v
// }

// func TestPostCreateTaskHandler(t *testing.T) {
// 	now := time.Now().UTC()

// 	t.Run("created", func(t *testing.T) {
// 		h := &TodoHandler{srvc: &fakeTaskService{
// 			createFunc: func(_ context.Context, req *dtomodels.CreateRequest) (*dtomodels.CreateResponse, error) {
// 				return &dtomodels.CreateResponse{PID: "p-1", Value: req.Value, CreatedAt: now, Completed: false}, nil
// 			},
// 		}}

// 		rec := serve(newTodoRouter(h), http.MethodPost, "/tasks", `{"value":"  buy milk  "}`)
// 		if rec.Code != http.StatusCreated {
// 			t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
// 		}

// 		body := decodeResponse[dtomodels.CreateResponse](t, rec)
// 		if body.PID != "p-1" || body.Value != "buy milk" || body.Completed {
// 			t.Errorf("unexpected body: %+v", body)
// 		}
// 		if !body.CreatedAt.Equal(now) {
// 			t.Errorf("created_at = %v, want %v", body.CreatedAt, now)
// 		}
// 	})

// 	t.Run("invalid json returns 400", func(t *testing.T) {
// 		h := &TodoHandler{srvc: &fakeTaskService{
// 			createFunc: func(_ context.Context, _ *dtomodels.CreateRequest) (*dtomodels.CreateResponse, error) {
// 				t.Fatal("createFunc should not be called")
// 				return nil, nil
// 			},
// 		}}

// 		rec := serve(newTodoRouter(h), http.MethodPost, "/tasks", `{bad`)
// 		if rec.Code != http.StatusBadRequest {
// 			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
// 		}

// 		body := decodeResponse[dtomodels.ErrorResponse](t, rec)
// 		if body.Message != "invalid json" {
// 			t.Errorf("message = %q, want %q", body.Message, "invalid json")
// 		}
// 	})

// 	t.Run("whitespace value returns 400", func(t *testing.T) {
// 		h := &TodoHandler{srvc: &fakeTaskService{
// 			createFunc: func(_ context.Context, _ *dtomodels.CreateRequest) (*dtomodels.CreateResponse, error) {
// 				t.Fatal("createFunc should not be called")
// 				return nil, nil
// 			},
// 		}}

// 		rec := serve(newTodoRouter(h), http.MethodPost, "/tasks", `{"value":"   "}`)
// 		if rec.Code != http.StatusBadRequest {
// 			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
// 		}

// 		body := decodeResponse[dtomodels.ErrorResponse](t, rec)
// 		if body.Message != "missing value" {
// 			t.Errorf("message = %q, want %q", body.Message, "missing value")
// 		}
// 	})

// 	t.Run("service error returns 500", func(t *testing.T) {
// 		h := &TodoHandler{srvc: &fakeTaskService{
// 			createFunc: func(_ context.Context, _ *dtomodels.CreateRequest) (*dtomodels.CreateResponse, error) {
// 				return nil, errors.New("boom")
// 			},
// 		}}

// 		rec := serve(newTodoRouter(h), http.MethodPost, "/tasks", `{"value":"buy milk"}`)
// 		if rec.Code != http.StatusInternalServerError {
// 			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
// 		}

// 		body := decodeResponse[dtomodels.ErrorResponse](t, rec)
// 		if body.Message != "server error" {
// 			t.Errorf("message = %q, want %q", body.Message, "server error")
// 		}
// 	})
// }

// func TestGetAllTaskHandler(t *testing.T) {
// 	now := time.Now().UTC()

// 	t.Run("ok", func(t *testing.T) {
// 		h := &TodoHandler{srvc: &fakeTaskService{
// 			readAllFunc: func(_ context.Context) (*dtomodels.GetAllResponse, error) {
// 				return &dtomodels.GetAllResponse{Tasks: []dtomodels.GetTaskResponse{
// 					{PID: "p-1", Value: "first", Completed: false, CreatedAt: now},
// 					{PID: "p-2", Value: "second", Completed: true, CreatedAt: now},
// 				}}, nil
// 			},
// 		}}

// 		rec := serve(newTodoRouter(h), http.MethodGet, "/tasks", "")
// 		if rec.Code != http.StatusOK {
// 			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
// 		}

// 		body := decodeResponse[dtomodels.GetAllResponse](t, rec)
// 		if len(body.Tasks) != 2 {
// 			t.Fatalf("len(tasks) = %d, want 2", len(body.Tasks))
// 		}
// 		if body.Tasks[0].PID != "p-1" || body.Tasks[1].Completed != true {
// 			t.Errorf("unexpected body: %+v", body.Tasks)
// 		}
// 	})

// 	t.Run("service error returns 500", func(t *testing.T) {
// 		h := &TodoHandler{srvc: &fakeTaskService{
// 			readAllFunc: func(_ context.Context) (*dtomodels.GetAllResponse, error) {
// 				return nil, errors.New("boom")
// 			},
// 		}}

// 		rec := serve(newTodoRouter(h), http.MethodGet, "/tasks", "")
// 		if rec.Code != http.StatusInternalServerError {
// 			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
// 		}

// 		body := decodeResponse[dtomodels.ErrorResponse](t, rec)
// 		if body.Message != "server error" {
// 			t.Errorf("message = %q, want %q", body.Message, "server error")
// 		}
// 	})
// }

// func TestPutUpdateTaskHandler(t *testing.T) {
// 	now := time.Now().UTC()
// 	existing := &dtomodels.GetTaskResponse{PID: "p-1", Value: "old", Completed: false, CreatedAt: now}

// 	t.Run("updated", func(t *testing.T) {
// 		h := &TodoHandler{srvc: &fakeTaskService{
// 			readByPIDFunc: func(_ context.Context, pid string) (*dtomodels.GetTaskResponse, error) {
// 				return existing, nil
// 			},
// 			updateFunc: func(_ context.Context, _ *dtomodels.UpdateRequest, pid string) (*dtomodels.UpdateResponse, error) {
// 				return &dtomodels.UpdateResponse{PID: pid, Value: "new", Completed: true, CreatedAt: now}, nil
// 			},
// 		}}

// 		rec := serve(newTodoRouter(h), http.MethodPut, "/tasks/p-1", `{"value":"  new  ","completed":true}`)
// 		if rec.Code != http.StatusOK {
// 			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
// 		}

// 		body := decodeResponse[dtomodels.UpdateResponse](t, rec)
// 		if body.PID != "p-1" || body.Value != "new" || !body.Completed {
// 			t.Errorf("unexpected body: %+v", body)
// 		}
// 	})

// 	t.Run("not found returns 404", func(t *testing.T) {
// 		h := &TodoHandler{srvc: &fakeTaskService{
// 			readByPIDFunc: func(_ context.Context, _ string) (*dtomodels.GetTaskResponse, error) {
// 				return nil, sql.ErrNoRows
// 			},
// 		}}

// 		rec := serve(newTodoRouter(h), http.MethodPut, "/tasks/missing", `{"value":"new"}`)
// 		if rec.Code != http.StatusNotFound {
// 			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
// 		}

// 		body := decodeResponse[dtomodels.ErrorResponse](t, rec)
// 		if !strings.Contains(body.Message, "missing") {
// 			t.Errorf("message = %q, want it to mention pid", body.Message)
// 		}
// 	})

// 	t.Run("read error returns 500", func(t *testing.T) {
// 		h := &TodoHandler{srvc: &fakeTaskService{
// 			readByPIDFunc: func(_ context.Context, _ string) (*dtomodels.GetTaskResponse, error) {
// 				return nil, errors.New("db down")
// 			},
// 		}}

// 		rec := serve(newTodoRouter(h), http.MethodPut, "/tasks/p-1", `{"value":"new"}`)
// 		if rec.Code != http.StatusInternalServerError {
// 			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
// 		}
// 	})

// 	t.Run("invalid json returns 400", func(t *testing.T) {
// 		h := &TodoHandler{srvc: &fakeTaskService{
// 			readByPIDFunc: func(_ context.Context, _ string) (*dtomodels.GetTaskResponse, error) {
// 				return existing, nil
// 			},
// 		}}

// 		rec := serve(newTodoRouter(h), http.MethodPut, "/tasks/p-1", `{bad`)
// 		if rec.Code != http.StatusBadRequest {
// 			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
// 		}

// 		body := decodeResponse[dtomodels.ErrorResponse](t, rec)
// 		if body.Message != "invalid json" {
// 			t.Errorf("message = %q, want %q", body.Message, "invalid json")
// 		}
// 	})

// 	t.Run("whitespace value returns 400", func(t *testing.T) {
// 		h := &TodoHandler{srvc: &fakeTaskService{
// 			readByPIDFunc: func(_ context.Context, _ string) (*dtomodels.GetTaskResponse, error) {
// 				return existing, nil
// 			},
// 		}}

// 		rec := serve(newTodoRouter(h), http.MethodPut, "/tasks/p-1", `{"value":"   "}`)
// 		if rec.Code != http.StatusBadRequest {
// 			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
// 		}

// 		body := decodeResponse[dtomodels.ErrorResponse](t, rec)
// 		if body.Message != "missing value" {
// 			t.Errorf("message = %q, want %q", body.Message, "missing value")
// 		}
// 	})

// 	t.Run("update error returns 500", func(t *testing.T) {
// 		h := &TodoHandler{srvc: &fakeTaskService{
// 			readByPIDFunc: func(_ context.Context, _ string) (*dtomodels.GetTaskResponse, error) {
// 				return existing, nil
// 			},
// 			updateFunc: func(_ context.Context, _ *dtomodels.UpdateRequest, _ string) (*dtomodels.UpdateResponse, error) {
// 				return nil, errors.New("boom")
// 			},
// 		}}

// 		rec := serve(newTodoRouter(h), http.MethodPut, "/tasks/p-1", `{"value":"new"}`)
// 		if rec.Code != http.StatusInternalServerError {
// 			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
// 		}
// 	})
// }

// func TestDeleteTaskHandler(t *testing.T) {
// 	now := time.Now().UTC()
// 	existing := &dtomodels.GetTaskResponse{PID: "p-1", Value: "old", Completed: false, CreatedAt: now}

// 	t.Run("deleted", func(t *testing.T) {
// 		h := &TodoHandler{srvc: &fakeTaskService{
// 			readByPIDFunc: func(_ context.Context, _ string) (*dtomodels.GetTaskResponse, error) {
// 				return existing, nil
// 			},
// 			deleteFunc: func(_ context.Context, _ string) error {
// 				return nil
// 			},
// 		}}

// 		rec := serve(newTodoRouter(h), http.MethodDelete, "/tasks/p-1", "")
// 		if rec.Code != http.StatusNoContent {
// 			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
// 		}
// 	})

// 	t.Run("not found returns 404", func(t *testing.T) {
// 		h := &TodoHandler{srvc: &fakeTaskService{
// 			readByPIDFunc: func(_ context.Context, _ string) (*dtomodels.GetTaskResponse, error) {
// 				return nil, sql.ErrNoRows
// 			},
// 		}}

// 		rec := serve(newTodoRouter(h), http.MethodDelete, "/tasks/missing", "")
// 		if rec.Code != http.StatusNotFound {
// 			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
// 		}

// 		body := decodeResponse[dtomodels.ErrorResponse](t, rec)
// 		if !strings.Contains(body.Message, "missing") {
// 			t.Errorf("message = %q, want it to mention pid", body.Message)
// 		}
// 	})

// 	t.Run("read error returns 500", func(t *testing.T) {
// 		h := &TodoHandler{srvc: &fakeTaskService{
// 			readByPIDFunc: func(_ context.Context, _ string) (*dtomodels.GetTaskResponse, error) {
// 				return nil, errors.New("db down")
// 			},
// 		}}

// 		rec := serve(newTodoRouter(h), http.MethodDelete, "/tasks/p-1", "")
// 		if rec.Code != http.StatusInternalServerError {
// 			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
// 		}
// 	})

// 	t.Run("delete error returns 500", func(t *testing.T) {
// 		h := &TodoHandler{srvc: &fakeTaskService{
// 			readByPIDFunc: func(_ context.Context, _ string) (*dtomodels.GetTaskResponse, error) {
// 				return existing, nil
// 			},
// 			deleteFunc: func(_ context.Context, _ string) error {
// 				return errors.New("boom")
// 			},
// 		}}

// 		rec := serve(newTodoRouter(h), http.MethodDelete, "/tasks/p-1", "")
// 		if rec.Code != http.StatusInternalServerError {
// 			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
// 		}
// 	})
// }
