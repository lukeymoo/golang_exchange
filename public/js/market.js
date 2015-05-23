'use strict';

$(function() {

	/** First post-action-menu coloring on hover **/
	$(document).on({
		mouseenter: function() {
			colorPostActionMenuItem($(this));
		},
		mouseleave: function() {
			uncolorPostActionMenuItem($(this));
		}
	}, '.post-action-menu-list li:first-child');

	/** Close action menus on external click **/
	$(document).on('click', function(e) {
		if(!$(e.target).is($('.post-action-menu-button'))) {
			closeAllPostActionMenus();
		} else {
			closeAllPostActionMenusExcept($(e.target));
		}
	});

	/** Post description enter key **/
	$(document).on('keydown', '.post-form-description[contenteditable="true"]', function(e) {
		if(e.which == 13) {
			document.execCommand('insertHTML', false, '<br><br>');
			return false;
		}
	});

	/** Selection of images **/
	$(document).on('click', '.photo-handler', function() {
		$('#' + $(this).attr('data-for')).click();
	});

	/** Validating selection **/
	$(document).on('change', '.post-image', function() {
		// Get file from hidden input field
		var file = this.files[0];
		// Get handler name & prepend with # to create JQuery ID selector
		var handlerID = '#' + $(this).attr('data-for');
		// If selection was cancelled, reset handler
		if(!file) {
			// Grab data-for & place `#` before it
			resetHandler(handlerID);
			return;
		}
		// if bad extension reset handler & input
		if(!validImageExt(file.name)) {
			createAlert('Not a valid image extension `.' + getExtension(file.name) + '`', 'medium');
			return;
		}

		/** Validate the image **/
		loadImage(file, function(result) {
			if(!result) {
				createAlert('An error occurred', 'medium');
				// Reset handler
				resetHandler(handlerID);
				return false;
			}
			if(result == 'invalid') {
				createAlert('Invalid file selected', 'medium');
				resetHandler(handlerID);
				return false;
			}
			if(result == 'small') {
				createAlert('Image must be at least 100x100', 'medium');
				resetHandler(handlerID);
				return false;
			}
			$(handlerID).attr('src', result);
			showNextImageHandler();
		});
	});

});




















function showNextImageHandler() {
	// Iterates through all photo-handlers and ensures they're placed in
	// chronological order
	var i = 1;
	var spawned = false;
	$('.photo-handler').each(function() {
		// If inactive & spawned is false show it
		if($(this).attr('data-active') == 'false' && !spawned) {
			$(this).attr('data-active', 'true');
			spawned = true;
		}
		// Ensure the handlers are in order
		if(i != parseInt($(this).attr('id').split('photo-handler')[1])) {
			// If they're not equal, its out of order..Remove it and insert where correct
			var placement = parseInt($(this).attr('id').split('photo-handler')[1]);
			var handlerDOM = $(this)[0].outerHTML;
			$(this).remove();
			switch(placement) {
				case 1:
					// Insert before photo2
					$(handlerDOM).insertBefore($('#photo-handler2'));
					break;
				case 2:
					// Insert after photo1
					$(handlerDOM).insertAfter($('#photo-handler1'));
					break;
				case 3:
					// Insert after photo2
					$(handlerDOM).insertAfter($('#photo-handler2'));
					break;
				case 4:
					// Insert after photo3
					$(handlerDOM).insertAfter($('#photo-handler3'));
					break;
			}
		}
		i++;
	});
	return;
}

function organizeHandlers() {
	// Iterate through
	return;
}

/** Loads an image and ensures it can load **/
function loadImage(file, callback) {
	var reader = new FileReader();
	reader.onerror = function() {
		callback(null);
	};
	reader.onloadend = function() {

		// Create image
		var pre = new Image();
		if(reader.result == 'data:') {
			callback('invalid');
			return;
		}
		pre.onerror = function() {
			callback(null);
			return;
		};
		pre.onload = function() {

			// Ensure image width x height at minimum is 100x100
			if(parseInt(pre.width) < 100 || parseInt(pre.height) < 100) {
				callback('small');
				return;
			}

			callback(reader.result);
			return;
		};
		pre.src = reader.result;
		return;
	};
	reader.readAsDataURL(file);
	return;
}

/** Returns the extension as a string **/
function getExtension(filename) {
	var parts = filename.split('.');
	return parts[parts.length - 1].toLowerCase();
}

/** Determines if the extension is a valid image type **/
function validImageExt(filename) {
	var parts = filename.split('.');
	var ext = parts[parts.length - 1].toLowerCase();
	return (ext == 'jpg' || ext == 'png'
			|| ext == 'bmp' || ext == 'jpeg'
			|| ext == 'gif' || ext == 'tiff') ? true : false;
}

function resetImageInput(inputNoHash) {
	var inputDOM = $('#' + inputNoHash)[0].outerHTML;
	var placement = parseInt($('#' + inputNoHash).attr('data-for').split('photo-handler')[1]);
	$('#' + inputNoHash).remove();
	switch(placement) {
		case 1:
			// Insert before photo2
			$(inputDOM).insertBefore($('#photo2'));
			break;
		case 2:
			// Insert after photo1
			$(inputDOM).insertAfter($('#photo1'));
			break;
		case 3:
			// Insert after photo2
			$(inputDOM).insertAfter($('#photo2'));
			break;
		case 4:
			// Insert after photo3
			$(inputDOM).insertAfter($('#photo3'));
			break;
	}
	return;
}

/** Returns the photo handler ( displays currently selected image ) to default `cross` image **/
function resetHandler(id) {
	$(id).attr('src', '/img/cross.png');
	resetImageInput($(id).attr('data-for'));
	return;
}






/**** INCOMPLETE ****/


/** Displays post form to upload product **/
function createPostForm() {
	var DOM = "";
	return;
}








/** Closes all post action menus except the specified one **/
function closeAllPostActionMenusExcept(button) {
	$('.post-action-menu-button').each(function() {
		if($(this).parents('.feed-post').find('.info-post-id').html() != $(button).parents('.feed-post').find('.info-post-id').html()) {
			closePostActionMenu($(this));
		}
	});
	return;
}

/** Iterates all posts and closes open action menus **/
function closeAllPostActionMenus() {
	$('.post-action-menu-list').each(function() {
		if($(this).attr('data-state') == 'opened') {
			closePostActionMenu($(this).parent().find('.post-action-menu-list'));
		}
	});
	return;
}

/** Opens/Closes post action menu **/
function togglePostActionMenu(button) {
	if($(button).parents('.feed-post').find('.post-action-menu-list').attr('data-state') == 'opened') {
		console.log('closing');
		closePostActionMenu(button);
	} else if($(button).parents('.feed-post').find('.post-action-menu-list').attr('data-state') == 'closed') {
		openPostActionMenu(button);
	}
	return;
}

/** Opens a specified post action menu **/
function openPostActionMenu(button) {
	$(button).parents('.feed-post').find('.post-action-menu-list').attr('data-state', 'opened');
	$(button).parents('.feed-post').find('.post-action-menu-list').show();
	return;
}

/** Closes a specified post action menu **/
function closePostActionMenu(button) {
	$(button).parents('.feed-post').find('.post-action-menu-list').attr('data-state', 'closed');
	$(button).parents('.feed-post').find('.post-action-menu-list').hide();
	return;
}

/** Used to color the speech bubble end of a post action menu **/
function colorPostActionMenuItem(item) {
	if(!$(item).hasClass('hover')){
		$(item).addClass('hover');
	}
	return;
}

/** Used to return speech bubble end of a post action menu to default color **/
function uncolorPostActionMenuItem(item) {
	if($(item).hasClass('hover')) {
		$(item).removeClass('hover');
	}
	return;
}