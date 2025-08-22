package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var jwtSecret []byte



type RedisRateLimiter struct {
	RedisClient *redis.Client
	MaxRequests int        // 35
	Window      time.Duration // 1 daqiqa
	BlockTime   time.Duration // 20 daqiqa
}

func NewRedisRateLimiter(redisClient *redis.Client) *RedisRateLimiter {
	return &RedisRateLimiter{
		RedisClient: redisClient,
		MaxRequests: 35,
		Window:      time.Minute,
		BlockTime:   20 * time.Minute,
	}
}

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not set in .env file")
	}

	jwtSecret = []byte(secret)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const bearerPrefix = "Bearer "
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token topilmadi"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Noto‘g‘ri Authorization formati"})
			c.Abort()
			return
		}

		tokenString := authHeader[len(bearerPrefix):]
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token mavjud emas"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token muddati o'tgan"})
				c.Abort()
				return
			}
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Imzo noto‘g‘ri"})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Tokenni tekshirishda xato: " + err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token yaroqsiz"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token ma'lumotlari o'qib bo'lmadi"})
			c.Abort()
			return
		}

		// user_id olish va stringga o'tkazish
		var userID string
		if uid, ok := claims["user_id"].(string); ok {
			userID = uid
		} else if uidFloat, ok := claims["user_id"].(float64); ok {
			userID = strconv.FormatInt(int64(uidFloat), 10)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token ichida user_id mavjud emas yoki noto'g'ri formatda"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := []string{
			"http://localhost:5173",
			"http://localhost:3000",
		}

		origin := c.Request.Header.Get("Origin")
		var allowedOrigin string
		for _, o := range allowedOrigins {
			if o == origin {
				allowedOrigin = o
				break
			}
		}

		if allowedOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Writer.Header().Add("Vary", "Origin")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}


func AdminOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
            c.Abort()
            return
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

        // Tokenni ParseUnverified bilan parse qilish
        token, _, err := jwt.NewParser().ParseUnverified(tokenStr, jwt.MapClaims{})
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
            c.Abort()
            return
        }

        role, ok := claims["role"].(string)
        if !ok || role != "admin" {
            c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
            c.Abort()
            return
        }

        // Hammasi to‘g‘ri bo‘lsa, route davom etadi
        c.Next()
    }
}



// internal/handler/middleware.go

func (r *RedisRateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		// Endpointni aniqlash: method + path (masalan: GET:/api/v1/profile)
		routeKey := fmt.Sprintf("%s:%s", c.Request.Method, c.FullPath())

		// Agar FullPath bo'sh bo'lsa (masalan, route nomlanmagan), RequestURI yoki Path
		if routeKey == ":" {
			routeKey = fmt.Sprintf("%s:%s", c.Request.Method, c.Request.URL.Path)
		}

		ctx := context.Background()

		// Har bir endpoint uchun alohida kalit
		key := "rate_limit:" + clientIP + ":" + routeKey
		blockKey := "blocked:" + clientIP + ":" + routeKey // Har bir endpoint uchun alohida block

		// 1. Bloklanganligini tekshir
		isBlocked, err := r.RedisClient.Get(ctx, blockKey).Result()
		if err == nil && isBlocked == "1" {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Siz ushbu so'rovni yuborish chegarasiga yetdingiz. Iltimos, 20 daqiqa kutib turing.",
			})
			c.Abort()
			return
		}

		// 2. Joriy tokenlar
		result, err := r.RedisClient.Get(ctx, key).Result()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Serverda xatolik yuz berdi."})
			c.Abort()
			return
		}

		var tokens int
		var resetTime time.Time

		if err == redis.Nil {
			// Birinchi marta so'rov
			tokens = r.MaxRequests - 1
			resetTime = time.Now().Add(r.Window)

			r.RedisClient.Set(ctx, key, tokens, r.Window)
			r.RedisClient.Set(ctx, key+":reset", resetTime.Unix(), r.Window)
		} else {
			tokens, _ = strconv.Atoi(result)
			resetUnix, _ := r.RedisClient.Get(ctx, key+":reset").Int64()
			resetTime = time.Unix(resetUnix, 0)

			if tokens <= 0 {
				// Limit tugadi → bloklaymiz
				r.RedisClient.Set(ctx, blockKey, "1", r.BlockTime)
				r.RedisClient.Del(ctx, key, key+":reset")

				c.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Siz ushbu amalni bajarish chegarasiga yetdingiz. 20 daqiqa davomida so'rov yuborolmaysiz.",
				})
				c.Abort()
				return
			}

			// Tokenni kamaytiramiz
			tokens--
			r.RedisClient.Set(ctx, key, tokens, time.Until(resetTime))
		}

		// Qo'shimcha: header orqali info berish
		c.Header("X-RateLimit-Limit", strconv.Itoa(r.MaxRequests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(tokens))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		c.Next()
	}
}