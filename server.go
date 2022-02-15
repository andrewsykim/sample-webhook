package main

import (
	"io"
	"net/http"
	"time"

	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/klog/v2"
)

type server struct {
	scheme *runtime.Scheme
}

func (s *server) validate(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	data, err := io.ReadAll(req.Body)
	if err != nil {
		klog.ErrorS(err, "failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	codec := json.NewSerializerWithOptions(json.DefaultMetaFactory, s.scheme, s.scheme, json.SerializerOptions{})
	obj, _, err := codec.Decode(data, nil, nil)
	if err != nil {
		klog.ErrorS(err, "failed to decode request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	admissionReview, ok := obj.(*admissionv1.AdmissionReview)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	time.Sleep(5 * time.Second)

	admissionReview.Response = &admissionv1.AdmissionResponse{
		Allowed: true,
	}

	if err := codec.Encode(admissionReview, w); err != nil {
		klog.ErrorS(err, "error encoding admission response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *server) mutate(w http.ResponseWriter, r *http.Request) {
}
