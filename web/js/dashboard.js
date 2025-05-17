const editPassModal = new bootstrap.Modal($('#editPasswordModal'))
const editUsernameModal = new bootstrap.Modal($('#editUsernameModal'))
const toggle2FAConfirmModal = new bootstrap.Modal($('#toggle2FAConfirmModal'))
const editPhoneModal = new bootstrap.Modal($('#editPhoneModal'))
const resetPhoneModal = new bootstrap.Modal($('#resetPhoneModal'))
const phoneOtpVerificationModal = new bootstrap.Modal($('#phoneOtpVerificationModal'))
let USERNAME = ""
let USER_ID = 0
let IS_LOGGED_IN = false
let CREDIT_TOKEN = 0
let LAST_FIRST_LLM_USED = ''
let PHONE_NUMBER = ''
let STATUS2FA = false
let USING2FA = false
let TEMP_NUMBER_CHANGE = ''

// OTP input handling
function setupOtpInputs() {
    const inputs = document.querySelectorAll('#phoneOtpVerificationModal input[type="text"]');
    const hiddenInput = document.getElementById('fullOtpInput');
    
    inputs.forEach((input, index) => {
        // Auto-focus to next input when a digit is entered
        input.addEventListener('input', function(e) {
            const value = e.target.value;
            
            if (value.length === 1) {
                if (index < inputs.length - 1) {
                    inputs[index + 1].focus();
                }
                
                // Update the hidden input with all values
                hiddenInput.value = Array.from(inputs).map(input => input.value).join('');
            }
        });
        
        // Handle backspace key
        input.addEventListener('keydown', function(e) {
            if (e.key === 'Backspace' && !e.target.value) {
                if (index > 0) {
                    inputs[index - 1].focus();
                }
            }
        });
        
        // Handle paste event
        input.addEventListener('paste', function(e) {
            e.preventDefault();
            const paste = (e.clipboardData || window.clipboardData).getData('text');
            
            if (paste.match(/^\d+$/) && paste.length === inputs.length) {
                // If pasted content is a numeric string with correct length
                for (let i = 0; i < inputs.length; i++) {
                    inputs[i].value = paste[i];
                }
                hiddenInput.value = paste;
                inputs[inputs.length - 1].focus();
            }
        });
    });
}

async function redirectAuth(app_name) {
    event.preventDefault()
    showLoader()

    try {
        if (!IS_LOGGED_IN) {
            throw new Error('You need to log in first to access this application.')
        }

        const resp = await fetch("/api/redirect?app=" + app_name, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
        })
        const response = await resp.json()

        if (response.error) {
            throw new Error(response.message)
        } else {
            const url = response.data.redirect_url + '?token=' + response.data.token
            window.location.href = url
        }

    } catch (e) {
        showInfoModal(app_name + ' Failed: <b>' + e.message + "</b>", app_name + ' Failed')
    } finally {
        hideLoader()
    }
}

async function loadDataDashboard() {
    showLoader()

    try {
        const resp = await fetch("/api/dashboard", {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
        })
        const response = await resp.json()

        if (response.error) {
            throw new Error(response.message)
        } else {
            const data = response.data
            IS_LOGGED_IN = data.is_logged_in

            // check new status for 2mfa check
            USING2FA = data.using_mfa

            // check if using 2fa and logged in true, so throw to 2fa page
            if (USING2FA && !IS_LOGGED_IN) {
                window.location.href = '/multifa'
            }

            if (IS_LOGGED_IN) {
                // Show logged-in content and set username
                USERNAME = data.user.username
                USER_ID = data.user.id
                CREDIT_TOKEN = data.user.credit_token
                LAST_FIRST_LLM_USED = data.user.last_first_llm_used
                PHONE_NUMBER = data.user.phone_number
                STATUS2FA = data.user.multifa_enabled
                $('#logged-in').removeClass('hidden')
                $('#not-logged-in').addClass('hidden')
                $('#username').text(USERNAME)
                $('#credit-token').text(CREDIT_TOKEN)
                updatePhoneNumberUI()
                $('#status2FA').text(STATUS2FA ? 'On' : 'Off')

                if (LAST_FIRST_LLM_USED === '') {
                    $('#last-llm-used').text('User has max token daily and has not used any LLM feature today');
                } else {

                    let nextDay = new Date(LAST_FIRST_LLM_USED)
                    nextDay.setHours(nextDay.getHours() + 24)

                    let now = new Date()

                    let diff = nextDay - now

                    let hours = Math.floor(diff / (1000 * 60 * 60))
                    let minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
                    let seconds = Math.floor((diff % (1000 * 60)) / 1000)

                    let time_until_reset_token = `${hours} hours, ${minutes} minutes, ${seconds} seconds`

                    $('#last-llm-used').text(time_until_reset_token);
                }

                $('#usernameEditCurrentInput').val(USERNAME)
                $('#usernameEditInput').val(USERNAME)

            } else {
                // Show not-logged-in content
                $('#logged-in').addClass('hidden');
                $('#not-logged-in').removeClass('hidden');
            }
        }

    } catch (e) {
        showInfoModal('Failed to load data: ' + e.message, 'Failed to Load Data')
    } finally {
        hideLoader()
    }
}

// Update UI based on phone number
function updatePhoneNumberUI() {
    if (PHONE_NUMBER === '' || PHONE_NUMBER === null) {
        // Hide reset phone button if no phone number is set
        $('#resetPhoneBtn').hide();
        $('#phoneNumber').text('You have not set your phone number yet');
    } else {
        // Show reset phone button if phone number exists
        $('#resetPhoneBtn').show();
        $('#phoneNumber').text("+62 " + PHONE_NUMBER);
    }
}

// Initialize the 2FA toggle state based on current status
async function update2FAToggleUI() {
    if (STATUS2FA) {
        $('#toggle2FACheckbox').prop('checked', true);
        $('.slider').addClass('bg-purple-500');
        $('.slider').removeClass('bg-gray-300');
        $('.slider').addClass('before:translate-x-5');
    } else {
        $('#toggle2FACheckbox').prop('checked', false);
        $('.slider').removeClass('bg-purple-500');
        $('.slider').addClass('bg-gray-300');
        $('.slider').removeClass('before:translate-x-5');
    }
    $('#status2FA').text(STATUS2FA ? 'On' : 'Off');
}

$('document').ready(async function () {
    hideLoader()
    await loadDataDashboard()
    // Initialize the toggle state on page load
    await update2FAToggleUI()

    // reset phone number modal
    $('#resetPhoneBtn').click(function (e) {
        e.preventDefault()

        resetPhoneModal.show()
    })

    // Setup the confirmation modal before showing it
    $('#toggle2FA').click(function (e) {
        // Prevent default action if using data-bs-toggle
        e.preventDefault()

        // first check here if the phone number is no set, user can'nt set 2FA
        if (PHONE_NUMBER === '' || PHONE_NUMBER === null) {
            showInfoModal('You need to set your phone number first to enable 2FA', 'Enable 2FA Failed')
            return
        }

        const newStatus = !STATUS2FA

        // Update confirmation modal content based on current state
        if (newStatus) {
            $('#toggle2FAMessage').text('Are you sure you want to enable 2FA?')
            $('#toggle2FADescription').text('Two-factor authentication adds an extra layer of security to your account.')

        } else {
            $('#toggle2FAMessage').text('Are you sure you want to disable 2FA?')
            $('#toggle2FADescription').text('Disabling two-factor authentication will reduce the security of your account.')
        }

        toggle2FAConfirmModal.show()

    })

    // Handle the confirmation button click
    $('#confirm2FAToggle').click(async function () {
        const newStatus = !STATUS2FA
        showLoader()

        try {
            const resp = await fetch("/api/users/2fa", {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    id: parseInt(USER_ID),
                    multifa_enabled: newStatus
                })
            })

            const response = await resp.json()

            if (response.error) {
                throw new Error(response.message)
            } else {
                // Update the status and UI
                STATUS2FA = newStatus;
                update2FAToggleUI();

                toggle2FAConfirmModal.hide()
                showInfoModal(
                    `Two-factor authentication has been ${newStatus ? 'enabled' : 'disabled'} successfully.`,
                    '2FA Status Updated'
                )
            }
        } catch (e) {
            showInfoModal('Failed to update 2FA status: ' + e.message, '2FA Update Failed')
        } finally {
            hideLoader()
        }
    })

    $('#echonotes').click(async function () {
        await redirectAuth('echonotes')
    })

    $('#gochat').click(async function () {
        await redirectAuth('gochat')
    })

    $('#llm').click(async function () {
        await redirectAuth('llm')
    })


    // --------------------------------------------- EDIT USERNAME, PHONE NUMBER AND PASSWORD
    // edit username
    $('#editUsernameForm').submit(async function () {
        event.preventDefault()

        const username = $('#usernameEditInput').val().trim()

        showLoader()

        try {
            const resp = await fetch("/api/users", {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    id: parseInt(USER_ID),
                    username: username
                })
            })
            const response = await resp.json()

            editUsernameModal.hide()

            if (response.error) throw new Error(response.message)
            else {
                hideLoader()
                showInfoModal('Success edit username', 'Edit Username Success')
                setTimeout(() => {
                    window.location.reload()
                }, 1000)
            }

        } catch (e) {
            showInfoModal('Failed to edit username: ' + e.message, 'Edit Username Failed')
            hideLoader()
        }
    })

    async function sendOTPEditPhoneNumber(phoneNumber) {
            const response = await fetch("/api/users/phone/otp", {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    phone_number: phoneNumber
                })
            })
            const resp = await response.json()

            return resp
    }

    async function verifyOTPAndEditPhoneNumber(otp_code, phoneNumber) {
        const response = await fetch("/api/users/phone", {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                otp_code: otp_code,
                phone_number: phoneNumber
            })
        })
        const resp = await response.json()

        return resp
    }

    $('#editPhoneForm').submit(async function () {
        event.preventDefault()

        const phoneNumber = $('#phoneNumberEditInput').val().trim()
        TEMP_NUMBER_CHANGE = phoneNumber

        // check phone number for now must start with 8
        if (phoneNumber[0] !== '8') {
            showInfoModal('Phone number must start with 8 and at least 8 digits', 'Edit Phone Number Failed')
            return
        }

        if (phoneNumber === PHONE_NUMBER) {
            showInfoModal('Phone number is the same as before', 'Edit Phone Number Success')
            return
        }

        showLoader()

        try {
            const sendOTPNum = await sendOTPEditPhoneNumber(phoneNumber)
            if (sendOTPNum.error) throw new Error(sendOTPNum.message)
            
            // Hide the phone edit modal and show OTP verification modal
            editPhoneModal.hide()
            
            // Reset OTP inputs before showing modal
            const otpInputs = document.querySelectorAll('#phoneOtpVerificationModal input[type="text"]');
            otpInputs.forEach(input => input.value = '');
            document.getElementById('fullOtpInput').value = '';
            
            // Setup OTP input behavior
            setupOtpInputs();
            
            // Show the OTP verification modal
            phoneOtpVerificationModal.show();
            
            // Focus on the first input
            otpInputs[0].focus();

            // reset the phone number input
            $('#phoneNumberEditInput').val('')
            
            hideLoader()

        } catch (e) {
            showInfoModal('Failed to send OTP: ' + e.message, 'Edit Phone Number Failed')
            hideLoader()
        }
    })

    // Verify phone OTP handler
    $('#verifyPhoneOtpBtn').click(async function() {
        const otpCode = document.getElementById('fullOtpInput').value
        
        if (otpCode.length !== 6 || !/^\d+$/.test(otpCode)) {
            showInfoModal('Please enter a valid 6-digit OTP code', 'Verification Failed')
            return
        }
        
        showLoader()
        
        try {
            // check the otp verification and edit the phone number
            const updateResp = await verifyOTPAndEditPhoneNumber(otpCode, TEMP_NUMBER_CHANGE);
            
            if (updateResp.error) throw new Error(updateResp.message)

            // Close the OTP modal
            phoneOtpVerificationModal.hide()
            
            // Show success message
            showInfoModal('Phone number updated successfully!', 'Success')
            
            // Reload page after a short delay
            setTimeout(() => {
                window.location.reload()
            }, 1500)
            
        } catch (e) {
            showInfoModal('Verification failed: ' + e.message, 'Verification Failed')
        } finally {
            hideLoader()
        }
    })

    // edit password
    $('#editPasswordForm').submit(async function () {
        event.preventDefault()

        const passwordNow = $('#passwordEditNowInput').val().trim()
        const newPassword = $('#passwordEditNewInput').val().trim()
        const confirmPassword = $('#passwordEditConfirmInput').val().trim()

        if (newPassword !== confirmPassword) {
            showInfoModal('Password and password confirmation not same', 'Edit Password Failed')
            return
        }

        showLoader()

        try {
            const resp = await fetch("/api/users/password", {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    id: parseInt(USER_ID),
                    password: passwordNow,
                    new_password: newPassword
                })
            })
            const response = await resp.json()

            if (response.error) throw new Error(response.message)
            else {
                hideLoader()
                editPassModal.hide()
                $('#passwordEditNowInput').val('')
                $('#passwordEditNewInput').val('')
                $('#passwordEditConfirmInput').val('')
                showInfoModal('Success edit password', 'Edit Password Success')
            }

        } catch (e) {
            showInfoModal('Failed to edit password: ' + e.message, 'Edit Password Failed')
            hideLoader()
        }
    })

    // logout
    $('#logout').click(async function () {
        event.preventDefault()

        showLoader()

        try {
            const resp = await fetch("/api/logout", {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
            })
            const response = await resp.json()

            if (!response.error) {
                window.location.href = '/login'
            } else {
                showInfoModal('Logout Failed', 'Logout Failed')
            }
        } catch (error) {
            showInfoModal('Logout Failed', 'Logout Failed')
        } finally {
            hideLoader()
        }
    })

    // reset phone number modal
    $('#confirmResetPhone').click(async function () {
        showLoader()

        try {
            const resp = await fetch("/api/users/phone/reset", {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json'
                }
            })

            const response = await resp.json()

            if (response.error) {
                throw new Error(response.message)
            } else {
                // Close the modal
                resetPhoneModal.hide()

                // Show success message
                showInfoModal(
                    'Your phone number has been reset successfully. This has also disabled 2FA if it was enabled.',
                    'Phone Number Reset Success'
                )

                // Reload the page after a short delay to update the UI
                setTimeout(() => {
                    window.location.reload()
                }, 1500)
            }
        } catch (e) {
            showInfoModal('Failed to reset phone number: ' + e.message, 'Reset Phone Number Failed')
        } finally {
            hideLoader()
        }
    })
})