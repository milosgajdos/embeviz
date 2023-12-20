package http

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	v1 "github.com/milosgajdos/embeviz/api/v1"
)

func (s *Server) registerProviderRoutes(r fiber.Router) {
	routes := fiber.New()
	// get all providers stored in the database
	routes.Get("/providers", s.GetAllProviders)
	// get a provider by UID
	routes.Get("/providers/:uid", s.GetProviderByUID)
	// get provider embeddings
	routes.Get("/providers/:uid/embeddings", s.GetProviderEmbeddings)
	// get provider projections
	routes.Get("/providers/:uid/projections", s.GetProviderProjections)
	// update existing provider embeddings
	routes.Put("/providers/:uid/embeddings", s.UpdateProviderEmbeddings)
	// mount graph routes at the root of r
	r.Mount("/", routes)
}

// GetAllProviders returns all available providers.
// @Summary Get all providers.
// @Description Get all available providers.
// @Tags providers
// @Produce json
// @Param offset query int false "Result offset"
// @Param limit query int false "Result limit"
// @Success 200 {object} v1.ProvidersResponse
// @Failure 500 {object} v1.ErrorResponse
// @Router /v1/providers [get]
func (s *Server) GetAllProviders(c *fiber.Ctx) error {
	var filter v1.ProviderFilter
	filter.Limit = v1.DefaultLimit

	// NOTE(milosgajdos): we don't care if the conversion fails
	offset := c.QueryInt("offset")
	limit := c.QueryInt("limit")

	if offset > 0 {
		filter.Offset = offset
	}

	if limit > 0 {
		filter.Limit = limit
	}

	providers, n, err := s.ProvidersService.GetProviders(context.TODO(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(v1.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(v1.ProvidersResponse{
		Providers: providers,
		N:         n,
	})
}

// GetProviderByUID returns the provider with the given UID.
// @Summary Get provider by UID.
// @Description Returns embeddings provider with the given UID.
// @Tags providers
// @Produce json
// @Param id path string true "Provider UID"
// @Success 200 {object} v1.Provider
// @Failure 400 {object} v1.ErrorResponse
// @Failure 404 {object} v1.ErrorResponse
// @Failure 500 {object} v1.ErrorResponse
// @Router /v1/providers/{uid} [get]
func (s *Server) GetProviderByUID(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("uid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(v1.ErrorResponse{
			Error: err.Error(),
		})
	}

	provider, err := s.ProvidersService.GetProviderByUID(context.TODO(), uid.String())
	if err != nil {
		if code := v1.ErrorCode(err); code == v1.ENOTFOUND {
			return c.Status(fiber.StatusNotFound).JSON(v1.ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(v1.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(provider)
}

// GetProviderEmbeddings returns all stored embeddings for the given provider.
// @Summary Get provider embedding by UID.
// @Description Returns embeddings for the provider with the given UID.
// @Tags providers
// @Produce json
// @Param id path string true "Provider UID"
// @Success 200 {object} v1.EmbeddingsResponse
// @Failure 400 {object} v1.ErrorResponse
// @Failure 404 {object} v1.ErrorResponse
// @Failure 500 {object} v1.ErrorResponse
// @Router /v1/providers/{uid}/embeddings [get]
func (s *Server) GetProviderEmbeddings(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("uid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(v1.ErrorResponse{
			Error: err.Error(),
		})
	}

	var filter v1.ProviderFilter
	filter.Limit = v1.DefaultLimit

	// NOTE(milosgajdos): we don't care if the conversion fails
	offset := c.QueryInt("offset")
	limit := c.QueryInt("limit")

	if offset > 0 {
		filter.Offset = offset
	}

	if limit > 0 {
		filter.Limit = limit
	}

	embeddings, n, err := s.ProvidersService.GetProviderEmbeddings(context.TODO(), uid.String(), filter)
	if err != nil {
		if code := v1.ErrorCode(err); code == v1.ENOTFOUND {
			return c.Status(fiber.StatusNotFound).JSON(v1.ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(v1.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(v1.EmbeddingsResponse{
		Embeddings: embeddings,
		N:          n,
	})
}

// GetProviderProjections returns all stored embedding projections for the given provider.
// @Summary Get provider embedding projections by UID.
// @Description Returns embedding projections for the provider with the given UID.
// @Tags providers
// @Produce json
// @Param id path string true "Provider UID"
// @Success 200 {object} v1.ProjectionsResponse
// @Failure 400 {object} v1.ErrorResponse
// @Failure 404 {object} v1.ErrorResponse
// @Failure 500 {object} v1.ErrorResponse
// @Router /v1/providers/{uid}/projections [get]
func (s *Server) GetProviderProjections(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("uid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(v1.ErrorResponse{
			Error: err.Error(),
		})
	}

	var filter v1.ProviderFilter
	filter.Limit = v1.DefaultLimit

	// NOTE(milosgajdos): we don't care if the conversion fails
	offset := c.QueryInt("offset")
	limit := c.QueryInt("limit")
	if offset > 0 {
		filter.Offset = offset
	}
	if limit > 0 {
		filter.Limit = limit
	}
	dim := c.Query("dim")
	if dim != "" {
		switch strings.ToUpper(dim) {
		case string(v1.Dim2D):
			dim2d := v1.Dim2D
			filter.Dim = &dim2d
		case string(v1.Dim3D):
			dim3d := v1.Dim3D
			filter.Dim = &dim3d
		}
	}

	embeddings, n, err := s.ProvidersService.GetProviderProjections(context.TODO(), uid.String(), filter)
	if err != nil {
		if code := v1.ErrorCode(err); code == v1.ENOTFOUND {
			return c.Status(fiber.StatusNotFound).JSON(v1.ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(v1.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(v1.ProjectionsResponse{
		Embeddings: embeddings,
		N:          n,
	})
}

// UpdateProviderEmbeddings fetches embeddings and updates provider records.
// @Summary Fetch embeddings and update the store for the provider with the given UID.
// @Description Update provider embeddings.
// @Tags providers
// @Accept json
// @Produce json
// @Param id path string true "Provider UID"
// @Param provider body v1.EmbeddingUpdate true "Update a provider"
// @Success 200 {object} v1.Embedding
// @Failure 400 {object} v1.ErrorResponse
// @Failure 404 {object} v1.ErrorResponse
// @Failure 500 {object} v1.ErrorResponse
// @Router /v1/providers/{uid}/embeddings [put]
func (s *Server) UpdateProviderEmbeddings(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("uid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(v1.ErrorResponse{
			Error: err.Error(),
		})
	}

	embedder, ok := s.Embedders[uid.String()]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(v1.ErrorResponse{
			Error: fmt.Sprintf("%s provider not found", uid.String()),
		})
	}

	// TODO: validate payload
	req := new(v1.EmbeddingUpdate)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(v1.ErrorResponse{
			Error: err.Error(),
		})
	}

	if req.Projection != v1.PCA && req.Projection != v1.TSNE {
		return c.Status(fiber.StatusBadRequest).JSON(v1.ErrorResponse{
			Error: fmt.Sprintf("invalid projection: %v", req.Projection),
		})
	}

	resp, err := FetchEmbeddings(context.TODO(), embedder, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(v1.ErrorResponse{
			Error: err.Error(),
		})
	}

	if req.Metadata == nil {
		req.Metadata = make(map[string]any)
	}
	req.Metadata["projection"] = req.Projection
	if req.Label != "" {
		req.Metadata["label"] = req.Label
	}
	resp.Metadata = req.Metadata

	proj := req.Projection

	emb, err := s.ProvidersService.UpdateProviderEmbeddings(context.TODO(), uid.String(), *resp, proj)
	if err != nil {
		if code := v1.ErrorCode(err); code == v1.EINVALID {
			return c.Status(fiber.StatusBadRequest).JSON(v1.ErrorResponse{
				Error: err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(v1.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(emb)
}
