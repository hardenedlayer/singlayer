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

$(() => {

});
