package v2

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/freekieb7/gopenehr/database"
	"github.com/freekieb7/gopenehr/rest"
)

func Error(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func HandleCreateUser(db *database.Database) rest.HandlerFunc {
	type RequestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Name     string `json:"name"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		var reqBody RequestBody
		data, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(data, &reqBody); err != nil {
			Error(w, err)
			return nil
		}

		return nil

		// passwordHash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
		// if err != nil {
		// 	Error(w, err)
		// 	return nil
		// }

		// user, err := db.CreateUser(r.Context(), database.CreateUserParams{
		// 	Username:     reqBody.Username,
		// 	PasswordHash: string(passwordHash),
		// 	Email:        reqBody.Email,
		// 	Name:         reqBody.Name,
		// })
		// if err != nil {
		// 	Error(w, err)
		// 	return nil
		// }

		// return rest.JSON(w, user)
	}
}

func HandleListUsers() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}

func HandleGetUserById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleUpdateUserById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleDeleteUserById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleListFolders() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleCreateFolder() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleGetFolderById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleUpdateFolderById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleDeleteFolderById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleCreateDocument() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleListDocuments() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleGetDocumentById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleUpdateDocumentById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleDeleteDocumentById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleListDocumentRevisions() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleGetDocumentRevisionById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleUpdateDocumentRevisionById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleDeleteDocumentRevisionById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleGetAuditLog() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleCreateWebhook() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleListWebhooks() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleGetWebhookById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleUpdateWebhookById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleDeleteWebhookById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleListWebhookEvents() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleGetWebhookEventById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleGetCalendar() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleCreateCalendarEvent() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleListCalendarEvents() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleGetCalendarEventById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleUpdateCalendarEventById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}

func HandleDeleteCalendarEventById() rest.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		// todo
		return nil
	}
}
