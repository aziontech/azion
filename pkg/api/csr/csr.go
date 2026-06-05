package csr

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"go.uber.org/zap"
)

// Create generates a Certificate Signing Request. The CSR endpoint is
// create-only and returns a Certificate whose Csr field holds the generated
// PEM-encoded signing request.
func (c *Client) Create(ctx context.Context, req *CreateRequest) (sdk.Certificate, error) {
	logger.Debug("Create Certificate Signing Request")

	request := c.apiClient.DigitalCertificatesCertificateSigningRequestsAPI.
		CreateCertificateSigningRequest(ctx).
		CertificateSigningRequest(req.CertificateSigningRequest)

	certResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Certificate Signing Request", zap.Error(err), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.Certificate{}, err
			}
		}
		return sdk.Certificate{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return certResponse.Data, nil
}

// Get retrieves a CSR's certificate. The CSR endpoint has no read operation,
// so the standard Digital Certificates endpoint is used instead.
func (c *Client) Get(ctx context.Context, id int64) (sdk.Certificate, error) {
	logger.Debug("Get Certificate Signing Request")
	request := c.apiClient.DigitalCertificatesCertificatesAPI.RetrieveCertificate(ctx, id)

	res, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting a Certificate Signing Request", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.Certificate{}, err
			}
		}
		return sdk.Certificate{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return res.Data, nil
}

// Delete removes a CSR's certificate. The CSR endpoint has no delete operation,
// so the standard Digital Certificates endpoint is used instead.
func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete Certificate Signing Request")
	req := c.apiClient.DigitalCertificatesCertificatesAPI.DeleteCertificate(ctx, id)

	_, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a Certificate Signing Request", zap.Error(err), zap.Any("ID", id))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}
