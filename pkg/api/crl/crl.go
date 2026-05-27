package crl

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"go.uber.org/zap"
)

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedCertificateRevocationList, error) {
	logger.Debug("List Certificate Revocation Lists")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}
	resp, httpResp, err := c.apiClient.DigitalCertificatesCertificateRevocationListsAPI.ListCertificateRevocationLists(ctx).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).
		Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing the certificate revocation lists", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}

func (c *Client) Get(ctx context.Context, id int64) (sdk.CertificateRevocationList, error) {
	logger.Debug("Get Certificate Revocation List")
	request := c.apiClient.DigitalCertificatesCertificateRevocationListsAPI.RetrieveCertificateRevocationList(ctx, id)

	res, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting a Certificate Revocation List", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.CertificateRevocationList{}, err
			}
		}
		return sdk.CertificateRevocationList{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return res.Data, nil
}

func (c *Client) Create(ctx context.Context, req *CreateRequest) (sdk.CertificateRevocationList, error) {
	logger.Debug("Create Certificate Revocation List")

	request := c.apiClient.DigitalCertificatesCertificateRevocationListsAPI.
		CreateCertificateRevocationList(ctx).
		CertificateRevocationList(req.CertificateRevocationList)

	crlResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Certificate Revocation List", zap.Error(err), zap.Any("Name", req.Name))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.CertificateRevocationList{}, err
			}
		}
		return sdk.CertificateRevocationList{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return crlResponse.Data, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest, id int64) (sdk.CertificateRevocationList, error) {
	logger.Debug("Update Certificate Revocation List", zap.Any("Certificate Revocation List ID", id))
	request := c.apiClient.DigitalCertificatesCertificateRevocationListsAPI.
		PartialUpdateCertificateRevocationList(ctx, id).
		PatchedCertificateRevocationList(req.PatchedCertificateRevocationList)

	crlResponse, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a Certificate Revocation List", zap.Error(err), zap.Any("ID", id))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.CertificateRevocationList{}, err
			}
		}
		return sdk.CertificateRevocationList{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return crlResponse.Data, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	logger.Debug("Delete Certificate Revocation List")
	req := c.apiClient.DigitalCertificatesCertificateRevocationListsAPI.DeleteCertificateRevocationList(ctx, id)

	_, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a Certificate Revocation List", zap.Error(err), zap.Any("ID", id))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}
