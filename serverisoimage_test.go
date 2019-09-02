package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetServerIsoImageList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "isoimages")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerIsoImageListHTTPGet())
	})
	res, err := client.GetServerIsoImageList(dummyUUID)
	assert.Nil(t, err, "GetServerIsoImageList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockServerIsoImage()), fmt.Sprintf("%v", res))
}

func TestClient_GetServerIsoImage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "isoimages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerIsoImageHTTPget())
	})
	res, err := client.GetServerIsoImage(dummyUUID, dummyUUID)
	assert.Nil(t, err, "GetServerIsoImage returned an error %v", err)
	assert.Equal(t, fmt.Sprintf("%v", getMockServerIsoImage()), fmt.Sprintf("%v", res))
}

func TestClient_CreateServerIsoImage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "isoimages")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.CreateServerIsoImage(dummyUUID, ServerIsoImageRelationCreateRequest{
		ObjectUUID: dummyUUID,
	})
	assert.Nil(t, err, "CreateServerIsoImage returned an error %v", err)
}

func TestClient_UpdateServerIsoImage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "isoimages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.UpdateServerIsoImage(dummyUUID, dummyUUID, ServerIsoImageRelationUpdateRequest{
		BootDevice: true,
		Name:       "test",
	})
	assert.Nil(t, err, "UpdateServerIsoImage returned an error %v", err)
}

func TestClient_DeleteServerIsoImage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "isoimages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.DeleteServerIsoImage(dummyUUID, dummyUUID)
	assert.Nil(t, err, "DeleteServerIsoImage returned an error %v", err)
}

func TestClient_LinkIsoImage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "isoimages")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.LinkIsoImage(dummyUUID, dummyUUID)
	assert.Nil(t, err, "LinkIsoImage returned an error %v", err)
}

func TestClient_UnlinkIsoImage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "isoimages", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.UnlinkIsoImage(dummyUUID, dummyUUID)
	assert.Nil(t, err, "UnlinkIsoImage returned an error %v", err)
}

func getMockServerIsoImage() ServerIsoImageRelationProperties {
	mock := ServerIsoImageRelationProperties{
		ObjectUUID: dummyUUID,
		ObjectName: "test",
		Private:    false,
		CreateTime: dummyTime,
		Bootdevice: true,
	}
	return mock
}

func prepareServerIsoImageListHTTPGet() string {
	iso := getMockServerIsoImage()
	res, _ := json.Marshal(iso)
	return fmt.Sprintf(`{"isoimage_relations": [%s]}`, string(res))
}

func prepareServerIsoImageHTTPget() string {
	iso := getMockServerIsoImage()
	res, _ := json.Marshal(iso)
	return fmt.Sprintf(`{"isoimage_relation": %s}`, string(res))
}
