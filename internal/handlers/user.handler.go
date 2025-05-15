package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/momokii/go-sso-web/internal/database"
	"github.com/momokii/go-sso-web/pkg/models"
	"github.com/momokii/go-sso-web/pkg/repository/user"
	"github.com/momokii/go-sso-web/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userRepo user.UserRepo
}

func NewUserHandler(userRepo user.UserRepo) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

func (h *UserHandler) Change2FAStatus(c *fiber.Ctx) error {

	// change 2fa just reversed it, if true so make it false and the otherwise

	user := c.Locals("user").(models.UserSession)

	tx, err := database.DB.Begin()
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to start transaction")
	}
	defer func() {
		database.CommitOrRollback(tx, c, err)
	}()

	user_data, err := h.userRepo.FindByID(tx, user.Id)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to check user")
	}

	if user_data.Id == 0 {
		return utils.ResponseError(c, fiber.StatusBadRequest, "User not found")
	}

	if user_data.Id != user.Id {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid request")
	}

	// just reserve the 2fa staus on update
	var reversed_multifa_status bool
	if user_data.MultiFAEnabled {
		reversed_multifa_status = false
	} else {
		reversed_multifa_status = true
	}
	user_update := models.User{
		Id:             user_data.Id,
		Username:       user_data.Username,
		PhoneNumber:    user_data.PhoneNumber,
		MultiFAEnabled: reversed_multifa_status,
	}

	if err := h.userRepo.Update(tx, &user_update); err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to change 2FA status")
	}

	return utils.ResponseMessage(c, fiber.StatusOK, "Success Change 2FA Status")
}

func (h *UserHandler) ResetPhoneNumber(c *fiber.Ctx) error {

	// this function is to reset the phone number to "" and return success

	user := c.Locals("user").(models.UserSession)

	tx, err := database.DB.Begin()
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to start transaction")
	}
	defer func() {
		database.CommitOrRollback(tx, c, err)
	}()

	// check user validity
	user_check, err := h.userRepo.FindByID(tx, user.Id)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to check user")
	}

	if user_check.Id == 0 {
		return utils.ResponseError(c, fiber.StatusBadRequest, "User not found")
	}

	// update reset the phone number to "" and automatically set multifa to false
	user_update := models.User{
		Id:             user_check.Id,
		Username:       user_check.Username,
		MultiFAEnabled: false,
		PhoneNumber:    "",
	}

	if err := h.userRepo.Update(tx, &user_update); err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to change phone number")

	}

	return utils.ResponseMessage(c, fiber.StatusOK, "Success Change 2FA Status")
}

func (h *UserHandler) ChangePhoneNumber(c *fiber.Ctx) error {

	user := c.Locals("user").(models.UserSession)

	userInput := new(models.UserChangePhoneInput)
	if err := c.BodyParser(&userInput); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid request")
	}

	if err := utils.ValidateStruct(userInput); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Id":
				return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid request, invalid id")
			case "PhoneNumber":
				return utils.ResponseError(c, fiber.StatusBadRequest, "Phone number must be filled and valid")
			}
		}
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to start transaction")
	}
	defer func() {
		database.CommitOrRollback(tx, c, err)
	}()

	user_data, err := h.userRepo.FindByID(tx, user.Id)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to check user")
	}

	// check user validity by input data
	if user_data.Id == 0 {
		return utils.ResponseError(c, fiber.StatusBadRequest, "User not found")
	}

	if user_data.Id != user.Id {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid request")
	}

	// if phone number input just same as before just return it success
	if user_data.PhoneNumber == userInput.PhoneNumber {
		return utils.ResponseMessage(c, fiber.StatusOK, "Success Change Phone Number")
	}

	// if different so, check if phone number already exist
	user_check_number, err := h.userRepo.FindByPhoneNumber(tx, userInput.PhoneNumber)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to check phone number")
	}

	// if user_check_number is is different and is not 0 so the data is exist, failed the request
	if user_check_number.Id != 0 && user_check_number.Id != user_data.Id {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Phone number already registered on another account")
	}

	// else here, so maybe if user just using "" so just update the phone number to "" and return success
	updateUser := models.User{
		Id:             user.Id,
		Username:       user.Username,
		PhoneNumber:    userInput.PhoneNumber,
		MultiFAEnabled: user_data.MultiFAEnabled,
	}

	if err := h.userRepo.Update(tx, &updateUser); err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to change phone number")
	}

	return utils.ResponseMessage(c, fiber.StatusOK, "Success Change Phone Number")
}

func (h *UserHandler) ChangeUsername(c *fiber.Ctx) error {
	// var txError error
	user := c.Locals("user").(models.UserSession)

	userInput := new(models.UserChangeUsernameInput)
	if err := c.BodyParser(&userInput); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid request")
	}

	if err := utils.ValidateStruct(userInput); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Id":
				return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid request, invalid id")
			case "Username":
				return utils.ResponseError(c, fiber.StatusBadRequest, "Username must be alphanumeric and between 3-25 characters")
			}
		}
	}

	// if new username is same with current username, just return success
	if user.Username == userInput.Username {
		return utils.ResponseMessage(c, fiber.StatusOK, "Success Change Username")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to start transaction")
	}
	defer func() {
		database.CommitOrRollback(tx, c, err)
	}()

	// check if user is exist or not
	userCheck, err := h.userRepo.FindByID(tx, user.Id)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to check user")
	}

	if userCheck.Id == 0 {
		return utils.ResponseError(c, fiber.StatusBadRequest, "User not found")
	}

	// check if id inputted is same with user id
	if userCheck.Id != user.Id {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid request")
	}

	// check if username already exist
	isUsernameExist, err := h.userRepo.FindByUsername(tx, userInput.Username)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to check username")
	}

	if isUsernameExist.Id > 0 {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Username already exist")
	}

	// if free to use, update username
	updateUser := models.User{
		Id:             user.Id,
		Username:       userInput.Username,
		PhoneNumber:    userCheck.PhoneNumber,
		MultiFAEnabled: userCheck.MultiFAEnabled,
	}

	if err := h.userRepo.Update(tx, &updateUser); err != nil {
		// txError = err
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to change username")
	}

	return utils.ResponseMessage(c, fiber.StatusOK, "Success Change Username")
}

func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserSession)

	passInput := new(models.UserChangePasswordInput)
	if err := c.BodyParser(&passInput); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid request")
	}

	if err := utils.ValidateStruct(passInput); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Id":
				return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid request")
			case "Password":
				return utils.ResponseError(c, fiber.StatusBadRequest, "Password must be filled")
			case "NewPassword":
				return utils.ResponseError(c, fiber.StatusBadRequest, "New Password must be alphanumeric and between 6-50 characters, contains number and uppercase")
			}
		}
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to start transaction")
	}
	defer func() {
		database.CommitOrRollback(tx, c, err)
	}()

	// check if id inputted is same with user id
	if user.Id != passInput.Id {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid request")
	}

	// check if user exist
	userCheck, err := h.userRepo.FindByID(tx, user.Id)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to check user")
	}

	if userCheck.Id == 0 {
		return utils.ResponseError(c, fiber.StatusBadRequest, "User not found")
	}

	// check if password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(userCheck.Password), []byte(passInput.Password)); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid current password value")
	}

	// hash password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(passInput.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to hash new password")
	}

	// update password
	if err := h.userRepo.UpdatePassword(tx, &models.User{
		Id:       user.Id,
		Username: userCheck.Username,
		Password: string(hashedPass),
	}); err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, "Failed to change password")
	}

	return utils.ResponseMessage(c, fiber.StatusOK, "Success Change Password")
}
