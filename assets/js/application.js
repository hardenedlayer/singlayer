require('expose-loader?$!expose-loader?jQuery!jquery');
require("bootstrap/dist/js/bootstrap.js");

// enabling popovers and tooltips
$(document).ready(function(){
	$('[data-toggle="popover"]').popover();
});

$(document).ready(function(){
	$('[data-toggle="tooltip"]').tooltip();
});

// auto closing of alerts. but how can I avoid it for specific alert?
$(document).ready(function(){
	window.setTimeout(function() {
		$(".alert").alert('close');
	}, 10000);
});

// actor selector control:
$(document).ready(function(){
	$(".selector").on("click", function() {
		$("#hl_actor").text($(this).text());
		$(".selector").parent().removeClass("active");
		$(this).parent().addClass("active");
		document.cookie = "_singlayer_actor=" + $(this).text() + "; path=/";
		// auto reload is somewhat annoying sometime.
		//location.reload();
	});
});

// navbar highlight on current location:
$(document).ready(function() {
	$(".nav a:not('.selector')").parent().removeClass("active");
	$(".nav a:not('.selector')").each(function(index) {
		if ($(this).attr('href') == document.location.pathname) {
			$(this).parent().addClass("active");
		}
	});
});

$(() => {

});
