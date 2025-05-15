hideLoader()

        const COUNTDOWN_RESET = 60

        $('document').ready(function() {
            // Focus the first input on page load
            $('#otp1').focus()
            
            // Setup OTP input functionality
            setupOtpInputs()
            
            // Setup countdown timer
            startCountdown(COUNTDOWN_RESET)

            function getFullOtp() {
                let otp = ''
                for (let i = 1; i <= 6; i++) {
                    otp += $('#otp' + i).val() || ''
                }
                return otp
            }
            
            // Setup form submission
            $("#otpForm").submit(async function(event) {
                event.preventDefault()
                
                // Collect OTP from individual inputs
                const otp = getFullOtp()
                
                if (!otp || otp.length !== 6) {
                    showInfoModal("Invalid Code", "Please enter a complete 6-digit verification code")
                    return
                }
                
                showLoader()
                
                try {
                    const resp = await fetch("/api/multifa/verify", {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({
                            otp: otp
                        }),
                    })
                    const response = await resp.json()
                    
                    if (!response.error) window.location.href = '/'
                    else throw new Error(response.message)
                    
                } catch(e) {
                    showInfoModal('Verification Failed: ' + e.message, 'OTP Verification Failed')
                } finally {
                    hideLoader()
                }
            })
            
            // Setup resend button
            $("#resendBtn").on('click', async function() {
                $(this).prop('disabled', true)
                
                showLoader()
                
                try {
                    const resp = await fetch("/api/multifa/resend", {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        }
                    })
                    const response = await resp.json()
                    
                    if (!response.error) {
                        showInfoModal('A new verification code has been sent to your device', 'Code Sent')
                        startCountdown(COUNTDOWN_RESET)
                        
                        // Clear OTP fields
                        $('.otp-input').val('')
                        $('#otp1').focus()
                    } else {
                        throw new Error(response.message)
                    }
                    
                } catch(e) {
                    showInfoModal('Failed to resend code: ' + e.message, 'Resend Failed')
                    $(this).prop('disabled', false)
                } finally {
                    hideLoader()
                }
            })
            
            // Setup logout button (New functionality)
            $("#logoutBtn").on('click', async function(event) {
                event.preventDefault()
                
                showLoader()
                
                try {
                    const resp = await fetch("/api/logout", {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        }
                    })
                    const response = await resp.json()
                    
                    if (!response.error) {
                        window.location.href = '/login'
                    } else {
                        throw new Error(response.message)
                    }
                } catch(e) {
                    showInfoModal('Logout Failed: ' + e.message, 'Logout Failed')
                } finally {
                    hideLoader()
                }
            })
        })
        
        function setupOtpInputs() {
            // Auto-focus next input and only allow numbers
            $('input[id^="otp"]').on('input', function(e) {
                const val = $(this).val()
                
                // Allow only numbers
                if (!/^\d*$/.test(val)) {
                    $(this).val('')
                    return
                }
                
                // Auto focus next input
                if (val && $(this).attr('id') !== 'otp6') {
                    const nextInputId = 'otp' + (parseInt($(this).attr('id').replace('otp', '')) + 1)
                    $('#' + nextInputId).focus()
                }
                
                // Update hidden input with full OTP
                $('#fullOtp').val(getFullOtp())
            })
            
            // Handle backspace to go to previous input
            $('input[id^="otp"]').on('keydown', function(e) {
                if (e.key === 'Backspace' && !$(this).val()) {
                    const inputNum = parseInt($(this).attr('id').replace('otp', ''))
                    if (inputNum > 1) {
                        const prevInputId = 'otp' + (inputNum - 1)
                        $('#' + prevInputId).focus().val('')
                    }
                }
            })
            
            // Handle paste event for the entire OTP
            $('input[id^="otp"]').on('paste', function(e) {
                e.preventDefault()
                const clipboardData = e.originalEvent.clipboardData || window.clipboardData
                const pastedData = clipboardData.getData('text')
                
                // If pasted content is a 6-digit number
                if (/^\d{6}$/.test(pastedData)) {
                    for (let i = 0; i < 6; i++) {
                        $('#otp' + (i + 1)).val(pastedData[i])
                    }
                    $('#fullOtp').val(pastedData)
                }
            })
        }
        
        function startCountdown(seconds) {
            let remainingSeconds = seconds
            
            // Reset button state
            $('#resendBtn').prop('disabled', true)
            
            // Update countdown display
            function updateCountdown() {
                const minutes = Math.floor(remainingSeconds / 60)
                const seconds = remainingSeconds % 60
                $('#countdown').text(
                    (minutes < 10 ? '0' + minutes : minutes) + ':' + 
                    (seconds < 10 ? '0' + seconds : seconds)
                )
            }
            
            // Initial display
            updateCountdown()
            
            // Clear any existing interval
            if (window.countdownInterval) {
                clearInterval(window.countdownInterval)
            }
            
            // Start the countdown
            window.countdownInterval = setInterval(function() {
                remainingSeconds--
                updateCountdown()
                
                if (remainingSeconds <= 0) {
                    clearInterval(window.countdownInterval)
                    $('#resendBtn').prop('disabled', false)
                    $('#countdown').text('00:00')
                }
            }, 1000)
        }