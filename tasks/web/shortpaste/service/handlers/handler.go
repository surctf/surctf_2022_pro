package handlers

import (
	"encoding/base64"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/safehtml/template"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"shortpaste/config"
	"shortpaste/middleware"
	"shortpaste/models"
)

type Handlers struct {
	DB        *gorm.DB
	Templates map[string]*template.Template
}

func getTemplateClone(t *template.Template) *template.Template {
	clone, err := t.Clone()
	if err != nil {
		panic(err)
	}

	return clone
}

func NewHandlers(db *gorm.DB) Handlers {
	templates := make(map[string]*template.Template)

	base := template.Must(template.ParseFiles("templates/base.html"))
	templates["index.html"] = template.Must(getTemplateClone(base).ParseFiles("templates/index.html"))
	templates["paste.html"] = template.Must(getTemplateClone(base).ParseFiles("templates/paste.html"))
	return Handlers{DB: db, Templates: templates}
}

type Content struct {
	UserAuth  middleware.UserAuth
	Error     string
	WithError bool
}

type PasteContent struct {
	Content
	Paste          models.Paste
	AuthorUsername string
}

func (h *Handlers) GetPasteHandler(w http.ResponseWriter, r *http.Request) {
	userAuth, ok := r.Context().Value("userAuth").(middleware.UserAuth)
	if !ok {
		http.Error(w, "Can't get value with key 'userAuth' from context", http.StatusInternalServerError)
		return
	}
	respContent := Content{UserAuth: userAuth}

	pasteB64, ok := mux.Vars(r)["paste_id"]
	if !ok || pasteB64 == "" {
		http.Error(w, "No such paste((", http.StatusNotFound)
	}

	pasteID, err := base64.URLEncoding.DecodeString(pasteB64)
	if err != nil {
		http.Error(w, "No such paste((", http.StatusNotFound)
		return
	}

	var paste models.Paste
	if h.DB.Model(&paste).Find(&paste, "id = ?", pasteID).RowsAffected == 0 {
		http.Error(w, "No such paste((", http.StatusNotFound)
		return
	}

	var author models.User
	h.DB.Model(&author).Find(&author, "id = ?", paste.UserID)

	err = h.Templates["paste.html"].ExecuteTemplate(w, "base", PasteContent{
		Content:        respContent,
		Paste:          paste,
		AuthorUsername: config.USERNAME_PREFIX + strconv.Itoa(int(author.ID)),
	})

	if err != nil {
		log.Print(err)
		http.Error(w, "Can't execute 'paste.html' template", http.StatusInternalServerError)
		return
	}

}

func (h *Handlers) CreatePasteHandler(w http.ResponseWriter, r *http.Request) {
	userAuth, ok := r.Context().Value("userAuth").(middleware.UserAuth)
	if !ok {
		http.Error(w, "Can't get value with key 'userAuth' from context", http.StatusInternalServerError)
		return
	}
	respContent := Content{UserAuth: userAuth}

	if !userAuth.IsAuthorized {
		http.Error(w, "You need to be authorized to create pastes!", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Print(err)
		http.Error(w, "Can't parse 'sign_up' form", http.StatusInternalServerError)
		return
	}

	// Получаем название пасты из формы, проверяем на пустоту и длину
	title := r.PostForm.Get("title")
	if title == "" || len(title) > 32 {
		respContent.Error = "Длина названия пасты должна быть от 1 до 32 символоов!"
		respContent.WithError = true
		h.Templates["index.html"].ExecuteTemplate(w, "base", respContent)
		return
	}

	// Получаем текст пасты из формы, проверяем на пустоту и длину
	content := r.PostForm.Get("content")
	if content == "" || len(content) > 140 {
		respContent.Error = "Длина пасты должна быть от 1 до 140 символоов!"
		respContent.WithError = true
		h.Templates["index.html"].ExecuteTemplate(w, "base", respContent)
		return
	}

	err := userAuth.User.AddPaste(&models.Paste{Title: title, Content: content}, h.DB)
	if err != nil {
		log.Print(err)
		http.Error(w, "Can't create this paste((", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handlers) IndexHandler(w http.ResponseWriter, r *http.Request) {
	userAuth, ok := r.Context().Value("userAuth").(middleware.UserAuth)
	if !ok {
		http.Error(w, "Can't get value with key 'userAuth' from context", http.StatusInternalServerError)
		return
	}

	err := h.Templates["index.html"].ExecuteTemplate(w, "base", Content{UserAuth: userAuth})
	if err != nil {
		log.Print(err)
		http.Error(w, "Can't execute 'index.html' template", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	respContent := Content{UserAuth: middleware.UserAuth{IsAuthorized: false}}
	if err := r.ParseForm(); err != nil {
		log.Print(err)
		http.Error(w, "Can't parse 'sign_up' form", http.StatusInternalServerError)
		return
	}

	// Получаем пароль из формы, проверяем на пустоту
	password := r.PostForm.Get("password")
	if password == "" {
		respContent.Error = "Пароль не может быть пустым!"
		respContent.WithError = true
		h.Templates["index.html"].ExecuteTemplate(w, "base", respContent)
		return
	}

	// Хэшируем пароль
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		http.Error(w, "Can't hash password", http.StatusInternalServerError)
		return
	}

	// Создаем нового юзера
	newUser := models.User{PasswordHash: passwordHash}
	tx := h.DB.Create(&newUser)
	if tx.Error != nil {
		log.Print(err)
		http.Error(w, "Can't create new user", http.StatusInternalServerError)
		return
	}

	// Генерим jwt токен и подписываем его
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": newUser.ID,
	})
	tokenString, err := token.SignedString([]byte(config.HMAC_SECRET))
	if err != nil {
		log.Print(err)
		http.Error(w, "Can't sign token", http.StatusInternalServerError)
		return
	}

	// Устанавливаем jwt токен в куки 'auth' и редиректим в '/'
	tokenCookie := http.Cookie{Name: "auth", Value: tokenString}
	http.SetCookie(w, &tokenCookie)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handlers) SignInHandler(w http.ResponseWriter, r *http.Request) {
	respContent := Content{UserAuth: middleware.UserAuth{IsAuthorized: false}}
	if err := r.ParseForm(); err != nil {
		log.Print(err)
		http.Error(w, "Can't parse 'sign_in' form", http.StatusInternalServerError)
		return
	}

	// Получаем юзернейм из формы, проверяем на пустоту
	username := r.PostForm.Get("username")
	if username == "" {
		respContent.Error = "Юзернейм не может быть пустым!"
		respContent.WithError = true
		h.Templates["index.html"].ExecuteTemplate(w, "base", respContent)
		return
	}

	var user models.User
	// Проверяем, что юзернейм имеет правильный USERNAME_PREFIX, если всё ок, пытаемся спарсить ID и юзернейма
	if strings.HasPrefix(strings.ToLower(username), config.USERNAME_PREFIX) {
		user_id, err := strconv.Atoi(username[len(config.USERNAME_PREFIX):])
		if err == nil {
			if h.DB.Model(&user).Find(&user, "id = ?", user_id).Error != nil {
				http.Error(w, "Something went wrong when tried query DB", http.StatusInternalServerError)
				return
			}
		}
	}

	if user.ID == 0 {
		respContent.Error = "Неверные данные для входа!"
		respContent.WithError = true
		h.Templates["index.html"].ExecuteTemplate(w, "base", respContent)
		return
	}

	// Получаем пароль из формы, проверяем на пустоту
	password := r.PostForm.Get("password")
	if password == "" {
		respContent.Error = "Пароль не может быть пустым!"
		respContent.WithError = true
		h.Templates["index.html"].ExecuteTemplate(w, "base", respContent)
		return
	}

	// Сравниваем отправленный пароль с хэшем в бд
	cmpErr := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if cmpErr != nil {
		respContent.Error = "Неверные данные для входа!"
		respContent.WithError = true
		h.Templates["index.html"].ExecuteTemplate(w, "base", respContent)
		return
	}

	// Генерим jwt токен и подписываем его
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.ID,
	})
	tokenString, err := token.SignedString([]byte(config.HMAC_SECRET))
	if err != nil {
		log.Print(err)
		http.Error(w, "Can't sign token", http.StatusInternalServerError)
		return
	}

	// Устанавливаем jwt токен в куки 'auth' и редиректим в '/'
	tokenCookie := http.Cookie{Name: "auth", Value: tokenString}
	http.SetCookie(w, &tokenCookie)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handlers) LogOutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   "auth",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}
