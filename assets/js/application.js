require('expose-loader?$!expose-loader?jQuery!jquery');
require("bootstrap/dist/js/bootstrap.js");

$(document).ready(function(){
	$('[data-toggle="popover"]').popover();
});

$(document).ready(function(){
	$('[data-toggle="tooltip"]').tooltip();
});

$(document).ready(function(){
	window.setTimeout(function() {
		$(".alert").alert('close');
	}, 5000);
});

$(document).ready(function(){
	$(".selector").on("click", function() {
		$("#hl_actor").text($(this).text());
		$(".selector").parent().removeClass("active");
		$(this).parent().addClass("active");
		document.cookie = "_singlayer_actor=" + $(this).text() + "; path=/";
	});
});

$(() => {

});
