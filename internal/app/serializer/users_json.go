package serializer

import "web_backend/internal/app/ds"

type SignInRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignUpRequest struct {
	Login       string `json:"login" binding:"required"`
	Password    string `json:"password" binding:"required"`
	IsModerator bool   `json:"is_moderator"`
}

type SignUpResponse struct {
	Login       string `json:"login"`
	IsModerator bool   `json:"is_moderator"`
}

func SignUpResponseFromUser(user ds.Users) SignUpResponse {
	return SignUpResponse{
		Login:       user.Login,
		IsModerator: user.IsModerator,
	}
}

func SignUpRequestToUser(j SignUpRequest) ds.Users {
	return ds.Users{
		Login:       j.Login,
		Password:    j.Password,
		IsModerator: j.IsModerator,
	}
}

type UserJSON struct {
	ID          uint   `json:"id,omitempty"`
	Login       string `json:"login"`
	Password    string `json:"password,omitempty"`
	IsModerator bool   `json:"is_moderator"`
}

func UserToJSON(user ds.Users) UserJSON {
	return UserJSON{
		ID:          user.UserID,
		Login:       user.Login,
		Password:    user.Password,
		IsModerator: user.IsModerator,
	}
}

func UserFromJSON(j UserJSON) ds.Users {
	return ds.Users{
		Login:       j.Login,
		Password:    j.Password,
		IsModerator: j.IsModerator,
	}
}
