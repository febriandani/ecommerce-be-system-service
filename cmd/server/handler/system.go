package handler

import (
	"context"
	"time"

	"github.com/febriandani/ecommerce-be-system-service/cmd/server/utils"
	systemPb "github.com/febriandani/ecommerce-be-system-service/protogen/golang/proto/system"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	client systemPb.SystemsClient
}

func NewHandler(client systemPb.SystemsClient) *Handler {
	return &Handler{client: client}
}

func (h *Handler) HandlerSystem_GetProvinces(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Parse incoming JSON
	var reqBody systemPb.Filter
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&systemPb.EmptyResponse{
			Meta: utils.NewMetaFail("Payload tidak valid: "+err.Error(), ""),
		})
	}

	// 3. Call gRPC
	res, err := h.client.GetProvinces(ctx, &reqBody)
	if err != nil {
		return utils.GRPCErrorToFiber(c, err, "Gagal mengirim notifikasi", "", "")
	}

	// 4. Handle response Meta dari gRPC
	if res.GetMeta().GetCode() != 200 {
		return c.Status(fiber.StatusBadRequest).JSON(&systemPb.EmptyResponse{
			Meta: res.GetMeta(),
		})
	}

	// 5. Success
	return c.Status(fiber.StatusCreated).JSON(&systemPb.EmptyResponse{
		Meta: utils.NewMetaSuccess("Notifikasi berhasil dikirim", res.GetMeta().GetTraceId()),
	})
}
