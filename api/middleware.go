package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/October-9th/simple-bank/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"

	// Declare the key for the payload so when store in the context, it can be retrieved by the next handler
	// The payload will be stored in the context with this specific key
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the authorization header
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Get the form a authorization header
		// Split the authorization header by space
		fields := strings.Fields(authorizationHeader)

		if len(fields) < 2 {
			err := errors.New("invalid authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Get the type of authorization header
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type: %s, expected %s", authorizationType, authorizationTypeBearer)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Get the access token and use it to verify the payload
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// If the payload is valid, store it in the context before parsing it to the next handler
		ctx.Set(authorizationPayloadKey, payload)
		// Forward the request to the next handler
		ctx.Next()
	}
}
