package middleware

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"product-service/internal/auth"
	"product-service/pkg/logkey"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *Mid) Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// We get the current request context
		ctx := c.Request.Context()

		// Extract the traceId from the request context
		// We assert the type to string since context.Value returns an interface{}
		traceId, ok := ctx.Value(TraceIdKey).(string)

		if !ok {
			traceId = "unknown"
		}
		authHeader := c.Request.Header.Get("Authorization")

		// Splitting the Authorization header based on the space character.
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			// If the header format doesn't match required format, log and send an error
			err := errors.New("expected authorization header format: Bearer <token>")
			slog.Error("An error occurred",
				slog.Any(logkey.ERROR, err),
				slog.Any(logkey.TraceID, traceId),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		claims, err := m.k.ValidateToken(parts[1])
		if err != nil {
			slog.Error("Unauthorized User",
				slog.Any(logkey.ERROR, err),
				slog.Any(logkey.TraceID, traceId),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
			return
		}

		ctx = context.WithValue(ctx, auth.ClaimsKey, claims)
		c.Request = c.Request.WithContext(ctx)
		// Call the validate token from auth struct
		//put the validated claims in context
		// do the next thing in the chain

		c.Next()

	}
}

func (m *Mid) Authorize(next gin.HandlerFunc, requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		claims, ok := ctx.Value(auth.ClaimsKey).(auth.Claims)
		fmt.Println("printing claims roles ", claims.HasRoles())
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "claims missing from the context: Authorize called without/before Authenticate"})
			return
		}
		ok = claims.HasRoles(requiredRoles...)
		if !ok {
			slog.Error("An error occurred",
				slog.Any("claims", claims),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
			return
		}

		next(c)

	}
}
