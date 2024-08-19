package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/signature"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/store"
)

type NewSignatureDevice struct {
	SignatureAlg string `json:"signature_alg"`
	Label string `json:"label"`
}

type SignatureDevice struct {
	ID string `json:"id"`
	SignatureAlg string `json:"signature_alg"`
	Label string `json:"label"`
	SignatureCounter int `json:"signature_counter"`
}

type SignatureReq struct {
	DataToBeSigned string `json:"data_to_be_signed"`
}

type SignatureResp struct {
	Signature string `json:"signature"`
	SignedData string `json:"signed_data"`
}

type APIError struct {
	Message string `json:"message"`
}

func (s *Server) CreateSigningDevice(response http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")

	var device NewSignatureDevice
	err := json.NewDecoder(request.Body).Decode(&device)
	if err != nil {
		WriteAPIResponse(response, http.StatusBadRequest, APIError{
			Message: "invalid request payload",
		})
		return
	}

	if err = s.signatureService.CreateSignatureDevice(request.Context(), signature.NewSignatureDevice{
		ID: id,
		Tenant: "1", // we do not care about the tenant at this stage
		SignatureAlg: device.SignatureAlg,
		Label: device.Label,
	}); err != nil {
		WriteAPIResponse(response, http.StatusBadRequest, APIError{
			Message: fmt.Sprintf("error creating a signature device: %v", err), // we should handle that better
		})
	}
	

	WriteAPIResponse(response, http.StatusOK, "OK") // TODO: we should return a more structured response with a link to the resource
}

func (s *Server) GetSigningDevice(response http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")

	signDevice, err := s.signatureService.GetSignatureDevice(request.Context(), id)
	if err != nil {
		if errors.Is(store.ErrDeviceNotFound, err) {
			WriteAPIResponse(response, http.StatusNotFound, APIError{
				Message: "signature device not found",
			})
			return
		}

		WriteAPIResponse(response, http.StatusInternalServerError, APIError{
			Message: fmt.Sprintf("error getting signature device: %v", err),
		})
		return
	}

	WriteAPIResponse(response, http.StatusOK, SignatureDevice{
		ID: signDevice.ID,
		SignatureAlg: signDevice.SignatureAlg,
		Label: signDevice.Label,
		SignatureCounter: signDevice.SignatureCounter,
	}) // TODO: we should return a more structured response with a link to the resource
}

func (s *Server) SignData(response http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")

	var signatureReq SignatureReq
	err := json.NewDecoder(request.Body).Decode(&signatureReq)
	if err != nil {
		WriteAPIResponse(response, http.StatusBadRequest, APIError{
			Message: "invalid request payload",
		})
		return
	}

	signedData, err := s.signatureService.SignData(request.Context(), id, signatureReq.DataToBeSigned)
	if err != nil {
		if errors.Is(store.ErrDeviceNotFound, err) {
			WriteAPIResponse(response, http.StatusNotFound, APIError{
				Message: "signature device not found",
			})
			return
		}

		WriteAPIResponse(response, http.StatusInternalServerError, APIError{
			Message: fmt.Sprintf("error signing the message: %v", err),
		})
		return
	}

	WriteAPIResponse(response, http.StatusOK, SignatureResp{
		Signature: signedData.Signature,
		SignedData: signedData.SignedData,
	})
}