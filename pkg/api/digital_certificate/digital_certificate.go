package digitalcertificate

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"go.uber.org/zap"
)

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedCertificateList, error) {
	logger.Debug("List digital certificates")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}
	resp, httpResp, err := c.apiClient.DigitalCertificatesCertificatesAPI.ListCertificates(ctx).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).
		Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing the digital certificates", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}

func (c *Client) Get(ctx context.Context, id int64) (sdk.Certificate, error) {
	logger.Debug("Get Digital Certificate")
	request := c.apiClient.DigitalCertificatesCertificatesAPI.RetrieveCertificate(ctx, id)

	res, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting a Digital Certificate", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.Certificate{}, err
			}
		}
		return sdk.Certificate{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return res.Data, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete Digital Certificate")
	req := c.apiClient.DigitalCertificatesCertificatesAPI.DeleteCertificate(ctx, id)

	_, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a Digital Certificate", zap.Error(err), zap.Any("ID", id))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}

func (c *Client) Create(ctx context.Context, req *CreateRequest) (sdk.Certificate, error) {
	logger.Debug("Create Digital Certificate")

	request := c.apiClient.DigitalCertificatesCertificatesAPI.CreateCertificate(ctx).Certificate(req.Certificate)

	certResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Digital Certificate", zap.Error(err), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.Certificate{}, err
			}
		}
		return sdk.Certificate{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return certResponse.Data, nil
}

func (c *Client) Request(ctx context.Context, req *RequestRequest) (sdk.Certificate, error) {
	logger.Debug("Request Digital Certificate")

	request := c.apiClient.DigitalCertificatesRequestACertificateAPI.RequestCertificate(ctx).CertificateRequest(req.CertificateRequest)

	certResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while requesting a Digital Certificate", zap.Error(err), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.Certificate{}, err
			}
		}
		return sdk.Certificate{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return certResponse.Data, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest, id int64) (sdk.Certificate, error) {
	logger.Debug("Update Digital Certificate", zap.Any("Digital Certificate ID", id))
	request := c.apiClient.DigitalCertificatesCertificatesAPI.PartialUpdateCertificate(ctx, id).PatchedCertificate(req.PatchedCertificate)

	certResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a Digital Certificate", zap.Error(err), zap.Any("ID", id))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.Certificate{}, err
			}
		}
		return sdk.Certificate{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return certResponse.Data, nil
}
