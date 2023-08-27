package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sinakeshmiri/switchman/internal/ports"
)

type Adapter struct {
	api ports.APIPort
}

// NewAdapter creates a new Adapter
func NewAdapter(api ports.APIPort) *Adapter {
	return &Adapter{api: api}
}

func (httpa Adapter) checkHealth(ctx *fiber.Ctx) error {
	err := httpa.api.Check()
	if err != nil {
		return ctx.Status(500).JSON(nil)
	}
	return ctx.Status(201).JSON(fiber.Map{"message": "ok"})
}

func (httpa Adapter) Check() {
	r := fiber.New()
	r.Use(recover.New())

	r.Get("/", httpa.checkHealth)
	r.Listen(":9000")
}
