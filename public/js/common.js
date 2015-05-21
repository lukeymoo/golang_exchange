'use strict';

$(function() {
	var isLoginSubmitted = false;

	var password = {
		obj: $('#header-login-form #password'),
		val: $('#header-login-form #password').val()
	};

	/** Toggle header menu on click **/
	$(document).on('click', '#header-menu-button', function() {
		toggleHeaderMenu();
	});

	/** Toggle login form on click **/
	$(document).on('click', '#header-login-button', function() {
		toggleHeaderLoginForm();
	});

	/** Close alert on click **/
	$(document).on('click', '.close-notification-button', function(e) {
		removeAlert($(this).parent('.notification').attr('data-id'));
	});

	/** Close forms on escape key **/
	$(document).on('keydown', function(e) {
		// Close form on escape
		if(e.which == 27) {
			if(isHeaderMenuOpen()) {
				closeHeaderMenu();
			}
			if(isHeaderLoginFormOpen()) {
				closeHeaderLoginForm();
			}
		}
	});

	/** Close forms on click outside element **/
	$(document).on('click', function(e) {
		if(isHeaderMenuOpen()) {
			if(!$(e.target).is($('#header-menu-button'))) {
				closeHeaderMenu();
			}
		}
		if(isHeaderLoginFormOpen()) {
			if(!$(e.target).is($('#header-login-button'))
				&& !$(e.target).closest($('#header-login-container')).length) {
				closeHeaderLoginForm();
			}
		}
	});

	/** Header Login Form **/
	$('#header-login-form input').on('keydown', function(e) {
		if(e.which == 13) {
			$('#header-login-form-button').click();
		}
	});
	$('#header-login-form-button').on('click', function(e) {
		// Reset form
		resetLoginStyles();
		clearLoginErrors();
		if(!isLoginSubmitted) {
			// If valid login form
			if(validLoginForm()) {
				isLoginSubmitted = true;
				// Hide submit button
				$('#header-login-form-button').hide();
				// Submit
				tryLogin(function(res) {
					// On good login refresh page
					if(res.Status == 'DX-OK') {
						window.location.href = '/';
					} else {
						isLoginSubmitted = false;
						// Show login button
						$('#header-login-form-button').show();
						// Insert error after password field
						generateFormError(res.Message, password.obj);
					}
				});
			}
		}
	});

});

function tryLogin(callback) {
	$.ajax({
		type: 'POST',
		url: '/login/process',
		data: {
			u: $('#header-login-form #username-or-email').val(),
			p: $('#header-login-form #password').val()
		},
		error: function(err) {
			var res = {
				status: 'DX-FAILED',
				message: 'Server error'
			};
			if(err.status == 0) {
				res.Message = 'Server is currently down';
			}
			if(err.status == 404) {
				res.Message = 'Something is wrong on server';
			}
			callback(res);
		}
	}).done(function(res) {
		callback(res);
	});
}

function toggleHeaderMenu() {
	if(isHeaderMenuOpen()) {
		closeHeaderMenu();
	} else {
		openHeaderMenu();
	}
	return;
}

function openHeaderMenu() {
	$('#header-menu-button').attr('data-state', 'opened');
	$('#header-menu').show();
	return;
}

function closeHeaderMenu() {
	$('#header-menu-button').attr('data-state', 'closed');
	$('#header-menu').hide();
	return;
}

function isHeaderMenuOpen() {
	return ($('#header-menu-button').attr('data-state') == 'opened') ? true : false;
}

function toggleHeaderLoginForm() {
	if(isHeaderLoginFormOpen()) {
		closeHeaderLoginForm();
	} else {
		openHeaderLoginForm();
	}
	return;
}

function isHeaderLoginFormOpen() {
	return ($('#header-login-button').attr('data-state') == 'opened') ? true : false;
}

function openHeaderLoginForm() {
	$('#header-login-button').attr('data-state', 'opened');
	$('#header-login-container').show();
	$('#header-login-container input').first().focus();
	return;
}

function closeHeaderLoginForm() {
	$('#header-login-button').attr('data-state', 'closed');
	$('#header-login-container').hide();
}

function getParam(sParam) {
	var sPageURL = window.location.search.substring(1);
	var sURLVariables = sPageURL.split('&');
	for(var i = 0; i < sURLVariables.length; i++) {
		var sParameterName = sURLVariables[i].split('=');
		if (sParameterName[0] == sParam) {
			return sParameterName[1];
		}
	}
}

function goodStyle(obj) {
	$(obj).css('border', '1px solid rgba(0, 0, 0, 0.15)');
}

function badStyle(obj) {
	$(obj).css('border', '1px solid rgb(175, 0, 0)');
	return;
}

function validLoginForm() {
	var login_username = {obj:$('#header-login-form #username-or-email'), val:$('#header-login-form #username-or-email').val()};
	var login_password = {obj:$('#header-login-form #password'), val:$('#header-login-form #password').val()}

	var status = true;

	if(login_password.val.length) {
		if(!validPassword(login_password.val)) {
			status = false;
			generateFormError('Invalid password', login_password.obj);
		}
	} else {
		status = false;
		generateFormError('Password field is empty', login_password.obj);
	}

	if(login_username.val.length) {
		if(!validUsername(login_username.val)) {
			if(!validEmail(login_username.val)) {
				status = false;
				generateFormError('Invalid username or email', login_password.obj);
			}
		}
	} else {
		status = false;
		generateFormError('Username field is empty', login_password.obj);
	}

	return status;
}

function clearLoginErrors() {
	$('#header-login-form .form-error').each(function() {
		$(this).remove();
	});
	return;
}

function resetLoginStyles() {
	$('#header-login-form input').each(function() {
		goodStyle($(this));
	});
	return;
}

function generateFormError(string, obj) {
	$("<span class='form-error'>" + string + "</span>").insertAfter(obj);
	return;
}

function validUsername(string) {
	return (/^[A-Za-z0-9_]+$/.test(string)
		&& string.length >= 2
		&& string.length < 16) ? true : false;
}

function validEmail(string) {
	return (/^([a-zA-Z0-9_.+-])+\@(([a-zA-Z0-9-])+\.)+([a-zA-Z0-9]{2,4})+$/.test(string)
		&& string.length < 64) ? true : false;
}

function validPassword(string) {
	return (string.length > 2 && string.length < 32) ? true : false;
}

function createAlert(string, alertLevel) {
	alertLevel = alertLevel || 'low';
	var classLevel = '';
	if(alertLevel == 'high') {
		classLevel = 'high-alert-level';
	}
	if(alertLevel == 'medium') {
		classLevel = 'medium-alert-level';
	}
	if(alertLevel == 'low') {
		classLevel = 'low-alert-level';
	}
	// Gender an ID using Date.now() ( this is for setTimeout removal )
	var noticeID = Date.now();
	var DOM = 
	"<div data-id='" + noticeID + "' class='notification " + classLevel +  "'>" + 
		"<span class='close-notification-button'>&times;</span>" +
		"<span class='notification-text'>" + string +  "</span>" +
	"</div>";
	$('#notification-container').prepend(DOM);
	// Remove this notice after 6.5 seconds
	setTimeout(function() {
		$('#notification-container .notification').each(function() {
			if($(this).attr('data-id') == noticeID) {
				$(this).fadeOut('slow', function() {
					$(this).remove();
					return false;
				});
			}
		});
	}, 6500);
	return;
}

function removeAlert(id) {
	$('#notification-container .notification').each(function() {
		if($(this).attr('data-id') == id) {
			$(this).fadeOut('fast', function() {
				$(this).remove();
				return false;
			});
		}
	});
	return;
}