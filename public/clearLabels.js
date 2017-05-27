var infoLabels = document.getElementsByClassName("Info");
var passLabel = infoLabels[0];
var doiLabel = infoLabels[1];

var button = document.getElementsByClassName("login-button")[0];

// handling case when user corrects mistake in input fields
// clear both labels on button press, if anything is wrong
// server will rerender template and display error messages
button.addEventListener('click', function() {
    passLabel.innerHTML = "";
    doiLabel.innerHTML = "";
});