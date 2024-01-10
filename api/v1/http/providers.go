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
	// drop existing provider embeddings
	routes.Delete("/providers/:uid/embeddings", s.DropProviderEmbeddings)
	// compute existing provider projections
	routes.Patch("/providers/:uid/projections", s.ComputeProviderProjections)
	// mount graph routes at the root of r
	r.Mount("/", routes)
}

// GetAllProviders returns all available embeddings providers.
// @Summary Get all embeddings providers.
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

	// NOTE(milosgajdos): we don't care if the conversion fails here.
	// If it does we'll get 0 and use the default values.
	offset := c.QueryInt("offset")
	limit := c.QueryInt("limit")

	if offset > 0 {
		filter.Offset = offset
	}

	if limit > 0 {
		filter.Limit = limit
	}

	providers, page, err := s.ProvidersService.GetProviders(context.TODO(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(v1.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(v1.ProvidersResponse{
		Providers: providers,
		Page:      page,
	})
}

// GetProviderByUID returns the embeddings provider with the given UID.
// @Summary Get embeddings provider by UID.
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

// GetProviderEmbeddings returns all embeddings for the provider with the given UID.
// @Summary Get embeddings by provider UID.
// @Description Returns embeddings for the provider with the given UID.
// @Tags providers
// @Produce json
// @Param id path string true "Provider UID"
// @Param offset query int false "Result offset"
// @Param limit query int false "Result limit"
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

	// NOTE(milosgajdos): we don't care if the conversion fails here.
	// If it does we'll get 0 and use the default values.
	offset := c.QueryInt("offset")
	limit := c.QueryInt("limit")

	if offset > 0 {
		filter.Offset = offset
	}

	if limit > 0 {
		filter.Limit = limit
	}

	embeddings, page, err := s.ProvidersService.GetProviderEmbeddings(context.Background(), uid.String(), filter)
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
		Page:       page,
	})
}

// GetProviderProjections returns all stored embedding projections for the given provider.
// @Summary Get embeddings projections by provider UID.
// @Description Returns embedding projections for the provider with the given UID.
// @Tags providers
// @Produce json
// @Param id path string true "Provider UID"
// @Param offset query int false "Result offset"
// @Param limit query int false "Result limit"
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

	// NOTE(milosgajdos): we don't care if the conversion fails here.
	// If it does we'll get 0 and use the default values.
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

	projections, page, err := s.ProvidersService.GetProviderProjections(context.Background(), uid.String(), filter)
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
		Projections: projections,
		Page:        page,
	})
}

// UpdateProviderEmbeddings fetches embeddings and updates provider records.
// @Summary Fetch and store embeddings for the provider with the given UID.
// @Description Update provider embeddings.
// @Tags providers
// @Accept json
// @Produce json
// @Param id path string true "Provider UID"
// @Param provider body v1.EmbeddingsUpdate true "Update provider embeddings"
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

	req := new(v1.EmbeddingsUpdate)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(v1.ErrorResponse{
			Error: err.Error(),
		})
	}

	if req.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(v1.ErrorResponse{
			Error: fmt.Sprintf("empty text provided to %s provider", uid.String()),
		})
	}

	if req.Projection != v1.PCA && req.Projection != v1.TSNE {
		return c.Status(fiber.StatusBadRequest).JSON(v1.ErrorResponse{
			Error: fmt.Sprintf("invalid projection: %v", req.Projection),
		})
	}

	ctx := context.Background()
	resp, err := FetchEmbeddings(ctx, embedder, req)
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

	emb, err := s.ProvidersService.UpdateProviderEmbeddings(ctx, uid.String(), *resp, req.Projection)
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

// DropProviderEmbeddings drops all embeddings of the provider with the given uid.
// @Summary Delete embeddings by provider UID.
// @Description Delete embeddings by provider UID. This also drops projections.
// @Tags providers
// @Produce json
// @Param uid path string true "Provider UID"
// @Success 204 {string} status "Provider embeddings deleted successfully"
// @Failure 400 {object} v1.ErrorResponse
// @Failure 500 {object} v1.ErrorResponse
// @Router /v1/providers/{uid}/embeddings [delete]
func (s *Server) DropProviderEmbeddings(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("uid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(v1.ErrorResponse{
			Error: err.Error(),
		})
	}

	if err := s.ProvidersService.DropProviderEmbeddings(context.Background(), uid.String()); err != nil {
		if code := v1.ErrorCode(err); code == v1.ENOTFOUND {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(v1.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ComputeProviderProjections recomputes provider projections from scratch by UID.
// @Summary Recompute embeddings projections for a provider by UID and return them.
// @Description Recompute provider projections.
// @Tags providers
// @Accept json
// @Produce json
// @Param id path string true "Provider UID"
// @Param provider body v1.ProjectionsUpdate true "Update embeddings projections"
// @Success 200 {object} v1.ProjectionsResponse
// @Success 200 {object} v1.Embedding
// @Failure 400 {object} v1.ErrorResponse
// @Failure 404 {object} v1.ErrorResponse
// @Failure 500 {object} v1.ErrorResponse
// @Router /v1/providers/{uid}/projections [patch]
func (s *Server) ComputeProviderProjections(c *fiber.Ctx) error {
	uid, err := uuid.Parse(c.Params("uid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(v1.ErrorResponse{
			Error: err.Error(),
		})
	}

	req := new(v1.ProjectionsUpdate)
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

	if req.Metadata == nil {
		req.Metadata = make(map[string]any)
	}
	req.Metadata["projection"] = req.Projection

	if err := s.ProvidersService.ComputeProviderProjections(context.Background(), uid.String(), req.Projection); err != nil {
		if code := v1.ErrorCode(err); code == v1.ENOTFOUND {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(v1.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusAccepted)
}
