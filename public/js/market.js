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
		// Get file from input
		var file = this.files[0];
		// Get handler
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
			useValidator(handlerID, result)
		});
	});

	/** Post-form image validator error detection **/
	$('#post-form-validator').on('error', function() {
		resetImageValidator();
		// Reset handler
		createAlert('Not a valid image', 'medium');
	});

});








/**** INCOMPLETE ****/



/** Returns the photo handler ( displays currently selected image ) to default `cross` image **/
function resetHandler(id) {
	return;
}







/**** INCOMPLETE ****/


function loadImage(file, callback) {
	var reader = new FileReader();

	reader.onerror = function() {
		callback(false);
	};


	reader.onloadend = function() {
		callback(reader.result);
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

function useValidator(handlerID, url) {
	$('#post-form-validator').attr('data-for', handlerID);
	$('#post-form-validator').attr('src', result);
	return;
}

function resetImageValidator() {
	/** Reset handler & input **/
	$('#post-form-validator').attr('src', '/img/cross.png');
	return;
}

function resetImageInput(handlerID) {
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