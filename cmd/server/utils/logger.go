package utils

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	clientauth "github.com/febriandani/ecommerce-be-system-service/protogen/golang/proto/system"
	"github.com/gofiber/fiber/v2"
)

func NewTraceID() string {
	return uuid.NewString()
}

// LoggerWithTrace is a wrapper around logrus.Entry with pre-filled fields
type LoggerWithTrace struct {
	entry *logrus.Entry
}

// NewLoggerWithTrace returns a new logger with traceID and request info
func NewLoggerWithTrace(base *logrus.Logger, traceID string, request interface{}) *LoggerWithTrace {
	return &LoggerWithTrace{
		entry: base.WithFields(logrus.Fields{
			"traceID": traceID,
			"request": request,
		}),
	}
}

func (l *LoggerWithTrace) Info(msg string) {
	l.entry.Info(msg)
}

func (l *LoggerWithTrace) Error(err error, msg string) {
	l.entry.WithError(err).Error(msg)
}

func (l *LoggerWithTrace) Warn(msg string) {
	l.entry.Warn(msg)
}

func (l *LoggerWithTrace) Debug(msg string) {
	l.entry.Debug(msg)
}

func GRPCErrorToFiber(c *fiber.Ctx, err error, fallbackMsgID, fallbackMsgEN, traceId string) error {
	st, ok := status.FromError(err)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(&clientauth.EmptyResponse{
			Meta: NewMetaError(fiber.StatusInternalServerError, fallbackMsgID, err.Error(), traceId),
		})
	}

	switch st.Code() {
	case codes.AlreadyExists:
		return c.Status(fiber.StatusBadRequest).JSON(&clientauth.EmptyResponse{
			Meta: NewMetaError(fiber.StatusBadRequest, fallbackMsgID, st.Message(), traceId),
		})
	case codes.InvalidArgument:
		return c.Status(fiber.StatusBadRequest).JSON(&clientauth.EmptyResponse{
			Meta: NewMetaError(fiber.StatusBadRequest, fallbackMsgID, st.Message(), traceId),
		})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(&clientauth.EmptyResponse{
			Meta: NewMetaError(fiber.StatusInternalServerError, fallbackMsgID, st.Message(), traceId),
		})
	}
}
