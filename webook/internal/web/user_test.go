package web

//
//import (
//	"bytes"
//	"github.com/stretchr/testify/require"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestUserHandler_Signup(t *testing.T) {
//	testCases := []struct {
//		name string
//	}{}
//
//	req, err := http.NewRequest(
//		http.MethodPost,
//		"/users/signup",
//		bytes.NewBuffer([]byte(`
//{
//	"email" : "123@qq.com",
//	"password" : "123456"
//}
//`)))
//	require.NoError(t, err)
//
//	resp := httptest.NewRecorder()
//	resp.Header()
//	resp.Body
//
//	h := NewUserHandler()
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//
//		})
//	}
//}
