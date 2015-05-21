'use strict';

$(function() {

	/** Indexing form fields for error handling & form validation **/
	var 		  tos = {obj: $('#signup-form #tos-container')};
	var 		fname = {obj: $('#signup-form #firstname')};
	var 		lname = {obj: $('#signup-form #lastname')};
	var 		email = {obj: $('#signup-form #email')};
	var 	   gender = {obj: $('#signup-form #gender-select')};
	var 	  zipcode = {obj: $('#signup-form #zipcode')};
	var 	 username = {obj: $('#signup-form #username')};
	var      password = {obj: $('#signup-form #password')};
	var    emailAgain = {obj: $('#signup-form #email-again')};
	var passwordAgain = {obj: $('#signup-form #password-again')};

	/** Give focus to first name field on load **/
	$('#firstname').focus();

	/** Display errors from processing **/
	/**
		Code				Error Meaning
		-----------			--------------
		invalid_form		Missing form field
		UIN					Username in use
		F 					Firstname invalid
		L 					Lastname invalid
		U 					Username invalid
		E 					Email invalid
		EM 					Email's don't match
		P 					Password invalid
		PM 					Password's don't match
		G 					Gender is not m || f
		Z 					Zipcode invalid

	*/
	var errors = getParam('err');
	if(errors) {

		/** Fill in the fields **/
		fname.obj.val(getParam('f'));
		lname.obj.val(getParam('l'));
		username.obj.val(getParam('u'));
		email.obj.val(getParam('e'));
		emailAgain.obj.val(getParam('e'));
		zipcode.obj.val(getParam('z'));
		// Set gender
		var gen = document.getElementById('gender-select');
		for(var i = 0; i < gen.options.length; i++) {
			if(gen.options[i].getAttribute('value') == getParam('g')) {
				gen.selectedIndex = i;
				break;
			}
		}


		errors = errors.split('|');
		for(var e in errors) {
			switch(errors[e]) {
				case 'invalid_form':
					createAlert("Form was missing some fields", "high");
					break;
				case 'logged_in':
					createAlert("Must log out before creating a new account", "high");
					break;
				case 'F':
					badStyle(fname.obj);
					generateSignupError('Invalid name', lname.obj);
					goodStyle(lname.obj);
					clearField(fname.obj);
					break;
				case 'L':
					clearErrors();
					generateSignupError('Invalid name', lname.obj);
					clearField(lname.obj);
					break;
				case 'U':
					generateSignupError('Invalid username', username.obj);
					clearField(username.obj);
					break;
				case 'UIN':
					generateSignupError('Username is already in use', username.obj);
					clearField(username.obj);
					break;
				case 'P':
					generateSignupError('Password invalid', password.obj);
					break;
				case 'PM':
					generateSignupError('Password\'s don\'t match', password.obj);
					break;
				case 'E':
					generateSignupError('Email address invalid', email.obj);
					clearField(email.obj);
					clearField(emailAgain.obj);
					break;
				case 'EM':
					generateSignupError('Email addresses don\'t match', email.obj);
					clearField(email.obj);
					clearField(emailAgain.obj);
					break;
				case 'EIN':
					generateSignupError('Email address already in use', email.obj);
					clearField(email.obj);
					clearField(emailAgain.obj);
					break;
				case 'Z':
					generateSignupError("Invalid zipcode", zipcode.obj);
					clearField(zipcode.obj);
					break;
			}
		}
	}

	/** On submit Validate form **/
	$('#signup-form').on('submit', function(e) {
		// Reset field styles
		resetStyles();
		// Remove previous error messages
		clearErrors();
		// Validate form
		if(!validSignup()) {
			e.preventDefault();
		}
	});
});

function clearField(obj) {
	obj.val('');
	return;
}

function resetStyles() {
	$('#signup-form .field').each(function() {
		goodStyle($(this));
	});
}

function clearErrors() {
	$('.form-error').each(function() {
		$(this).remove();
	});
	return;
}

function generateSignupError(string, field) {
	$("<span class='form-error'>" + string + "</span>").insertAfter(field);
	badStyle(field);
}

function validateName(string) {
	return (/^[A-Za-z]+(([\'-])?[A-Za-z]+$)/.test(string)
		&& string.length >= 2 && string.length < 32) ? true : false;
}

function validateZipcode(string) {
	return (/[0-9]/.test(string) && string.length == 5) ? true : false;
}

function validateEmail(string) {
	return (/^([a-zA-Z0-9_.+-])+\@(([a-zA-Z0-9-])+\.)+([a-zA-Z0-9]{2,4})+$/.test(string)
		&& string.length < 64) ? true : false;
}

function validateUsername(string) {
	return (/^[A-Za-z0-9_]+$/.test(string)
		&& string.length >= 2
		&& string.length < 16) ? true : false;
}

function validatePassword(string) {
	return (string.length > 2 && string.length < 32) ? true : false;
}

function validSignup() {
	var status = true;

	var fname = {
		obj: $('#signup-form #firstname'),
		val: $('#signup-form #firstname').val()
	};
	var lname = {
		obj: $('#signup-form #lastname'),
		val: $('#signup-form #lastname').val()
	};
	var username = {
		obj: $('#signup-form #username'),
		val: $('#signup-form #username').val()
	};
	var password = {
		obj: $('#signup-form #password'),
		val: $('#signup-form #password').val()
	};
	var passwordAgain = {
		obj: $('#signup-form #password-again'),
		val: $('#signup-form #password-again').val()
	};
	var email = {
		obj: $('#signup-form #email'),
		val: $('#signup-form #email').val()
	};
	var emailAgain = {
		obj: $('#signup-form #email-again'),
		val: $('#signup-form #email-again').val()
	};
	var zipcode = {
		obj: $('#signup-form #zipcode'),
		val: $('#signup-form #zipcode').val()
	};
	var gender = {
		obj: $('#signup-form #gender-select'),
		val: $('#signup-form #gender-select').val()
	};
	var tos = {
		obj: $('#signup-form #tos-container'),
		val: $('#signup-form #tos-checkbox').is(':checked')
	};

	// Validate names
	if(!validateName(fname.val)) {
		status = false;
		generateSignupError('Invalid name', fname.obj);
	}

	if(!validateName(lname.val)) {
		status = false;
		clearErrors();
		generateSignupError('Invalid name', lname.obj);
	}

	// Validate username
	if(!validateUsername(username.val)) {
		status = false;
		generateSignupError('Invalid username', username.obj);
	}

	// Validate password
	if(validatePassword(password.val)) {
		// If valid password ensure both password fields match
		if(password.val != passwordAgain.val) {
			status = false;
			generateSignupError('Passwords don\'t match', passwordAgain.obj);
		}
	} else {
		status = false;
		generateSignupError('Invalid password', password.obj);
	}

	// Validate email
	if(validateEmail(email.val)) {
		// If valid email ensure both fields match
		if(email.val.toLowerCase() != emailAgain.val.toLowerCase()) {
			status = false;
			generateSignupError('Email addresses don\'t match', emailAgain.obj);
		}
	} else {
		status = false;
		generateSignupError('Invalid email', email.obj);
	}

	// Validate Zipcode
	if(!validateZipcode(zipcode.val)) {
		status = false;
		generateSignupError('Invalid zipcode', zipcode.obj);
	}

	// Ensure a gender was selected
	if(gender.val != 'm' && gender.val != 'f') {
		status = false;
		generateSignupError('Must select a gender', gender.obj);
	}

	// Ensure Terms of service has been agreed to
	if(!tos.val) {
		status = false;
		generateSignupError('Must agree to terms of service', tos.obj);
	}

	return status;
}