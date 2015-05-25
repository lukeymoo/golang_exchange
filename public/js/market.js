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
	$(document).on('keyup keydown', '#post-form-description', function(e) {
		/**
			Update form validity, this will enable or disable form button
		*/
		updatePostFormStatus();
	});

	/** Selection of images **/
	$(document).on('click', '.photo-handler', function() {
		$($(this).attr('data-for')).click();
	});

	/** Image removal on click **/
	$(document).on('click', '.photo-remover', function() {
		// RestoreInput
		restoreInput($(this).attr('data-for'));
		/**
			Update form validity, this will enable or disable form button
		*/
		updatePostFormStatus();
	});

	/** Switch post type **/
	$(document).on('click', '.post-type', function() {
		changePostType($(this));
		updatePostFormStatus();
	});

	/** Submit create form **/
	$(document).on('click', '#post-form-button', function() {
		// Ensure the form is valid before accepting submission
		if($('#post-form').attr('data-valid') == 'true') {
			alert('done');
		}
	});

	/** Validating selection **/
	$(document).on('change', '.post-image', function() {
		var file = this.files[0]; // Get file
		var inputID = '#' + $(this).attr('id'); // input field id
		var handlerContainer = $(this).attr('data-for');

		// If user cancelled selection, restore handler ( input field is already empty )
		if(!file) {
			// Redraw handlers
			restoreHandler(handlerContainer);
			/**
				Update form validity, this will enable or disable form button
			*/
			updatePostFormStatus();
			return;
		}

		/**
			If invalid extension, restore input field & restore handler
		*/
		if(!validImageExt(getExtension(file.name))) {
			createAlert('Not a valid image extension `.' + getExtension(file.name) + '`', 'medium');
			restoreInput(inputID);
			/**
				Update form validity, this will enable or disable form button
			*/
			updatePostFormStatus();
			return;
		}

		/**
			try validating image
			all errors => restore input ( calls restore handler automatically )
		*/
		loadImage(file, function(err, result) {
			if(err) {
				switch(err) {
					case 'empty': // result == data: <= attempted empty file upload
						createAlert('Select a valid image', 'medium');
						restoreInput(inputID);
						break;
					case 'failed': // something wrong happened
						createAlert('An error occurred, try again', 'medium');
						restoreInput(inputID);
						break;
					case 'invalid_image': // Image cannot be displayed = invalid
						createAlert('Invalid file selected', 'medium');
						restoreInput(inputID);
						break;
					case 'dim_small': // Image dimensions are too small
						createAlert('Image must be at least 100x100 in dimensions', 'medium');
						restoreInput(inputID);
						break;
					default:
						break;
				}
			} else {
				setHandler(handlerContainer, result);
			}
			/**
				Update form validity, this will enable or disable form button
			*/
			updatePostFormStatus();
		});

	});

});



/**
	validate extension
	validate image contents data ( can it be displayed? )
	validate image dimensions ( must be at least 100x100 )
	validate image size ( not more than 5MB ) <= might take time to upload for some users
*/



function changePostType(button) {
	// Deselect all post-types
	$('.post-type').each(function() {$(this).attr('data-selected', 'false');});
	// change post-form type
	$('#post-form').attr('data-type', $(button).attr('data-value'));
	// select button
	$(button).attr('data-selected', 'true');
	// change label
	if($('#post-form').attr('data-type') == 'sale') {
		$('#post-image-label').html('Sales <span style="color:rgb(190, 0, 0);">must</span> have at least 1 image');
	}
	if($('#post-form').attr('data-type') == 'general') {
		$('#post-image-label').html('Any images?');
	}
	return;
}



/**
	Update form status,
	enabling or disabling form post button
	
	ensure description has at least 5 characters

	If post type == sale {
		check if at least 1 image was selected
	}
*/
function updatePostFormStatus() {
	var validForm = true;
	
	// de-activate post button
	if(!validDescription($('#post-form-description').val())) {
		$('#post-form').attr('data-valid', 'false');
		return;
	}

	switch($('#post-form').attr('data-type')) {
		case 'sale':
			// ensure at least 1 file is selected
			var count = 0;
			$('.post-image').each(function() {
				if($(this)[0].files.length) {
					count++;
				}
			});
			if(!count) {
				$('#post-form').attr('data-valid', 'false');
				return;
			}
			break;
		case 'general':
			break;
		default:
			break;
	}
	$('#post-form').attr('data-valid', 'true');
	return;
}

function validDescription(string) {
	return (string.length >= 4 && string.length <= 2500) ? true : false;
}

/**
	set active='true', update child photo-handler src
	Call drawHandlers
*/

function setHandler(handlerContainerID, dataURL) {
	$(handlerContainerID).find('.photo-handler').attr('src', dataURL);
	$(handlerContainerID).attr('data-active', 'true');
	drawHandlers();
	return;
}

/**
	Loads image and validates contents
*/
function loadImage(file, callback) {
	var fileReader = new FileReader();

	fileReader.onerror = function() {
		callback('failed', null);
		return;
	};

	fileReader.onloadend = function() {

		// is empty ?
		if(fileReader.result == 'data:') {
			callback('empty', null);
			return;
		} else {
			// test image contents
			var image = new Image();

			image.onerror = function() {
				callback('invalid_image', null);
				return;
			};

			image.onload = function() {
				// test dimensions
				if(image.width < 400 && image.height < 400) {
					callback('dim_small', null);
					return;
				}
				callback(null, fileReader.result);
			};

			image.src = fileReader.result;
		}

		return;
	};

	fileReader.readAsDataURL(file); // attempt reading

	return;
}

/**
	removes and restores input field
	calls restoreHandler()
*/
function restoreInput(inputID) {
	var inputHTML = $(inputID)[0].outerHTML; // copy
	restoreHandler($(inputID).attr('data-for')); // restore handler
	$(inputID).remove(); // remove
	$(inputHTML).insertAfter($('#post-form-description')); // restore input
	return;
}

/**
	Complete handler reset

	Set active='false'
	Set visible='false'
	Set src='/img/cross.png'
	Call drawHandlers()
*/
function restoreHandler(handlerContainerID) {
	$(handlerContainerID).attr('data-active', 'false'); // data-active 
	$(handlerContainerID).attr('data-visible', 'false'); // data-visible 
	$(handlerContainerID).find('.photo-handler').attr('src', '/img/cross.png'); // src
	drawHandlers();
}

/**
	First sets all containers data-visible='false'
	Displays active handlers and 1 inactive handler for use
	Removes and appends handler to preview-container
*/
function drawHandlers() {
	var newHandlerReady = false;
	// Hide
	$('.handler-container').each(function() {
		$(this).attr('data-visible', 'false');
	});
	// Show all with data-active='true'
	$('.handler-container').each(function() {
		if($(this).attr('data-active') == 'true') {
			$(this).attr('data-visible', 'true');
		}
	});
	// Show 1 inactive handler
	$('.handler-container').each(function() {
		if($(this).attr('data-active') == 'false' && !newHandlerReady) {
			$(this).attr('data-visible', 'true');
			newHandlerReady = true;
			var outerHTML = $(this)[0].outerHTML;
			$(this).remove();
			$(outerHTML).appendTo($('#preview-container'));
			return false;
		}
	});
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