package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const maxUploadSize int64 = 8 << 20

// Claims 是应用使用的 JWT 声明。
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 使用密钥生成短时访问令牌。
func GenerateToken(userID int, secret []byte) (string, error) {
	if len(secret) < 32 {
		return "", errors.New("JWT_SECRET must contain at least 32 bytes")
	}
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-tutorial",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseToken 校验令牌签名算法、签名和声明有效期。
func ParseToken(value string, secret []byte) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(value, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	}, jwt.WithIssuer("go-tutorial"), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}

// NewRouter 创建文件上传和 JWT 示例的路由引擎。
func NewRouter(secret []byte, uploadDir string) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	// 限制 multipart 解析阶段使用的内存阈值；请求体上限仍由处理器单独控制。
	router.MaxMultipartMemory = maxUploadSize

	router.GET("/token", tokenHandler(secret))
	router.GET("/api/profile", JWTAuth(secret), profileHandler)
	router.POST("/api/upload", JWTAuth(secret), uploadHandler(uploadDir))
	router.Static("/uploads", uploadDir)

	return router
}

// JWTAuth 验证 Bearer 令牌并把用户 ID 写入请求上下文。
func JWTAuth(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		scheme, value, ok := strings.Cut(c.GetHeader("Authorization"), " ")
		if !ok || !strings.EqualFold(scheme, "Bearer") || value == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		claims, err := ParseToken(value, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

func tokenHandler(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := GenerateToken(1, secret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token configuration is invalid"})
			return
		}
		c.String(http.StatusOK, token)
	}
}

func profileHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user identity is missing"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}

func uploadHandler(uploadDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxUploadSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file"})
			return
		}
		if file.Size > maxUploadSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file too large"})
			return
		}

		ext := strings.ToLower(filepath.Ext(filepath.Base(file.Filename)))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported file type"})
			return
		}

		opened, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot inspect file"})
			return
		}
		header := make([]byte, 512)
		n, readErr := opened.Read(header)
		_ = opened.Close()
		if readErr != nil && !errors.Is(readErr, io.EOF) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cannot inspect file"})
			return
		}
		contentType := http.DetectContentType(header[:n])
		if !strings.HasPrefix(contentType, "image/") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file content is not an image"})
			return
		}

		if err := os.MkdirAll(uploadDir, 0o750); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upload directory unavailable"})
			return
		}
		name, err := randomName()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create file name"})
			return
		}
		destination := filepath.Join(uploadDir, name+ext)
		if err := c.SaveUploadedFile(file, destination); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save file"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"name": filepath.Base(destination)})
	}
}

func randomName() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate random name: %w", err)
	}
	return hex.EncodeToString(raw[:]), nil
}

func main() {
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) < 32 {
		log.Fatal("JWT_SECRET must contain at least 32 bytes")
	}
	if err := NewRouter(secret, "uploads").Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
